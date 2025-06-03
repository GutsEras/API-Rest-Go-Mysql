package main

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
)

func ConnectMockDB() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic("Erro ao criar sqlmock: " + err.Error())
	}
	return db, mock
}
