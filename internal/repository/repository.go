package repository

import "database/sql"

type DatabaseRepo interface {
	Connection() *sql.DB
}
