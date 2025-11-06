package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"scp-parser/pkg/cmd"
	"scp-parser/pkg/config"
	"scp-parser/server/domain"
	"scp-parser/server/repository"

	"github.com/jackc/pgx/v5/pgtype"
)

// SCPUnitDTO represents SCP unit data transfer object
// @Description SCP Unit information
type SCPUnitDTO struct {
	ID          int      `json:"id"`
	Name        *string  `json:"name"`
	Class       *string  `json:"class"`
	Structure   *string  `json:"structure"`
	Filial      *string  `json:"filial"`
	Anomaly     *string  `json:"anomaly"`
	Subject     []string `json:"subject"`
	Discription *string  `json:"discription"`
	SpecialCOD  *string  `json:"special_cod"`
	Property    []string `json:"property"`
	Link        *string  `json:"link"`
}

type SCPListResponse struct {
	Data       []*SCPUnitDTO `json:"data"`
	Pagination *Pagination   `json:"pagination"`
}

type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
	Pages int `json:"pages"`
}

type SCPService struct {
	repo *repository.SCPRepo
}

func NewSCPService(ctx context.Context, cfg *config.Config) (*SCPService, error) {
	conn, err := cmd.NewClient(ctx, &cfg.DB)
	if err != nil {
		slog.Error(fmt.Sprintf("Error when creating SCPService: %v", err))
		return nil, err
	}
	repository := repository.NewSCPRepository(conn)
	return &SCPService{
		repo: &repository,
	}, nil
}

func (r *SCPService) GetByID(ctx context.Context, id int) (*SCPUnitDTO, error) {
	if id <= 0 {
		slog.Error(fmt.Sprintf("Invalid ID: %d", id))
		return nil, domain.ErrBadParamInput
	}

	scpUnit, err := r.repo.GetByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result") {
			return nil, domain.ErrNotFound
		}
		slog.Error(fmt.Sprintf("Failed when get SCPunit from DB: %d", id))
		return nil, domain.ErrInternalServerError
	}

	if scpUnit == nil {
		slog.Error(fmt.Sprintf("SCP with ID %d not found", id))
		return nil, domain.ErrNotFound
	}

	dto := r.ConvertSCPUnitDBToSCPUnitDTO(scpUnit)
	return dto, nil
}

func (r *SCPService) GetListSCP(ctx context.Context, limit, offset int) (*SCPListResponse, error) {
	units, count, err := r.repo.GetListSCP(ctx, limit, offset)

	if err != nil {
		slog.Error(fmt.Sprintf("Error when Getting SCP units with limit: %d, offset: %d", limit, offset))
		slog.Error(fmt.Sprintf("%v", err))
		return nil, domain.ErrInternalServerError
	}

	dtos := make([]*SCPUnitDTO, len(units))

	for i, unit := range units {
		dtos[i] = r.ConvertSCPUnitDBToSCPUnitDTO(unit)
	}

	response := &SCPListResponse{
		Data: dtos,
		Pagination: &Pagination{
			Page:  offset,
			Limit: limit,
			Total: count,
			Pages: func(limit, count int) int {
				if count%limit >= 1 {
					return count/limit + 1
				} else {
					return count / limit
				}
			}(limit, count),
		},
	}
	slog.Warn(fmt.Sprintf("%v", response))
	return response, nil
}

func (r *SCPService) UpdateSCP(ctx context.Context, dto *domain.CreateSCPUnit, id int) (*SCPUnitDTO, error) {
	if dto.Name == "" || dto.Class == "" {
		return nil, domain.ErrBadParamInput
	}

	create, err := r.repo.UpdateSCPUnitByID(ctx, id, *dto)

	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil, domain.ErrNotFound
		}
		return nil, domain.ErrInternalServerError
	}

	unit := r.ConvertSCPUnitDBToSCPUnitDTO(create)

	return unit, nil
}

func (r *SCPService) DeleteSCP(ctx context.Context, id int) error {
	if id < 0 {
		slog.Error(fmt.Sprintf("Ivalid ID: %d", id))
		return domain.ErrBadParamInput
	}

	err := r.repo.DeleteByID(ctx, id)
	if err != nil {
		fmt.Println(err)
		if strings.Contains(err.Error(), "no rows in result") {
			return domain.ErrNotFound
		}
		slog.Error(fmt.Sprintf("Error when delete SCP: %v", err))
		return domain.ErrInternalServerError
	}

	return nil
}

func (r *SCPService) CreateSCP(ctx context.Context, dto *domain.CreateSCPUnit) (*SCPUnitDTO, error) {
	if dto.Name == "" || dto.Class == "" {
		return nil, domain.ErrBadParamInput
	}

	create, err := r.repo.Create(ctx, *dto)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return nil, domain.ErrConflict
		}
		return nil, domain.ErrInternalServerError
	}

	unit := r.ConvertSCPUnitDBToSCPUnitDTO(create)

	return unit, nil
}

func (r *SCPService) ConvertSCPUnitDBToSCPUnitDTO(unit *repository.GetSCPUnitDB) *SCPUnitDTO {
	return &SCPUnitDTO{
		ID:          unit.ID,
		Name:        r.convertPGTextToTextNil(unit.Name),
		Class:       r.convertPGTextToTextNil(unit.Class),
		Structure:   r.convertPGTextToTextNil(unit.Structure),
		Filial:      r.convertPGTextToTextNil(unit.Filial),
		Anomaly:     r.convertPGTextToTextNil(unit.Anomaly),
		Subject:     unit.Subject,
		Discription: r.convertPGTextToTextNil(unit.Discription),
		SpecialCOD:  r.convertPGTextToTextNil(unit.SpecialCOD),
		Property:    unit.Property,
		Link:        r.convertPGTextToTextNil(unit.Link),
	}
}

func (r *SCPService) convertPGTextToTextNil(val pgtype.Text) *string {
	if val.Valid && val.String != "" {
		return &val.String
	}
	return nil
}
