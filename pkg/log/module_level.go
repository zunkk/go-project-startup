package log

import (
	"context"
	"log/slog"
)

type ModuleLevel struct {
	module string
	slog.Handler
}

func (g *ModuleLevel) Enabled(ctx context.Context, level slog.Level) bool {
	moduleLevel, ok := globalModuleLevelMap[g.module]
	if ok && level < moduleLevel {
		return false
	}

	return g.Handler.Enabled(ctx, level)
}

func (g *ModuleLevel) Handle(ctx context.Context, record slog.Record) error {
	record.Add(slog.String("module", g.module))
	return g.Handler.Handle(ctx, record)
}

func (g *ModuleLevel) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &ModuleLevel{
		module:  g.module,
		Handler: g.Handler.WithAttrs(attrs),
	}
}

func (g *ModuleLevel) WithGroup(name string) slog.Handler {
	return &ModuleLevel{
		module:  g.module,
		Handler: g.Handler.WithGroup(name),
	}
}
