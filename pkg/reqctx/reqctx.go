package reqctx

import (
	"context"
	"log/slog"
	"sync"

	"github.com/samber/lo"
)

type ReqCtx struct {
	Ctx                    context.Context
	Logger                 *slog.Logger
	RequestID              int64
	Caller                 string
	Lock                   *sync.RWMutex
	values                 map[any]any
	customLogFields        map[string]any
	customLogFieldsOnError map[string]any
}

func NewReqCtx(ctx context.Context, logger *slog.Logger, requestID int64, caller string) *ReqCtx {
	return &ReqCtx{
		Ctx:                    ctx,
		Logger:                 logger.With("req_id", requestID),
		RequestID:              requestID,
		Caller:                 caller,
		Lock:                   new(sync.RWMutex),
		values:                 map[any]any{},
		customLogFields:        map[string]any{},
		customLogFieldsOnError: map[string]any{},
	}
}

func (ctx *ReqCtx) AddCustomLogField(key string, value any) {
	ctx.customLogFields[key] = value
}

func (ctx *ReqCtx) AddCustomLogFields(fields map[string]any) {
	for key, value := range fields {
		ctx.customLogFields[key] = value
	}
}

func (ctx *ReqCtx) AddCustomLogFieldOnError(key string, value any) {
	ctx.customLogFieldsOnError[key] = value
}

func (ctx *ReqCtx) AddCustomLogFieldsOnError(fields map[string]any) {
	for key, value := range fields {
		ctx.customLogFieldsOnError[key] = value
	}
}

func (ctx *ReqCtx) PutValue(key any, value any) {
	ctx.values[key] = value
}

func (ctx *ReqCtx) Clone() *ReqCtx {
	return &ReqCtx{
		Ctx:       ctx.Ctx,
		Logger:    ctx.Logger,
		RequestID: ctx.RequestID,
		Caller:    ctx.Caller,
		Lock:      new(sync.RWMutex),
		values: lo.MapEntries(ctx.values, func(key any, value any) (any, any) {
			return key, value
		}),
		customLogFields:        ctx.customLogFields,
		customLogFieldsOnError: ctx.customLogFieldsOnError,
	}
}

func (ctx *ReqCtx) CombineCustomLogFields(target []any) []any {
	for key, value := range ctx.customLogFields {
		target = append(target, key, value)
	}
	return target
}

func (ctx *ReqCtx) CombineCustomLogFieldsOnError(target []any) []any {
	for key, value := range ctx.customLogFieldsOnError {
		target = append(target, key, value)
	}
	return target
}

func GetValue[T any](ctx *ReqCtx, key any) T {
	return ctx.values[key].(T)
}
