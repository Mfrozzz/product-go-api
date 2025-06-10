package db

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func ConnectDB() (*sql.DB, error) {
	godotenv.Load(os.ExpandEnv("../.env"))
	var (
		host      = os.Getenv("DB_HOST")
		port, err = strconv.Atoi(os.Getenv("DB_PORT"))
		user      = os.Getenv("DB_USER")
		password  = os.Getenv("DB_PASSWORD")
		dbname    = os.Getenv("DB_NAME")
	)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to " + dbname)

	return db, nil
}
