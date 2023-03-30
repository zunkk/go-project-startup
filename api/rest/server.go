package rest

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/zunkk/go-project-startup/internal/coreapi"
	"github.com/zunkk/go-project-startup/internal/pkg/base"
	"github.com/zunkk/go-project-startup/internal/pkg/entity"
	"github.com/zunkk/go-project-startup/pkg/auth/jwt"
	"github.com/zunkk/go-project-startup/pkg/basic"
	"github.com/zunkk/go-project-startup/pkg/config"
	"github.com/zunkk/go-project-startup/pkg/errcode"
	"github.com/zunkk/go-project-startup/pkg/reqctx"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func init() {
	basic.RegisterComponents(New)
}

type Server struct {
	baseComponent *base.Component
	router        *gin.Engine
	listener      net.Listener
	hs            *http.Server
	*coreapi.CoreAPI
}

func New(baseComponent *base.Component, api *coreapi.CoreAPI) *Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	s := &Server{
		baseComponent: baseComponent,
		router:        router,
		hs: &http.Server{
			Addr:           fmt.Sprintf(":%d", baseComponent.Config.HTTP.Port),
			Handler:        router,
			ReadTimeout:    baseComponent.Config.HTTP.ReadTimeout.ToDuration(),
			WriteTimeout:   baseComponent.Config.HTTP.WriteTimeout.ToDuration(),
			MaxHeaderBytes: 1 << 20,
		},
		CoreAPI: api,
	}
	baseComponent.RegisterLifecycleHook(s)
	return s
}

func (s *Server) Start() error {

	err := s.init()
	if err != nil {
		return errors.Wrap(err, "register router failed")
	}

	s.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", s.baseComponent.Config.HTTP.Port))
	if err != nil {
		return err
	}
	printServerInfo := func() {
		s.baseComponent.Logger.Infof("Http server listen on: %d", s.baseComponent.Config.HTTP.Port)
	}

	s.baseComponent.SafeGoPersistentTask(func() {
		err := func() error {
			if s.baseComponent.Config.HTTP.TLSEnable {
				if _, err := os.Stat(s.baseComponent.Config.HTTP.TLSCertFilePath); err != nil {
					return errors.Wrapf(err, "tls_cert_file_path [%s] is invalid path", s.baseComponent.Config.HTTP.TLSCertFilePath)
				}
				if _, err := os.Stat(s.baseComponent.Config.HTTP.TLSKeyFilePath); err != nil {
					return errors.Wrapf(err, "tls_key_file_path [%s] is invalid path", s.baseComponent.Config.HTTP.TLSKeyFilePath)
				}
				printServerInfo()

				if err := s.hs.ServeTLS(s.listener, s.baseComponent.Config.HTTP.TLSCertFilePath, s.baseComponent.Config.HTTP.TLSKeyFilePath); err != nil {
					return err
				}
			} else {
				printServerInfo()
				if err := s.hs.Serve(s.listener); err != nil {
					return err
				}
			}
			return nil
		}()
		if err != nil && err != http.ErrServerClosed {
			s.baseComponent.Logger.WithFields(logrus.Fields{"err": err, "port": s.baseComponent.Config.HTTP.Port}).Warn("Failed to start http server")
			s.baseComponent.ComponentShutdown()
			return
		}
		s.baseComponent.Logger.Info("Http server shutdown")
	})

	return nil
}

func (s *Server) Stop() error {
	return s.hs.Close()
}

func (s *Server) init() error {
	s.router.MaxMultipartMemory = s.baseComponent.Config.HTTP.MultipartMemory
	s.router.Use(s.crossOriginMiddleware)

	{
		//v := s.router.Group("/api/v1")
		//{
		//	g := v.Group("/user")
		//}
	}

	// dev enable pprof
	if s.baseComponent.IsDevVersion() {
		s.router.GET("/debug/pprof/", IndexHandler())
		s.router.GET("/debug/pprof/heap", HeapHandler())
		s.router.GET("/debug/pprof/goroutine", GoroutineHandler())
		s.router.GET("/debug/pprof/allocs", AllocsHandler())
		s.router.GET("/debug/pprof/block", BlockHandler())
		s.router.GET("/debug/pprof/threadcreate", ThreadCreateHandler())
		s.router.GET("/debug/pprof/cmdline", CmdlineHandler())
		s.router.GET("/debug/pprof/profile", ProfileHandler())
		s.router.GET("/debug/pprof/symbol", SymbolHandler())
		s.router.GET("/debug/pprof/trace", TraceHandler())
		s.router.GET("/debug/pprof/mutex", MutexHandler())
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
	reqID := s.baseComponent.UUIDGenerator.Generate()
	ctx := reqctx.NewReqCtx(c.Request.Context(), s.baseComponent.Logger, int64(reqID), "")
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

func (s *Server) apiHandlerWrap(handler func(ctx *reqctx.ReqCtx, c *gin.Context) (res interface{}, err error), opts ...apiConfigOption) func(c *gin.Context) {
	cfg := newAPIConfig(opts...)
	return func(c *gin.Context) {
		ctx := s.generateRequestContext(c)
		startTime := time.Now()
		reqURI := c.Request.URL.Path
		var res interface{}
		err := s.baseComponent.RecoverExecute(func() error {
			if cfg.needAuth || cfg.needAdmin {
				token := c.GetHeader(config.JWTTokenHeaderKey)
				if token == "" {
					return errcode.ErrAuthCode.Wrap("token is empty")
				}

				var customClaims entity.CustomClaims
				id, err := jwt.ParseWithHMACKey(s.baseComponent.Config.HTTP.JWTTokenHMACKey, token, &customClaims)
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

		latencyTime := fmt.Sprintf("%6v", endTime.Sub(startTime))
		reqMethod := c.Request.Method

		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		logFields := logrus.Fields{
			"http_code": statusCode,
			"time_cost": latencyTime,
			"ip":        clientIP,
			"method":    reqMethod,
			"uri":       reqURI,
		}
		if ctx.Caller != "" {
			logFields["caller"] = ctx.Caller
		}

		if err != nil {
			s.failResponseWithErr(ctx, c, err)
			ctx.CombineCustomLogFields(logFields)
			ctx.CombineCustomLogFieldsOnError(logFields)
			ctx.Logger.WithFields(logFields).Error("API request failed")
			return
		}
		ctx.CombineCustomLogFields(logFields)
		ctx.Logger.WithFields(logFields).Info("API request")
		s.successResponseWithData(c, res)
	}
}

func (s *Server) failResponseWithErr(ctx *reqctx.ReqCtx, c *gin.Context, err error) {
	code := errcode.DecodeError(err)
	msg := err.Error()

	ctx.AddCustomLogField("err_code", code)
	ctx.AddCustomLogField("err_msg", msg)

	httpCode := http.StatusOK
	if strings.Contains(config.Version, "test") {
		httpCode = http.StatusInternalServerError
	}

	c.JSON(httpCode, gin.H{
		"code":    code,
		"message": msg,
	})
}

func (s *Server) successResponseWithData(c *gin.Context, data interface{}) {
	res := gin.H{
		"code": 0,
	}
	if data != nil {
		res["data"] = data
	}
	c.JSON(http.StatusOK, res)
}
