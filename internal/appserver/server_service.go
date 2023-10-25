package appserver

import (
	"context"
	"github.com/KillReall666/yaproject/internal/fileutil"
	"github.com/KillReall666/yaproject/internal/handlers"
	"github.com/KillReall666/yaproject/internal/logger"
	"github.com/KillReall666/yaproject/internal/model"
	"github.com/KillReall666/yaproject/internal/storage"
	"github.com/KillReall666/yaproject/internal/storage/postgres"
)

type Service struct {
	repository Repository
	log        *logger.Logger
	fileIo     *fileutil.FileIoStruct
	db         *postgres.Database
	useDB      bool
}

type Repository interface {
	CountSetter(ctx context.Context, name string, count int64) error
	GaugeSetter(ctx context.Context, name string, gauge float64) error
	GaugeGetter(ctx context.Context, key string) (float64, error)
	CountGetter(ctx context.Context, key string) (int64, error)
	GetAllMetrics() string
}

func (s *Service) DBStatusCheck() error {
	err := s.db.DBStatusCheck()
	if err != nil {
		return err
	}
	return nil
}

func NewService(useDB bool, log *logger.Logger, fileIo *fileutil.FileIoStruct, db *postgres.Database, memStorage *storage.MemStorage) *Service {
	service := Service{
		log:    log,
		fileIo: fileIo,
		useDB:  useDB,
		db:     db,
	}
	if useDB {
		service.repository = db
	} else {
		service.repository = memStorage
	}
	return &service
}

func (s *Service) SaveMetrics(ctx context.Context, request *model.Metrics) error {
	if request.Counter != nil {
		err := s.repository.CountSetter(ctx, request.Name, handlers.ConvertToInt64(request.Counter))
		if err != nil {
			return err
		}
	}

	if request.Gauge != nil {
		err := s.repository.GaugeSetter(ctx, request.Name, handlers.ConvertToFloat64(request.Gauge))
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) GetFloatMetrics(ctx context.Context, request *model.Metrics) (float64, error) {
	value, err := s.repository.GaugeGetter(ctx, request.Name)
	return value, err

}

func (s *Service) GetCountMetrics(ctx context.Context, request *model.Metrics) (int64, error) {
	value, err := s.repository.CountGetter(ctx, request.Name)
	return value, err
}

func (s *Service) PrintForHTML() string {
	htmlPage := s.repository.GetAllMetrics()
	return htmlPage
}

func (s *Service) LogInfo(args ...interface{}) {
	s.log.Sugar.Info(args)
}
