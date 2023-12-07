package getmetrics

import (
	"net/http"

	"github.com/KillReall666/yaproject/internal/logger"
)

//go:generate go run github.com/vektra/mockery/v2@v2.35.4 --name=DBStatusChecker

type DBStatusChecker interface {
	DBStatusCheck() error
}

//go:generate go run github.com/vektra/mockery/v2@v2.35.4 --name=Logger
type Logger interface {
	LogInfo(args ...interface{})
}

type DBStatusHandler struct {
	db     DBStatusChecker
	logger Logger
}

func NewCheckDBStatusHandler(d DBStatusChecker, l *logger.Logger) *DBStatusHandler {
	return &DBStatusHandler{
		db:     d,
		logger: l,
	}
}

func (h *DBStatusHandler) DBStatusCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed!", http.StatusNotFound)
		return
	}

	err := h.db.DBStatusCheck()
	if err != nil {
		w.WriteHeader(500)
		h.logger.LogInfo("connection with postgres not available", err)
	}

	w.WriteHeader(200)
	h.logger.LogInfo("connection established")
}
