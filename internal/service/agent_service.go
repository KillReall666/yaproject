package service

import (
	"github.com/KillReall666/yaproject/internal/client"
	"github.com/KillReall666/yaproject/internal/logger"
)

type AgentService struct {
	log    *logger.Logger
	client *client.Client
}

func NewAgentService(log *logger.Logger, client *client.Client) *AgentService {
	return &AgentService{
		log:    log,
		client: client,
	}
}

func (s *AgentService) LogInfo(args ...interface{}) {
	s.log.Sugar.Info(args)
}
