package docs

import (
	"net/http"

	"dunhayat-api/pkg/logger"

	"github.com/swaggo/http-swagger/v2"
	"github.com/swaggo/swag"
	"go.uber.org/zap"
)

type Handler struct {
	logger logger.Interface
}

func NewHandler(logger logger.Interface) *Handler {
	return &Handler{
		logger: logger,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc(
		"GET /swagger/doc.json",
		h.handleSwaggerJSON,
	)

	mux.HandleFunc(
		"/swagger/",
		httpSwagger.Handler(
			httpSwagger.URL("/swagger/doc.json"),
		),
	)
}

func (h *Handler) handleSwaggerJSON(
	w http.ResponseWriter,
	req *http.Request,
) {
	h.logger.Info(
		"Swagger JSON requested",
		zap.String("method", req.Method),
		zap.String("path", req.URL.Path),
	)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)

	doc, err := swag.ReadDoc()
	if err != nil {
		h.logger.Error(
			"Failed to read Swagger doc",
			zap.Error(err),
		)
		http.Error(
			w,
			"Failed to generate API documentation",
			http.StatusInternalServerError,
		)
		return
	}

	if _, err := w.Write([]byte(doc)); err != nil {
		h.logger.Error(
			"Failed to write Swagger doc",
			zap.Error(err),
		)
	}
}
