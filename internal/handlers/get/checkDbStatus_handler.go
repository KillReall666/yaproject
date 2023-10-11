package get

import (
	"net/http"
)

type PingChecker interface {
	ConnectionCheck() error
	LogInfo(args ...interface{})
}

type CheckHandler struct {
	db     PingChecker
	logger PingChecker
}

func NewCheckConnHandler(s PingChecker) *CheckHandler {
	return &CheckHandler{
		db:     s,
		logger: s,
	}
}

func (h *CheckHandler) CheckPingWithDb(w http.ResponseWriter, r *http.Request) {
	err := h.db.ConnectionCheck()
	if err != nil {
		w.WriteHeader(500)
		h.logger.LogInfo("connection with db not available")
	}
	w.WriteHeader(200)
	h.logger.LogInfo("connection established")
}
