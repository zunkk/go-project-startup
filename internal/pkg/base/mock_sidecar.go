package base

import (
	"context"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"

	"github.com/zunkk/go-project-startup/internal/pkg/config"
	"github.com/zunkk/go-project-startup/pkg/frame"
)

type mockLifecycle struct {
}

func (l *mockLifecycle) Append(fx.Hook) {}

type mockShutdowner struct {
}

func (s *mockShutdowner) Shutdown(...fx.ShutdownOption) error {
	return nil
}

func NewMockCustomSidecar(t *testing.T) *CustomSidecar {
	cfg := config.DefaultConfig(filepath.Join(t.TempDir(), time.Now().String()))
	cfg.HTTP.Port = 0

	bc, err := frame.NewSidecar(&frame.BuildConfig{
		Ctx:       context.Background(),
		Wg:        new(sync.WaitGroup),
		Version:   "test",
		NodeIndex: 0,
	}, &mockLifecycle{}, &mockShutdowner{})
	assert.Nil(t, err)
	return &CustomSidecar{
		Sidecar: bc,
		Config:  cfg,
	}
}
