package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/pkg/errors"

	"github.com/zunkk/go-project-startup/pkg/util"
)

var fileCacheRegistry = map[string]struct{}{}

var fileCacheRegistryLock = new(sync.RWMutex)

type FileCache[V any] struct {
	cacheDir string
	id       string
}

func NewFileCache[V any](repoDir string, id string) (*FileCache[V], error) {
	fileCacheRegistryLock.Lock()
	defer fileCacheRegistryLock.Unlock()

	if _, ok := fileCacheRegistry[id]; ok {
		return nil, errors.Errorf("cache with id %s already exists", id)
	}
	cacheDir := filepath.Join(repoDir, "cache")
	if os.MkdirAll(cacheDir, 0755) != nil {
		return nil, errors.Errorf("failed to create cache dir %s", cacheDir)
	}

	fileCacheRegistry[id] = struct{}{}
	return &FileCache[V]{
		cacheDir: cacheDir,
		id:       id,
	}, nil
}

func (c *FileCache[V]) filePath() string {
	return filepath.Join(c.cacheDir, fmt.Sprintf("%s.json", c.id))
}

func (c *FileCache[V]) Put(v V) error {
	err := func() error {
		raw, err := json.Marshal(v)
		if err != nil {
			return err
		}
		if err := os.WriteFile(c.filePath(), raw, 0755); err != nil {
			return err
		}
		return nil
	}()
	if err != nil {
		return errors.Wrapf(err, "failed to put cache: %s", c.id)
	}

	return nil
}

func (c *FileCache[V]) Get() (*V, error) {
	v, err := func() (*V, error) {
		raw, err := os.ReadFile(c.filePath())
		if err != nil {
			return nil, err
		}
		var v V
		if err := json.Unmarshal(raw, &v); err != nil {
			return nil, err
		}
		return &v, nil
	}()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get cache: %s", c.id)
	}

	return v, nil
}

func (c *FileCache[V]) Has() bool {
	return util.FileExist(c.filePath())
}

func (c *FileCache[V]) Delete() {
	_ = os.Remove(c.filePath())
}
