package service

import (
	"github.com/KillReall666/yaproject/internal/db"
	"github.com/KillReall666/yaproject/internal/fileutil"
	"github.com/KillReall666/yaproject/internal/logger"
	"github.com/KillReall666/yaproject/internal/model"
	"github.com/KillReall666/yaproject/internal/storage"
)

type Service struct {
	repository *storage.MemStorage
	log        *logger.Logger
	fileIo     *fileutil.FileIoStruct
	db         *db.Database
}

func (s *Service) DbStatusCheck() error {
	s.db.DbStatusCheck()
	return nil
}

func NewService(repo *storage.MemStorage, log *logger.Logger, fileIo *fileutil.FileIoStruct, db *db.Database) *Service {
	return &Service{
		repository: repo,
		log:        log,
		fileIo:     fileIo,
		db:         db,
	}
}

func (s *Service) SaveMetrics(request *model.Metrics) error {
	if request.Counter != nil {
		s.repository.CountSetter(request.Name, *request.Counter)
		return nil
	}

	if request.Gauge != nil {
		s.repository.GaugeSetter(request.Name, *request.Gauge)
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

func (s *Service) MetricsPrint() {
	s.repository.Print()
}

func (s *Service) LogInfo(args ...interface{}) {
	s.log.Sugar.Info(args)
}
