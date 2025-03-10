package db

import (
	"context"
	"fmt"
	"os"
	"scp-parser/domain/scpunit"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repository struct {
	Client pgx.Conn
}

func (r *Repository) CreateDB(ctx context.Context) {
	var err any
	q := `CREATE TABLE scpunits (
		id SERIAL primary key,
		name varchar(500),
		class varchar(500),
		structure varchar(200),
		filial varchar(50),
		anomaly varchar(200),
		subject varchar(200) array,
		discription text,
		specialCOD text,
		property varchar(200) array,
		link varchar(150) not null)`
	conn, err = r.client.Exec(ctx, q).Scan()
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
}

func (r *Repository) CreateUnit(ctx context.Context, unit scpunit.SCPUnitDB) (string, error) {
	q := `INSERT INTO scpunits
		(name, class, structure, filial, anomaly, subject, discription, specialCOD, property, link)
	VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING id`
	if err := r.client.QueryRow(
		ctx,
		q,
		unit.Name,
		unit.Class,
		unit.Structure,
		unit.Filial,
		unit.Anomaly,
		unit.Subject,
		unit.Discription,
		unit.SpecialCOD,
		unit.Property,
		unit.Link,
	).Scan(scpunit.SCPUnitDB.ID); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf("SQL Error: %s, Detail: %s, Where: %s, Code: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code)
			fmt.Println(newErr)
			return "", newErr
		}
		return "", nil
	}
	return scpunit.SCPUnitDB.ID, nil
}

// func NewRepository(client pgx.Conn) {
// 	return &repository{
// 		client: client,
// 	}
// }
