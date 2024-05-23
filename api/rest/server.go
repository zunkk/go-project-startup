package rest

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/zunkk/go-project-startup/internal/coreapi"
	"github.com/zunkk/go-project-startup/internal/pkg/base"
	"github.com/zunkk/go-project-startup/internal/pkg/entity"
	"github.com/zunkk/go-project-startup/pkg/auth/jwt"
	"github.com/zunkk/go-project-startup/pkg/errcode"
	"github.com/zunkk/go-project-startup/pkg/frame"
	glog "github.com/zunkk/go-project-startup/pkg/log"
	"github.com/zunkk/go-project-startup/pkg/repo"
	"github.com/zunkk/go-project-startup/pkg/reqctx"
)

var log = glog.WithModule("api")

func init() {
	frame.RegisterComponents(New)
}

type Server struct {
	sidecar  *base.CustomSidecar
	router   *gin.Engine
	listener net.Listener
	server   *http.Server
	*coreapi.CoreAPI
}

func New(sidecar *base.CustomSidecar, api *coreapi.CoreAPI) *Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	s := &Server{
		sidecar: sidecar,
		router:  router,
		server: &http.Server{
			Addr:           fmt.Sprintf(":%d", sidecar.Repo.Cfg.HTTP.Port),
			Handler:        router,
			ReadTimeout:    sidecar.Repo.Cfg.HTTP.ReadTimeout.ToDuration(),
			WriteTimeout:   sidecar.Repo.Cfg.HTTP.WriteTimeout.ToDuration(),
			MaxHeaderBytes: 1 << 20,
		},
		CoreAPI: api,
	}
	sidecar.RegisterLifecycleHook(s)
	return s
}

func (s *Server) Start() error {
	if !s.sidecar.Repo.Cfg.HTTP.Enable {
		return nil
	}

	err := s.init()
	if err != nil {
		return errors.Wrap(err, "register router failed")
	}

	s.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", s.sidecar.Repo.Cfg.HTTP.Port))
	if err != nil {
		return err
	}
	printServerInfo := func() {
		log.Info(fmt.Sprintf("Http server listen on: %d", s.sidecar.Repo.Cfg.HTTP.Port))
	}

	s.sidecar.SafeGoPersistentTask(func() {
		err := func() error {
			if s.sidecar.Repo.Cfg.HTTP.TLSEnable {
				if _, err := os.Stat(s.sidecar.Repo.Cfg.HTTP.TLSCertFilePath); err != nil {
					return errors.Wrapf(err, "tls_cert_file_path [%s] is invalid path", s.sidecar.Repo.Cfg.HTTP.TLSCertFilePath)
				}
				if _, err := os.Stat(s.sidecar.Repo.Cfg.HTTP.TLSKeyFilePath); err != nil {
					return errors.Wrapf(err, "tls_key_file_path [%s] is invalid path", s.sidecar.Repo.Cfg.HTTP.TLSKeyFilePath)
				}
				printServerInfo()

				if err := s.server.ServeTLS(s.listener, s.sidecar.Repo.Cfg.HTTP.TLSCertFilePath, s.sidecar.Repo.Cfg.HTTP.TLSKeyFilePath); err != nil {
					return err
				}
			} else {
				printServerInfo()
				if err := s.server.Serve(s.listener); err != nil {
					return err
				}
			}
			return nil
		}()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Warn("Failed to start http server", "err", err, "port", s.sidecar.Repo.Cfg.HTTP.Port)
			s.sidecar.ComponentShutdown()
			return
		}
		log.Info("Http server shutdown")
	})

	return nil
}

func (s *Server) Stop() error {
	if !s.sidecar.Repo.Cfg.HTTP.Enable {
		return nil
	}
	return s.server.Close()
}

func (s *Server) init() error {
	s.router.MaxMultipartMemory = s.sidecar.Repo.Cfg.HTTP.MultipartMemory
	s.router.Use(s.crossOriginMiddleware)

	{
		v := s.router.Group("/api/v1")
		{
			v.GET("/ping", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "pong",
				})
			})
		}
	}
	return nil
}

func (s *Server) crossOriginMiddleware(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
	c.Header("Access-Control-Allow-Headers", "token, origin, content-type, accept, is_zh")
	c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")

	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusOK)
	}
}

func (s *Server) generateRequestContext(c *gin.Context) *reqctx.ReqCtx {
	reqID := s.sidecar.UUIDGenerator.Generate()
	ctx := reqctx.NewReqCtx(c.Request.Context(), s.sidecar.Logger, int64(reqID), "")
	return ctx
}

type apiConfig struct {
	needAuth  bool
	needAdmin bool
}

type apiConfigOption func(*apiConfig)

func apiNeedAuth() apiConfigOption {
	return func(c *apiConfig) {
		c.needAuth = true
	}
}

func apiNeedAdmin() apiConfigOption {
	return func(c *apiConfig) {
		c.needAdmin = true
	}
}

func newAPIConfig(opts ...apiConfigOption) apiConfig {
	apiCfg := &apiConfig{
		needAuth:  false,
		needAdmin: false,
	}
	for _, opt := range opts {
		opt(apiCfg)
	}
	return *apiCfg
}

func (s *Server) apiHandlerWrap(handler func(ctx *reqctx.ReqCtx, c *gin.Context) (res any, err error), opts ...apiConfigOption) func(c *gin.Context) {
	cfg := newAPIConfig(opts...)
	return func(c *gin.Context) {
		ctx := s.generateRequestContext(c)
		startTime := time.Now()
		reqURI := c.Request.URL.Path
		var res any
		err := s.sidecar.RecoverExecute(func() error {
			if cfg.needAuth || cfg.needAdmin {
				token := c.GetHeader(repo.JWTTokenHeaderKey)
				if token == "" {
					return errcode.ErrAuthCode.Wrap("token is empty")
				}

				var customClaims entity.CustomClaims
				id, err := jwt.ParseWithHMACKey(s.sidecar.Repo.Cfg.HTTP.JWTTokenHMACKey, token, &customClaims)
				if err != nil {
					return errcode.ErrAuthCode.Wrap(err.Error())
				}
				if id == "" {
					return errcode.ErrAuthCode.Wrap("internal error: token data invalid: id is empty")
				}

				ctx.Caller = id
			}

			var err error
			res, err = handler(ctx, c)
			return err
		})
		endTime := time.Now()

		timeCost := fmt.Sprintf("%6v", endTime.Sub(startTime))
		reqMethod := c.Request.Method

		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()

		logFields := []any{
			"http_code", statusCode,
			"time_cost", timeCost,
			"ip", clientIP,
			"method", reqMethod,
			"uri", reqURI,
		}
		if ctx.Caller != "" {
			logFields = append(logFields, "caller", ctx.Caller)
		}

		if err != nil {
			s.failResponseWithErr(ctx, c, err)
			logFields = ctx.CombineCustomLogFields(logFields)
			logFields = ctx.CombineCustomLogFieldsOnError(logFields)
			log.Error("API request failed", logFields...)
			return
		}
		logFields = ctx.CombineCustomLogFields(logFields)
		log.Info("API request", logFields...)
		s.successResponseWithData(c, res)
	}
}

func (s *Server) failResponseWithErr(ctx *reqctx.ReqCtx, c *gin.Context, err error) {
	code := errcode.DecodeError(err)
	msg := err.Error()

	ctx.AddCustomLogField("err_code", code)
	ctx.AddCustomLogField("err_msg", msg)

	httpCode := http.StatusOK
	if strings.Contains(repo.Version, "test") {
		httpCode = http.StatusInternalServerError
	}

	c.JSON(httpCode, gin.H{
		"code":    code,
		"message": msg,
	})
}

func (s *Server) successResponseWithData(c *gin.Context, data any) {
	res := gin.H{
		"code": 0,
	}
	if data != nil {
		res["data"] = data
	}
	c.JSON(http.StatusOK, res)
}
