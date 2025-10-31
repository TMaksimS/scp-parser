package cmd

import (
	"context"
	"fmt"
	"time"

	"scp-parser/pkg/config"

	"github.com/jackc/pgx/v5"
)

func doWithTries(fn func() error, attemtps int, delay time.Duration) (err error) {
	for attemtps > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attemtps--

			continue
		}
		return nil
	}
	return
}

func getDBLink(cfg *config.PGConfig) string {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName)
	return url
}

func NewClient(ctx context.Context, cfg *config.PGConfig) (conn *pgx.Conn, err error) {
	dsn := getDBLink(cfg)
	err = doWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, time.Duration(cfg.MaxAttemps)*time.Second)
		defer cancel()
		conn, err = pgx.Connect(ctx, dsn)
		if err != nil {
			return err
		}
		return nil
	}, cfg.MaxAttemps, time.Duration(cfg.MaxAttemps)*time.Second)
	if err != nil {
		fmt.Println("error do with tries db")
	}
	return conn, nil
}

func CreateDB(ctx context.Context, conn *pgx.Conn) error {
	q := `CREATE TABLE if not exists scpunits (
		id SERIAL primary key,
		name varchar(500) UNIQUE,
		class varchar(500),
		structure varchar(200),
		filial varchar(50),
		anomaly varchar(200),
		subject varchar(200) array,
		discription text,
		specialCOD text,
		property varchar(200) array,
		link varchar(150) not null)`
	_, err := conn.Exec(ctx, q)
	if err != nil {
		fmt.Println(fmt.Sprintf("some error when creating tables %v", err))
		return err
	}
	fmt.Println("Table has been created")
	return nil
}
