package service

import (
	"github.com/KillReall666/yaproject/internal/model"
	"github.com/KillReall666/yaproject/internal/storage"
	"net/http"
)

type Service struct {
	repository *storage.MemStorage
}

func NewService(repo *storage.MemStorage) *Service {
	return &Service{
		repository: repo,
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

func (s *Service) PrintMetrics(w http.ResponseWriter) {
	s.repository.GetAllMetrics(w)
}
