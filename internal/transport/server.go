package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/1260124186-cc/maintenance-work-order-service/internal/domain"
	"github.com/1260124186-cc/maintenance-work-order-service/internal/service"
)

type Server struct {
	service *service.WorkOrderService
	mux     *http.ServeMux
}

func NewServer(repository service.Repository) *Server {
	server := &Server{
		service: service.NewWorkOrderService(repository),
		mux:     http.NewServeMux(),
	}
	server.routes()
	return server
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /health", s.health)
	s.mux.HandleFunc("GET /assets", s.listAssets)
	s.mux.HandleFunc("POST /work-orders", s.createWorkOrder)
	s.mux.HandleFunc("POST /work-orders/{id}/assign", s.assignWorkOrder)
	s.mux.HandleFunc("POST /work-orders/{id}/complete", s.completeWorkOrder)
	s.mux.HandleFunc("GET /summary", s.dailySummary)
}

func (s *Server) health(writer http.ResponseWriter, _ *http.Request) {
	writeJSON(writer, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) listAssets(writer http.ResponseWriter, request *http.Request) {
	assets, err := s.service.ListAssets(request.Context())
	if err != nil {
		writeError(writer, err)
		return
	}
	writeJSON(writer, http.StatusOK, assets)
}

func (s *Server) createWorkOrder(writer http.ResponseWriter, request *http.Request) {
	var input domain.CreateWorkOrderInput
	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		writeJSON(writer, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}
	order, err := s.service.Create(request.Context(), input)
	if err != nil {
		writeError(writer, err)
		return
	}
	writeJSON(writer, http.StatusCreated, order)
}

func (s *Server) assignWorkOrder(writer http.ResponseWriter, request *http.Request) {
	var input struct {
		Technician string `json:"technician"`
	}
	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		writeJSON(writer, http.StatusBadRequest, map[string]string{"error": "invalid JSON body"})
		return
	}
	order, err := s.service.Assign(request.Context(), request.PathValue("id"), input.Technician)
	if err != nil {
		writeError(writer, err)
		return
	}
	writeJSON(writer, http.StatusOK, order)
}

func (s *Server) completeWorkOrder(writer http.ResponseWriter, request *http.Request) {
	order, err := s.service.Complete(request.Context(), request.PathValue("id"))
	if err != nil {
		writeError(writer, err)
		return
	}
	writeJSON(writer, http.StatusOK, order)
}

func (s *Server) dailySummary(writer http.ResponseWriter, request *http.Request) {
	date := strings.TrimSpace(request.URL.Query().Get("date"))
	if date == "" {
		writeJSON(writer, http.StatusBadRequest, map[string]string{"error": "date is required"})
		return
	}
	summary, err := s.service.DailySummary(request.Context(), date)
	if err != nil {
		writeError(writer, err)
		return
	}
	writeJSON(writer, http.StatusOK, summary)
}

func writeError(writer http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	switch {
	case errors.Is(err, domain.ErrAssetNotFound), errors.Is(err, domain.ErrWorkOrderNotFound):
		status = http.StatusNotFound
	case errors.Is(err, domain.ErrAssetUnavailable), errors.Is(err, domain.ErrInvalidWorkOrder),
		errors.Is(err, domain.ErrUnsupportedPriority), errors.Is(err, domain.ErrInvalidStatusChange),
		errors.Is(err, domain.ErrTechnicianRequired):
		status = http.StatusBadRequest
	case errors.Is(err, context.Canceled):
		status = http.StatusRequestTimeout
	}
	writeJSON(writer, status, map[string]string{"error": err.Error()})
}

func writeJSON(writer http.ResponseWriter, status int, value any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(value)
}
