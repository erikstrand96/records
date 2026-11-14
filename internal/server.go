package internal

import (
	"database/sql"
)

type Server struct {
	Db *sql.DB
}
