package repository

import (
	"context"
	"fmt"
	"scp-parser/server/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type SCPRepo struct {
	DB pgx.Conn
}

func NewSCPRepository(db *pgx.Conn) SCPRepo {
	return SCPRepo{DB: *db}
}

type createSCPUnitDB struct {
	Name        pgtype.Text
	Class       pgtype.Text
	Structure   pgtype.Text
	Filial      pgtype.Text
	Anomaly     pgtype.Text
	Subject     []string
	Discription pgtype.Text
	SpecialCOD  pgtype.Text
	Property    []string
	Link        pgtype.Text
}

type GetSCPUnitDB struct {
	ID int
	createSCPUnitDB
}

func nilIfEmpty(s string) bool {
	if s == "" {
		return false
	}
	return true
}

func validateSCPUnit(unit domain.CreateSCPUnit) createSCPUnitDB {
	return createSCPUnitDB{
		Name:        pgtype.Text{String: unit.Name, Valid: nilIfEmpty(unit.Name)},
		Class:       pgtype.Text{String: unit.Class, Valid: nilIfEmpty(unit.Class)},
		Structure:   pgtype.Text{String: unit.Structure, Valid: nilIfEmpty(unit.Structure)},
		Filial:      pgtype.Text{String: unit.Filial, Valid: nilIfEmpty(unit.Filial)},
		Anomaly:     pgtype.Text{String: unit.Anomaly, Valid: nilIfEmpty(unit.Anomaly)},
		Subject:     unit.Subject,
		Discription: pgtype.Text{String: unit.Discription, Valid: nilIfEmpty(unit.Discription)},
		SpecialCOD:  pgtype.Text{String: unit.SpecialCOD, Valid: nilIfEmpty(unit.SpecialCOD)},
		Property:    unit.Property,
		Link:        pgtype.Text{String: unit.Link, Valid: true},
	}
}

func (r *SCPRepo) Create(ctx context.Context, unit domain.CreateSCPUnit) (*GetSCPUnitDB, error) {
	data := validateSCPUnit(unit)
	q := `INSERT INTO scpunits
		(name, class, structure, filial, anomaly, subject, discription, specialCOD, property, link)
	VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING id, name, class, structure, filial, anomaly, subject, discription, specialCOD, property, link`

	row := r.DB.QueryRow(
		ctx,
		q,
		data.Name,
		data.Class,
		data.Structure,
		data.Filial,
		data.Anomaly,
		data.Subject,
		data.Discription,
		data.SpecialCOD,
		data.Property,
		data.Link,
	)

	unitDB, err := r.scanSCPRow(row)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf("SQL Error: %s, Detail: %s, Where: %s, Code: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code)
			fmt.Println(newErr)
			return nil, newErr
		}
		return nil, err
	}
	return unitDB, err
}

func (r *SCPRepo) GetByID(ctx context.Context, id int) (*GetSCPUnitDB, error) {
	q := "SELECT id, name, class, structure, filial, anomaly, subject, discription, specialCOD, property, link FROM scpunits WHERE id = $1"

	row := r.DB.QueryRow(ctx, q, id)

	unit, err := r.scanSCPRow(row)

	if err != nil {
		return nil, err
	}

	return unit, nil
}

func (r *SCPRepo) GetListSCP(ctx context.Context, limit, offset int) ([]*GetSCPUnitDB, error) {
	q := "SELECT id, name, class, structure, filial, anomaly, subject, discription, specialCOD, property, link FROM scpunits ORDER BY name LIMIT $1 OFFSET $2"

	rows, err := r.DB.Query(ctx, q, limit, offset*limit)
	if err != nil {
		return nil, fmt.Errorf("Failed to get SCP list: %v", err)
	}

	defer rows.Close()

	var units []*GetSCPUnitDB

	for rows.Next() {
		unit, err := r.scanSCPRow(rows)

		if err != nil {
			return nil, err
		}

		units = append(units, unit)
	}

	return units, nil
}

func (r *SCPRepo) DeleteByID(ctx context.Context, id int) error {
	query := `DELETE FROM scpunits WHERE id = $1`

	result, err := r.DB.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete SCP with ID %d: %w", id, err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("SCP with ID %d not found", id)
	}

	return nil
}

func (r *SCPRepo) UpdateSCPUnitByID(ctx context.Context, id int, unit domain.CreateSCPUnit) (*GetSCPUnitDB, error) {
	data := validateSCPUnit(unit)
	query := `UPDATE scpunits SET 
	name = $1, class = $2, structure = $3, filial = $4, anomaly = $5, subject = $6, discription = $7, specialCOD = $8, property = $9, link = $10 
	WHERE id = $11 
	RETURNING id, name, class, structure, filial, anomaly, subject, discription, specialCOD, property, link`

	row := r.DB.QueryRow(ctx, query,
		data.Name,
		data.Class,
		data.Structure,
		data.Filial,
		data.Anomaly,
		data.Subject,
		data.Discription,
		data.SpecialCOD,
		data.Property,
		data.Link,
		id,
	)
	unitDB, err := r.scanSCPRow(row)
	if err != nil {
		return nil, err
	}
	return unitDB, nil
}

func (r *SCPRepo) scanSCPRow(query pgx.Row) (*GetSCPUnitDB, error) {
	var unit GetSCPUnitDB
	err := query.Scan(
		&unit.ID,
		&unit.Name,
		&unit.Class,
		&unit.Structure,
		&unit.Filial,
		&unit.Anomaly,
		&unit.Subject,
		&unit.Discription,
		&unit.SpecialCOD,
		&unit.Property,
		&unit.Link,
	)

	if err != nil {
		return nil, err
	}

	return &unit, nil
}
