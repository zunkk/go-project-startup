package reqctx

import (
	"context"
	"sync"

	"github.com/sirupsen/logrus"
)

type ReqCtx struct {
	Ctx                    context.Context
	Logger                 logrus.FieldLogger
	RequestID              int64
	Caller                 string
	Lock                   *sync.RWMutex
	values                 map[interface{}]interface{}
	customLogFields        map[string]interface{}
	customLogFieldsOnError map[string]interface{}
}

func NewReqCtx(ctx context.Context, logger logrus.FieldLogger, requestID int64, caller string) *ReqCtx {
	return &ReqCtx{
		Ctx: ctx,
		Logger: logger.WithFields(logrus.Fields{
			"req_id": requestID,
		}),
		RequestID:              requestID,
		Caller:                 caller,
		Lock:                   new(sync.RWMutex),
		values:                 map[interface{}]interface{}{},
		customLogFields:        map[string]interface{}{},
		customLogFieldsOnError: map[string]interface{}{},
	}
}

func (ctx *ReqCtx) AddCustomLogField(key string, value interface{}) {
	ctx.customLogFields[key] = value
}

func (ctx *ReqCtx) AddCustomLogFields(fields map[string]interface{}) {
	for key, value := range fields {
		ctx.customLogFields[key] = value
	}
}

func (ctx *ReqCtx) AddCustomLogFieldOnError(key string, value interface{}) {
	ctx.customLogFieldsOnError[key] = value
}

func (ctx *ReqCtx) AddCustomLogFieldsOnError(fields map[string]interface{}) {
	for key, value := range fields {
		ctx.customLogFieldsOnError[key] = value
	}
}

func (ctx *ReqCtx) PutValue(key interface{}, value interface{}) {
	ctx.values[key] = value
}

func (ctx *ReqCtx) Clone() *ReqCtx {
	c := make(map[interface{}]interface{}, len(ctx.values))
	for s, i := range ctx.values {
		c[s] = i
	}
	return &ReqCtx{
		Ctx:                    ctx.Ctx,
		Logger:                 ctx.Logger,
		RequestID:              ctx.RequestID,
		Caller:                 ctx.Caller,
		Lock:                   new(sync.RWMutex),
		values:                 c,
		customLogFields:        ctx.customLogFields,
		customLogFieldsOnError: ctx.customLogFieldsOnError,
	}
}

func (ctx *ReqCtx) CombineCustomLogFields(target map[string]interface{}) {
	for key, value := range ctx.customLogFields {
		target[key] = value
	}
}

func (ctx *ReqCtx) CombineCustomLogFieldsOnError(target map[string]interface{}) {
	for key, value := range ctx.customLogFieldsOnError {
		target[key] = value
	}
}

func GetValue[T any](ctx *ReqCtx, key interface{}) T {
	return ctx.values[key].(T)
}
