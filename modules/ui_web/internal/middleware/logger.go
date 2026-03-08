package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// LoggerMiddleware cria um middleware para logging de requisições HTTP
type LoggerMiddleware struct {
	logger *slog.Logger
}

// NewLoggerMiddleware cria uma nova instância do middleware de logging
func NewLoggerMiddleware(logger *slog.Logger) *LoggerMiddleware {
	return &LoggerMiddleware{logger: logger}
}

// responseWriter é um wrapper para capturar o status code

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Handler retorna o middleware de logging
func (m *LoggerMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)

		m.logger.Info("HTTP request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", wrapped.statusCode),
			slog.Duration("duration", duration),
			slog.String("remote_addr", r.RemoteAddr),
		)
	})
}
