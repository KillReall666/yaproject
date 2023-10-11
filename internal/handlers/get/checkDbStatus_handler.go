package get

import (
	"github.com/KillReall666/yaproject/internal/logger"
	"net/http"
)

type DBStatusChecker interface {
	DBStatusCheck() error
}

type Logger interface {
	LogInfo(args ...interface{})
}

type DbStatusHandler struct {
	db     DBStatusChecker
	logger Logger
}

func NewCheckDBStatusHandler(d DBStatusChecker, l *logger.Logger) *DbStatusHandler {
	return &DbStatusHandler{
		db:     d,
		logger: l,
	}
}

func (h *DbStatusHandler) DBStatusCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusNotFound)
		return
	}

	err := h.db.DBStatusCheck()
	if err != nil {
		w.WriteHeader(500)
		h.logger.LogInfo("connection with db not available")
	}

	w.WriteHeader(200)
	h.logger.LogInfo("connection established")
}
