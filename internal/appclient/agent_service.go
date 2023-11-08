package appclient

import (
	"github.com/KillReall666/yaproject/internal/client"
	"github.com/KillReall666/yaproject/internal/logger"
)

type AgentService struct {
	log    *logger.Logger
	client *client.Client
}

func (s *AgentService) LogInfo(args ...interface{}) {
	s.log.Sugar.Info(args)
}
