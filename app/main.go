package main

import (
	"context"
	"fmt"
	"os"
	"scp-parser/internal/repository/postgresql"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
}

func main() {
	ctx := context.Background()
	dbconfig := postgresql.StorageConfig{
		DBUser:     os.Getenv("DBUser"),
		DBPass:     os.Getenv("DBPass"),
		DBHost:     os.Getenv("DBHost"),
		DBPort:     os.Getenv("DBPort"),
		DBName:     os.Getenv("DBName"),
		MaxAttemps: 5,
	}
	conn, err := postgresql.NewClient(ctx, dbconfig)
	if err != nil {
		fmt.Println("some error when pg client created")
	}
	q := `CREATE TABLE if not exists scpunits (
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
	res, err := conn.Exec(ctx, q)
	if err != nil {
		fmt.Println(fmt.Sprintf("some error when creating tables %v", err))
	}
	fmt.Println(res)

	// repo := db.Repository{
	// 	Client: *conn
	// }

	// conn, err := postgresql.NewClient(ctx, dbconfig)
	// if err != nil {
	// 	fmt.Println("some error when pg client created")
	// }
	// fmt.Print(conn.Ping(ctx))
	// var name string
	// err = conn.QueryRow(ctx, "SELECT name FROM spunits WHERE id='1'").Scan(&name)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	// 	os.Exit(1)
	// }
	// conn.Exec(ctx, "SELECT name FROM spunits")
	// fmt.Println(name)

	// fmt.Println(res)
	// db_url := dbconfig.UrlDb()
	// client := scpclient.ScpClient{
	// 	URL: os.Getenv("URL"),
	// 	Headers: map[string]string{
	// 		"User-Agent": os.Getenv("UserAgent"),
	// 		"Accept":     "application/json",
	// 	},
	// 	Client: &http.Client{},
	// }
	// data := []string{"/scp-1122", "/scp-001", "/scp-002"}
	// // arraySCP := scp_parser.ParseGetListSCP()[0:2]
	// for _, item := range data {
	// 	unit := client.ParseGetCurrentSCP(item)
	// 	fmt.Println(unit)
	// }
	// res := scp_parser.ParseGetCurrentSCP(data[0])
	// fmt.Println(res)
	// fmt.Println(res)
	// arraySCP := parseGetListSCP()
	// parseGetCurrentSCP(arraySCP)
	// fmt.Println(arraySCP)
}
