package postgresql

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func UrlDb() string {
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", os.Getenv("DBUser"), os.Getenv("DBPass"), os.Getenv("DBHost"), os.Getenv("DBPort"), os.Getenv("DBName"))
	return url
}

func CreateTable() {
	url := UrlDb()
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		fmt.Errorf("error when connection to db: %v", err)
	}
	defer conn.Close(ctx)
	err = conn.Ping(ctx)
	if err != nil {
		fmt.Errorf("Ошибка при проверке подключения: %v", err)
	}
	// var name string
	err = conn.QueryRow(ctx, "SELECT * FROM pg_catalog.pg_tables").Scan()
	if err != nil {
		fmt.Println("error select: %v", err)
	}
	// fmt.Println(name)
}
