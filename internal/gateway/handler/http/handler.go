package http

import (
	_ "embed"
	"encoding/json"
	"net/http"
	"strings"

	gatewayusecase "github.com/vlm326/golang-course-hw-78/internal/gateway/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:embed openapi.json
var openAPISpec []byte

//go:embed swagger_index.html
var swaggerIndexHTML []byte

type Handler struct {
	getRepository gatewayusecase.GetRepositoryUseCase
}

func NewHandler(getRepository gatewayusecase.GetRepositoryUseCase) *Handler {
	return &Handler{getRepository: getRepository}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/repositories/{owner}/{repo}", h.handleGetRepository)
	mux.HandleFunc("GET /openapi.json", h.handleOpenAPI)
	mux.HandleFunc("GET /swagger/index.html", h.handleSwaggerUI)
	mux.HandleFunc("GET /swagger/", h.handleSwaggerUI)
}

func (h *Handler) handleGetRepository(w http.ResponseWriter, r *http.Request) {
	resp, err := h.getRepository.Execute(r.Context(), r.PathValue("owner"), r.PathValue("repo"))
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) handleOpenAPI(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(openAPISpec)
}

func (h *Handler) handleSwaggerUI(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(swaggerIndexHTML)
}

func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(payload)
}

func writeError(w http.ResponseWriter, err error) {
	grpcStatus, ok := status.FromError(err)
	if !ok {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "internal server error",
		})
		return
	}

	statusCode := http.StatusInternalServerError
	switch grpcStatus.Code() {
	case codes.InvalidArgument:
		statusCode = http.StatusBadRequest
	case codes.NotFound:
		statusCode = http.StatusNotFound
	case codes.DeadlineExceeded:
		statusCode = http.StatusGatewayTimeout
	case codes.Unavailable:
		statusCode = http.StatusBadGateway
	}

	message := grpcStatus.Message()
	if strings.TrimSpace(message) == "" || grpcStatus.Code() == codes.Internal {
		message = "internal server error"
	}

	writeJSON(w, statusCode, map[string]string{
		"error": message,
	})
}
