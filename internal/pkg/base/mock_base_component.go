package base

import (
	"context"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"

	"github.com/zunkk/go-project-startup/internal/pkg/config"
	"github.com/zunkk/go-project-startup/pkg/basic"
)

type mockLifecycle struct {
}

func (l *mockLifecycle) Append(fx.Hook) {

}

type mockShutdowner struct {
}

func (s *mockShutdowner) Shutdown(...fx.ShutdownOption) error {
	return nil
}

func NewMockBaseComponent(t *testing.T) *Component {
	cfg := config.DefaultConfig(filepath.Join(t.TempDir(), time.Now().String()))
	cfg.HTTP.Port = 0

	bc, err := basic.NewBaseComponent(&basic.BuildConfig{
		Ctx:       context.Background(),
		Logger:    logrus.New(),
		Wg:        new(sync.WaitGroup),
		Version:   "test",
		NodeIndex: 0,
	}, &mockLifecycle{}, &mockShutdowner{})
	assert.Nil(t, err)
	return &Component{
		BaseComponent: bc,
		Config:        cfg,
	}
}
