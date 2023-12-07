package storage

import (
	"context"

	"github.com/KillReall666/yaproject/internal/config"
	"github.com/KillReall666/yaproject/internal/fileutil"
	"github.com/KillReall666/yaproject/internal/logger"
	"github.com/KillReall666/yaproject/internal/storage/memstore"
	"github.com/KillReall666/yaproject/internal/storage/postgres"
)

type Repository interface {
	CountSetter(ctx context.Context, name string, count int64) error
	GaugeSetter(ctx context.Context, name string, gauge float64) error
	GaugeGetter(ctx context.Context, key string) (float64, error)
	CountGetter(ctx context.Context, key string) (int64, error)
	GetAllMetrics() string
}

func NewStore(cfg config.RunConfig, log *logger.Logger) (Repository, error) {
	if cfg.UseDB {
		log.LogInfo("use database")
		return postgres.NewDB(cfg.DefaultDBConnStr)
	} else {
		store := memstore.NewMemStorage()
		fl := fileutil.NewFileIo(cfg, store, log)
		fl.Run()
		return store, nil
	}
}
