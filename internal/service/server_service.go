package service

import (
	"github.com/KillReall666/yaproject/internal/db"
	"github.com/KillReall666/yaproject/internal/fileutil"
	"github.com/KillReall666/yaproject/internal/handlers"
	"github.com/KillReall666/yaproject/internal/logger"
	"github.com/KillReall666/yaproject/internal/model"
	"github.com/KillReall666/yaproject/internal/storage"
)

type Service struct {
	repository Repository //memstorage
	log        *logger.Logger
	fileIo     *fileutil.FileIoStruct
	db         *db.Database
	useDB      bool
}

type Repository interface {
	CountSetter(name string, count int64) error
	GaugeSetter(name string, gauge float64) error
	GaugeGetter(key string) (float64, error)
	CountGetter(key string) (int64, error)
	GetAllMetrics() string
}

func (s *Service) DBStatusCheck() error {
	s.db.DBStatusCheck()
	return nil
}

func NewService(useDB bool, log *logger.Logger, fileIo *fileutil.FileIoStruct, db *db.Database, memStorage *storage.MemStorage) *Service {
	service := Service{
		log:    log,
		fileIo: fileIo,
		useDB:  useDB,
	}
	if useDB {
		service.repository = db
	} else {
		service.repository = memStorage
	}
	return &service
}

func (s *Service) SaveMetrics(request *model.Metrics) error {
	if request.Counter != nil {
		s.repository.CountSetter(request.Name, handlers.ConvertToInt64(request.Counter))
		return nil
	}

	if request.Gauge != nil {
		s.repository.GaugeSetter(request.Name, handlers.ConvertToFloat64(request.Gauge))
		return nil
	}
	return nil
}

func (s *Service) GetFloatMetrics(request *model.Metrics) (float64, error) {
	value, err := s.repository.GaugeGetter(request.Name)
	return value, err

}

func (s *Service) GetCountMetrics(request *model.Metrics) (int64, error) {
	value, err := s.repository.CountGetter(request.Name)
	return value, err
}

func (s *Service) PrintForHTML() string {
	htmlPage := s.repository.GetAllMetrics()
	return htmlPage
}

func (s *Service) LogInfo(args ...interface{}) {
	s.log.Sugar.Info(args)
}
