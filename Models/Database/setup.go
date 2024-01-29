package Database

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"

)

func createIndex(DB *sql.DB, tableName, columnName string) error {
	indexName := fmt.Sprintf("idx_%s_%s", tableName, columnName)
	query := fmt.Sprintf("CREATE UNIQUE INDEX IF NOT EXISTS %s ON %s (%s)", indexName, tableName, columnName)

	_, err := DB.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating index on %s column in %s table: %w", columnName, tableName, err)
	}
	return nil
}

func CreateTables(DB *sql.DB) error {
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS GoIris (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		Plaintext TEXT,
		MD4 CHAR(32),
		MD5 CHAR(32),
		SHA1 CHAR(40),
		SHA256 CHAR(64),
		SHA512 CHAR(128)
	);
	`

	_, err := DB.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("error creating GoIris table: %w", err)
	}

	indexColumns := map[string][]string{
		"GoIris": {"MD4", "MD5", "SHA1", "SHA256", "SHA512"},
	}

	for table, columns := range indexColumns {
		for _, column := range columns {
			if err := createIndex(DB, table, column); err != nil {
				return err
			}
		}
	}

	return nil
}
