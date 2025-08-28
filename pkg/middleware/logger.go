package middleware

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"dunhayat-api/pkg/logger"

	"go.uber.org/zap"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode        int
	body              *bytes.Buffer
	teeWriter         io.Writer
	writeCalled       bool
	writeHeaderCalled bool
}

func (rw *responseWriter) WriteHeader(code int) {
	if !rw.writeHeaderCalled {
		rw.statusCode = code
		rw.writeHeaderCalled = true
		rw.ResponseWriter.WriteHeader(code)
	}
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.writeHeaderCalled {
		rw.WriteHeader(http.StatusOK)
	}

	rw.writeCalled = true

	_, err := rw.teeWriter.Write(b)
	if err != nil {
		return 0, err
	}

	return rw.ResponseWriter.Write(b)
}

func RequestLogger(log logger.Interface) func(http.Handler) http.Handler {
	return func(
		next http.Handler,
	) http.Handler {
		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			start := time.Now()

			if r.URL.Path == "/health" || r.URL.Path == "/" {
				next.ServeHTTP(w, r)
				return
			}

			var requestBody []byte
			if r.Body != nil &&
				r.ContentLength > 0 &&
				r.ContentLength < 1024*1024 {
				requestBody, _ = io.ReadAll(r.Body)
				r.Body = io.NopCloser(
					bytes.NewBuffer(requestBody),
				)
			}

			bodyBuffer := &bytes.Buffer{}
			responseWriter := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
				body:           bodyBuffer,
				teeWriter:      io.MultiWriter(w, bodyBuffer),
			}

			next.ServeHTTP(responseWriter, r)

			if !responseWriter.writeHeaderCalled {
				responseWriter.statusCode = http.StatusOK
			}

			duration := time.Since(start)

			fields := []zap.Field{
				zap.String(
					"method",
					r.Method,
				),
				zap.String(
					"path",
					r.URL.Path,
				),
				zap.String(
					"query",
					r.URL.RawQuery,
				),
				zap.String(
					"remote_addr",
					r.RemoteAddr,
				),
				zap.String(
					"user_agent",
					r.UserAgent(),
				),
				zap.String(
					"referer",
					r.Referer(),
				),
				zap.Int(
					"status_code",
					responseWriter.statusCode,
				),
				zap.Duration(
					"duration",
					duration,
				),
				zap.Int64(
					"content_length",
					r.ContentLength,
				),
				zap.Bool(
					"write_header_called",
					responseWriter.writeHeaderCalled,
				),
				zap.Bool(
					"write_called",
					responseWriter.writeCalled,
				),
			}

			if r.Header.Get("Authorization") != "" {
				fields = append(
					fields,
					zap.String("authorization", "***"),
				)
			}
			if r.Header.Get("Cookie") != "" {
				fields = append(
					fields,
					zap.String("cookie", "***"),
				)
			}

			if len(requestBody) > 0 && len(requestBody) <= 1024 {
				fields = append(
					fields,
					zap.String(
						"request_body",
						string(requestBody),
					),
				)
			}

			responseBodyLen := bodyBuffer.Len()
			if responseBodyLen > 0 && responseBodyLen <= 1024 {
				fields = append(
					fields,
					zap.String(
						"response_body",
						bodyBuffer.String(),
					),
				)
			}

			switch {
			case responseWriter.statusCode >= 500:
				log.Error("HTTP Request Error", fields...)
			case responseWriter.statusCode >= 400:
				log.Warn("HTTP Request Warning", fields...)
			case responseWriter.statusCode >= 300:
				log.Info("HTTP Request Redirect", fields...)
			default:
				log.Info("HTTP Request", fields...)
			}

			if !responseWriter.writeHeaderCalled {
				log.Warn(
					"WriteHeader never called - potential handler bug",
					zap.String("path", r.URL.Path),
					zap.String("method", r.Method),
				)
			}
		})
	}
}
