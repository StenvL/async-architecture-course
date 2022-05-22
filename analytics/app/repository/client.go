package repository

import "github.com/jmoiron/sqlx"

type Repository struct {
	db *sqlx.DB
}

func New() (Repository, error) {
	db, err := sqlx.Connect(
		"postgres",
		"host=localhost port=5432 user=postgres password=password dbname=analytics sslmode=disable",
	)
	if err != nil {
		return Repository{}, err
	}

	return Repository{db: db}, nil
}
