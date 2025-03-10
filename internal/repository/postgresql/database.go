package postgresql

import (
	"context"
	"fmt"
	repeatable "scp-parser/internal/repository/utils"
	"time"

	"github.com/jackc/pgx/v5"
)

type StorageConfig struct {
	DBUser     string
	DBPass     string
	DBHost     string
	DBPort     string
	DBName     string
	MaxAttemps int
}

func (c *StorageConfig) LinkDB() string {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", c.DBUser, c.DBPass, c.DBHost, c.DBPort, c.DBName)
	return url
}

func NewClient(ctx context.Context, c StorageConfig) (conn *pgx.Conn, err error) {
	dsn := c.LinkDB()
	err = repeatable.DoWithTries(func() error {
		ctx, cancel := context.WithTimeout(ctx, time.Duration(c.MaxAttemps)*time.Second)
		defer cancel()
		conn, err = pgx.Connect(ctx, dsn)
		if err != nil {
			return err
		}
		return nil
	}, c.MaxAttemps, time.Duration(c.MaxAttemps)*time.Second)
	if err != nil {
		fmt.Println("error do with tries db")
	}
	return conn, nil
}
