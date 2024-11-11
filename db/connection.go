package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect() (*sql.DB, error) {
	url := os.Getenv("S2_DATABASE_URL")

	if url == "" {
		log.Fatalf("missing or empty env variable S2_DATABASE_URL\n")
	}

	db, err := sql.Open("pgx", url)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
