package Models

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

var DB *sql.DB
