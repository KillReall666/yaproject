package get

import (
	"github.com/KillReall666/yaproject/internal/logger"
	"net/http"
)

type DbStatusChecker interface {
	DbStatusCheck() error
}

type Logger interface {
	LogInfo(args ...interface{})
}

type DbStatusHandler struct {
	db     DbStatusChecker
	logger Logger
}

func NewCheckDbStatusHandler(d DbStatusChecker, l *logger.Logger) *DbStatusHandler {
	return &DbStatusHandler{
		db:     d,
		logger: l,
	}
}

func (h *DbStatusHandler) DbStatusCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusNotFound)
		return
	}

	err := h.db.DbStatusCheck()
	if err != nil {
		w.WriteHeader(500)
		h.logger.LogInfo("connection with db not available")
	}

	w.WriteHeader(200)
	h.logger.LogInfo("connection established")
}
