package Models

import (
	"GoIris/Models/Database"
	"GoIris/Models/Hashes"
	"database/sql"
	"log"

	_ "modernc.org/sqlite"

)

func init() {
	Hashes.CompileRegex()
	var err error
	DB, err = sql.Open("sqlite", "GoIris.db")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	if err := Database.CreateTables(DB); err != nil {
		log.Fatalf("Error creating tables: %v", err)
	}
}
