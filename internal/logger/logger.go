package logger

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Logger struct {
	Logger *zap.Logger
	Sugar  zap.SugaredLogger
}

type responseData struct {
	status int
	size   int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func InitLogger() (*Logger, error) {
	mylogger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}

	defer mylogger.Sync()
	sugar := *mylogger.Sugar()

	return &Logger{
		Logger: mylogger,
		Sugar:  sugar,
	}, nil
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func (l Logger) MyLogger(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		respData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   respData,
		}
		h.ServeHTTP(&lw, r)

		duration := time.Since(start)

		l.Sugar.Infoln(
			"uri:", r.RequestURI,
			"method:", r.Method,
			"duration:", duration,
			"status:", respData.status,
			"size:", respData.size,
		)

	}
	return http.HandlerFunc(logFn)
}