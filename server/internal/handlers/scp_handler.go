package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
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

// ListSCP list all existing units
// @Summary Get list of SCP units
// @Description Get paginated list of SCP units
// @Tags SCP
// @Accept json
// @Produce json
// @Param limit query int false "Limit the number of results" default(50) minimum(1) maximum(100)
// @Param offset query int false "Offset for pagination" default(0) minimum(0)
// @Success 200 {object} service.SCPUnitDTO
// @Failure 500
// @Router /scp/ [get]
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

// CreateSCP Create new SCP unit
// @Summary Create SCP unit
// @Description Create a new SCP unit
// @Tags SCP
// @Accept json
// @Produce json
// @Param request body domain.CreateSCPUnit true "SCP data"
// @Success 200 {object} service.SCPUnitDTO
// @Failure 400
// @Failure 500
// @Router /scp/ [post]
func (h *SCPHandler) CreateSCP(w http.ResponseWriter, r *http.Request) {
	var body domain.CreateSCPUnit

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		slog.Error(fmt.Sprintf("Error when decode SCP: %v", err))
		http.Error(w, domain.ErrBadParamInput.Error(), getStatusCode(domain.ErrBadParamInput))
		return
	}

	// TODO validate body to avoid 500

	createdSCP, err := h.service.CreateSCP(r.Context(), &body)

	if err != nil {
		http.Error(w, err.Error(), getStatusCode(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdSCP)

}

// ListSCP current SCP
// @Summary Get SCP unit
// @Description Get SCP unit
// @Tags SCP
// @Accept json
// @Produce json
// @Param id path int true "SCP ID"
// @Success 200 {object} service.SCPUnitDTO
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /scp/{id} [get]
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

// UpdateSCP Update SCP unit
// @Summary Update SCP unit
// @Description Update SCP unit
// @Tags SCP
// @Accept json
// @Produce json
// @Param id path int true "SCP ID"
// @Param request body domain.CreateSCPUnit true "SCP data"
// @Success 200 {object} service.SCPUnitDTO
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /scp/{id} [put]
func (h *SCPHandler) UpdateSCP(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid SCP ID", http.StatusBadRequest)
		return
	}

	var body domain.CreateSCPUnit

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		slog.Error(fmt.Sprintf("Error when decode SCP: %v", err))
		http.Error(w, domain.ErrBadParamInput.Error(), getStatusCode(domain.ErrBadParamInput))
		return
	}
	// TODO validate body to avoid 500

	updatedSCP, err := h.service.UpdateSCP(r.Context(), &body, id)
	if err != nil {
		http.Error(w, err.Error(), getStatusCode(err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedSCP)

}

// DeleteSCP Delete SCP unit
// @Summary Delete SCP unit
// @Description Delete SCP unit
// @Tags SCP
// @Accept json
// @Produce json
// @Param id path int true "SCP ID"
// @Success 204
// @Failure 404
// @Failure 500
// @Router /scp/{id} [delete]
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
