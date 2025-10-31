package repository

import (
	"context"
	"fmt"
	"scp-parser/parser/service"

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

func nilIfEmpty(s string) bool {
	if s == "" {
		return false
	}
	return true
}

func validateSCPUnit(unit service.SCPUnit) createSCPUnitDB {
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

func (r *SCPRepo) Create(ctx context.Context, unit service.SCPUnit) error {
	data := validateSCPUnit(unit)
	q := `INSERT INTO scpunits
		(name, class, structure, filial, anomaly, subject, discription, specialCOD, property, link)
	VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING id`
	err := r.DB.QueryRow(
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
	).Scan("id")
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf("SQL Error: %s, Detail: %s, Where: %s, Code: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code)
			fmt.Println(newErr)
			return newErr
		}
		return nil
	}
	return nil
}
