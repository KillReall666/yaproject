package get

import (
	"github.com/KillReall666/yaproject/internal/handlers/get/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDBStatusHandler_DBStatusCheck(t *testing.T) {
	type fields struct {
		db     DBStatusChecker
		logger Logger
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "simple test",
			args: args{
				w: httptest.NewRecorder(),
				r: &http.Request{
					Method: http.MethodGet,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := mocks.NewDBStatusChecker(t)
			logger := mocks.NewLogger(t)

			db.On("DBStatusCheck").Return(nil)
			logger.On("LogInfo", "connection established").Return(nil)

			h := &DBStatusHandler{
				db:     db,
				logger: logger,
			}
			h.DBStatusCheck(tt.args.w, tt.args.r)
		})
	}
}
