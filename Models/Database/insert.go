package Database

import (
	// "GoIris/Models"
	"GoIris/Models/Hashes"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func InsertSingleString(DB *sql.DB, data string, options Hashes.HashOptions) error {
	query, values := prepareBaseInsertQuery(options)
	args := []interface{}{data}

	if options.IncludeMD4 {
		args = append(args, Hashes.GetMD4Hash(data))
	}
	if options.IncludeMD5 {
		args = append(args, Hashes.GetMD5Hash(data))
	}
	if options.IncludeSHA1 {
		args = append(args, Hashes.GetSHA1Hash(data))
	}
	if options.IncludeSHA256 {
		args = append(args, Hashes.GetSHA256Hash(data))
	}
	if options.IncludeSHA512 {
		args = append(args, Hashes.GetSHA512Hash(data))
	}

	_, err := DB.Exec(query+" "+values, args...)
	if err != nil {
		return fmt.Errorf("error executing insert query: %w", err)
	}

	return nil
}

func InsertBatchString(DB *sql.DB, data []string, options Hashes.HashOptions) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	query, values := prepareBaseInsertQuery(options)
	stmt, err := tx.Prepare(query + " " + values)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, plaintext := range data {
		args := []interface{}{plaintext}

		if options.IncludeMD4 {
			args = append(args, Hashes.GetMD4Hash(plaintext))
		}
		if options.IncludeMD5 {
			args = append(args, Hashes.GetMD5Hash(plaintext))
		}
		if options.IncludeSHA1 {
			args = append(args, Hashes.GetSHA1Hash(plaintext))
		}
		if options.IncludeSHA256 {
			args = append(args, Hashes.GetSHA256Hash(plaintext))
		}
		if options.IncludeSHA512 {
			args = append(args, Hashes.GetSHA512Hash(plaintext))
		}

		_, err = stmt.Exec(args...)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	Mutex.Lock()
	TotalInserts += len(data)
	Mutex.Unlock()
	return tx.Commit()
}
