package logger

import (
	"Effective_Mobile/internal/logger"
	"net/http"
	"time"
)

func Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next(rw, r)

		duration := time.Since(start)

		if rw.statusCode >= 400 {
			logger.Error("→ %s %s | %d | %v | IP: %s | UA: %s | Err: %s",
				r.Method,
				r.URL.Path,
				rw.statusCode,
				duration,
				r.RemoteAddr,
				r.UserAgent(),
				rw.errMsg,
			)
		} else {
			logger.Info("→ %s %s | %d | %v | IP: %s | UA: %s",
				r.Method,
				r.URL.Path,
				rw.statusCode,
				duration,
				r.RemoteAddr,
				r.UserAgent(),
			)
		}
	}
}

type loggingResponseWriter struct {
	http.ResponseWriter
	errMsg     string
	statusCode int
	written    bool
}

func (w *loggingResponseWriter) WriteHeader(code int) {
	if !w.written {
		w.statusCode = code
		w.ResponseWriter.WriteHeader(code)
		w.written = true
	}
}

func (w *loggingResponseWriter) Write(b []byte) (int, error) {
	if !w.written {
		w.WriteHeader(http.StatusOK)
	}
	if w.statusCode >= 400 {
		const maxLen = 512
		if len(b) > maxLen {
			w.errMsg = string(b[:maxLen]) + "..."
		} else {
			w.errMsg = string(b)
		}
	}
	return w.ResponseWriter.Write(b)
}
