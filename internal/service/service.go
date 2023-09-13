package service

import (
	"github.com/KillReall666/yaproject/internal/metrics"
	"github.com/KillReall666/yaproject/internal/model"
	"github.com/KillReall666/yaproject/internal/storage"
)

type Service struct {
	repository     *storage.MemStorage
	metricsStorage *metrics.GaugeMetricsGetter
}

func NewService(repo *storage.MemStorage, memRepo *metrics.GaugeMetricsGetter) *Service {
	return &Service{
		repository:     repo,
		metricsStorage: memRepo,
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
