package src

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func UrlDb() string {
	err := godotenv.Load()
	if err != nil {
		fmt.Errorf("Failed load .env")
	}
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", os.Getenv("DBUser"), os.Getenv("DBPass"), os.Getenv("DBName"), os.Getenv("DBHost"), os.Getenv("DBPort"))
	return url
}
