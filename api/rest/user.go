package rest

import (
	"github.com/gin-gonic/gin"

	"github.com/zunkk/go-project-startup/pkg/reqctx"
)

type userQueryReq struct {
	ID int64 `uri:"id" binding:"required"`
}

func (s *Server) userQuery(ctx *reqctx.ReqCtx, c *gin.Context) (res any, err error) {
	req := &userQueryReq{}
	if err := c.ShouldBindUri(req); err != nil {
		return nil, err
	}
	return s.UserSrv.QueryByID(ctx.Ctx, req.ID)
}
