package service

import (
	"github.com/KillReall666/yaproject/internal/model"
	"github.com/KillReall666/yaproject/internal/storage"
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
		s.repository.Count(request.Name, *request.Counter)
		return nil
	}

	if request.Gauge != nil {
		s.repository.Gauge(request.Name, *request.Gauge)
		return nil
	}
	return nil
}
