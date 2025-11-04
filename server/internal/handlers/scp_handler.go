package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"scp-parser/pkg/config"
	"scp-parser/server/domain"
	"scp-parser/server/internal/service"

	"github.com/go-chi/chi/v5"
)

type SCPHandler struct {
	service *service.SCPService
}

func NewSCPHandler(ctx context.Context, cfg *config.Config) (*SCPHandler, error) {
	serviceSCP, err := service.NewSCPService(ctx, cfg)
	if err != nil {
		return nil, err
	}
	return &SCPHandler{
		service: serviceSCP,
	}, nil
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	// slog.Error(fmt.Sprintf("Error: %v", err))

	switch err {
	case domain.ErrNotFound:
		return http.StatusNotFound
	case domain.ErrBadParamInput:
		return http.StatusBadRequest
	case domain.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

func (h *SCPHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", h.GetSCPlist)
	r.Post("/", h.CreateSCP)
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", h.GetSCP)
		r.Put("/", h.UpdateSCP)
		r.Delete("/", h.DeleteSCP)
	})
	return r
}

func (h *SCPHandler) GetSCPlist(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 50
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}

	result, err := h.service.GetListSCP(r.Context(), limit, offset)
	if err != nil {
		http.Error(w, err.Error(), getStatusCode(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *SCPHandler) GetSCP(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid SCP ID", http.StatusBadRequest)
		return
	}

	scp, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), getStatusCode(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(scp)

}

func (h *SCPHandler) CreateSCP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var body domain.CreateSCPUnit

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, `{"error": "invalid JSON"}`, getStatusCode(domain.ErrBadParamInput))
		return
	}

	createdSCP, err := h.service.CreateSCP(r.Context(), &body)

	if err != nil {
		http.Error(w, err.Error(), getStatusCode(err))
		return
	}

	json.NewEncoder(w).Encode(createdSCP)

}

func (h *SCPHandler) UpdateSCP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid SCP ID", http.StatusBadRequest)
		return
	}

	var body domain.CreateSCPUnit

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, `{"error": "invalid JSON"}`, getStatusCode(domain.ErrBadParamInput))
		return
	}

	updatedSCP, err := h.service.UpdateSCP(r.Context(), &body, id)
	if err != nil {
		http.Error(w, err.Error(), getStatusCode(err))
		return
	}

	json.NewEncoder(w).Encode(updatedSCP)

}

func (h *SCPHandler) DeleteSCP(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid SCP ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteSCP(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), getStatusCode(err))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("OK")

}
