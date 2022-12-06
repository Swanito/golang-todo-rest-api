package server

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

func ConnectDatabase(config DBConfig) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.Username, config.Password, config.DBName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	err = db.Ping()
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return db, nil
}
