package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type DBInfo struct {
	Name     string
	Host     string
	Port     string
	User     string
	Password string
}

func (db *DBInfo) Connect(driver string) (*sqlx.DB, error) {
	return sqlx.Connect(driver, fmt.Sprintf(
		"dbname=%s host=%s port=%s user=%s password=%s sslmode=disable",
		db.Name,
		db.Host,
		db.Port,
		db.User,
		db.Password,
	))
}
