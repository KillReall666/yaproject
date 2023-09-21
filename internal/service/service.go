package service

import (
	"github.com/KillReall666/yaproject/internal/model"
	"github.com/KillReall666/yaproject/internal/storage"
)

type service struct {
	repository *storage.MemStorage
}

func NewService(repo *storage.MemStorage) *service {
	return &service{
		repository: repo,
	}
}

func (s *service) SaveMetrics(request *model.Metrics) error {
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

func (s *service) GetFloatMetrics(request *model.Metrics) (float64, error) {
	value, err := s.repository.GaugeGetter(request.Name)
	return value, err

}

func (s *service) GetCountMetrics(request *model.Metrics) (int64, error) {
	value, err := s.repository.CountGetter(request.Name)
	return value, err
}

func (s *service) PrintForHTML() string {
	htmlPage := s.repository.GetAllMetrics()
	return htmlPage
}

func (s *service) MetricsPrint() {
	s.repository.Print()
}
