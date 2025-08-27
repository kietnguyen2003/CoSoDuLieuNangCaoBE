package database

import (
	"clinic-management/internal/utils"
	"database/sql"
	"fmt"

	_ "github.com/denisenkom/go-mssqldb"
)

func Connect(databaseURL string) (*sql.DB, error) {
	fmt.Println("Connecting to database with url:", databaseURL)
	db, err := sql.Open("sqlserver", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}
	utils.InitializeCounters()

	return db, nil
}
