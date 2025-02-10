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
	"github.com/zunkk/go-sidecar/frame"
	"github.com/zunkk/go-sidecar/repo"
)

type mockLifecycle struct {
}

func (l *mockLifecycle) Append(fx.Hook) {}

type mockShutdowner struct {
	t testing.TB
}

func (s *mockShutdowner) Shutdown(...fx.ShutdownOption) error {
	s.t.Fatal("Shutdown called")
	return nil
}

func NewMockCustomSidecar(t testing.TB) *CustomSidecar {
	cfg := config.DefaultConfig()
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
		Repo: &repo.Repo[*config.Config]{
			RepoPath: filepath.Join(t.TempDir(), time.Now().String()),
			Cfg:      cfg,
		},
	}
}
