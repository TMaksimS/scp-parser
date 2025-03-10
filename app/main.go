package main

import (
	"fmt"
	"net/http"
	"os"
	"scp-parser/internal/repository/postgresql"
	scpclient "scp-parser/internal/repository/scp-parser"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
}

func main() {
	postgresql.CreateTable()
	client := scpclient.ScpClient{
		URL: os.Getenv("URL"),
		Headers: map[string]string{
			"User-Agent": os.Getenv("UserAgent"),
			"Accept":     "application/json",
		},
		Client: &http.Client{},
	}
	data := []string{"/scp-1122", "/scp-001", "/scp-002"}
	// arraySCP := scp_parser.ParseGetListSCP()[0:2]
	for _, item := range data {
		unit := client.ParseGetCurrentSCP(item)
		fmt.Println(unit)
	}
	// res := scp_parser.ParseGetCurrentSCP(data[0])
	// fmt.Println(res)
	// fmt.Println(res)
	// arraySCP := parseGetListSCP()
	// parseGetCurrentSCP(arraySCP)
	// fmt.Println(arraySCP)
}
