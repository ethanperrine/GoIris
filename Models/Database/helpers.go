package Database

import (
	"database/sql"
	"fmt"
	"log"

)

func CheckRowsAffected(result sql.Result) (bool, string) {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Sprintf("Error checking rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return false, "Operation failed: No rows affected."
	}
	return true, fmt.Sprintf("Operation successful: %d rows affected.", rowsAffected)
}

func setPragmas(db *sql.DB, pragmas []string, logErrors bool) error {
	for _, pragma := range pragmas {
		stmt, err := db.Prepare(pragma)
		if err != nil {
			if logErrors {
				log.Fatalf("Failed to prepare %s: %v", pragma, err)
			} else {
				return fmt.Errorf("failed to prepare PRAGMA %s: %w", pragma, err)
			}
		}

		_, err = stmt.Exec()
		if err != nil {
			if logErrors {
				log.Fatalf("Failed to execute %s: %v", pragma, err)
			} else {
				return fmt.Errorf("failed to execute PRAGMA %s: %w", pragma, err)
			}
		}

		stmt.Close()
	}

	return nil
}

// Shoutout To this project:
// https://github.com/onthegomap/planetiler/blob/db0ab02263baaae3038de058c4fb6a1bebd81e3c/planetiler-core/src/main/java/com/onthegomap/planetiler/mbtiles/Mbtiles.java#L151-L183

func SetPragmasForInsert(db *sql.DB) {
	insertPragmas := []string{
		"PRAGMA synchronous = OFF",  // faster writes at the cost of potential data loss in case of a crash.
		"PRAGMA journal_mode = OFF", // Disables the rollback journal, which can speed up write operations but at the risk of database corruption in the event of a crash.
		// "PRAGMA synchronous = NORMAL",
		// "PRAGMA journal_mode = WAL",
		"PRAGMA cache_size = -4000000",         // Sets the cache size to 1,000,000 pages, approximately 1GB (negative value indicates the number of pages).
		"PRAGMA temp_store = MEMORY",           // Use memory for temporary storage
		"PRAGMA locking_mode = EXCLUSIVE",      // Set locking mode to EXCLUSIVE
		"PRAGMA journal_size_limit = 67108864", // 64MB
		"PRAGMA mmap_size = 30000000000",
		"PRAGMA busy_timeout = 5000", // Set busy timeout to 5 seconds
		"PRAGMA auto_vacuum = FULL",
	}
	setPragmas(db, insertPragmas, true)
}

func SetPragmasForRead(db *sql.DB) {
	readPragmas := []string{
		"PRAGMA cache_size = -100000",          // Sets the cache size to 100,000 pages (negative value indicates the number of pages).
		"PRAGMA page_size = 32768",             // Set page size to 32,768 bytes
		"PRAGMA locking_mode = EXCLUSIVE",      // Set locking mode to EXCLUSIVE
		"PRAGMA temp_store = MEMORY",           // Use memory for temporary storage
		"PRAGMA journal_size_limit = 67108864", // 64MB
		"PRAGMA journal_mode = WAL",            // Sets the journal mode to Write-Ahead Logging.
		"PRAGMA busy_timeout = 5000",           // Set busy timeout to 5 seconds
		"PRAGMA synchronous = NORMAL",          // Set synchronous mode to NORMAL
		"PRAGMA threads = 2",                   // Set number of threads for concurrent access
	}
	setPragmas(db, readPragmas, true)
}

func SetPragmasForExit(db *sql.DB) {
	exitPragmas := []string{
		"PRAGMA synchronous = NORMAL",
		"PRAGMA journal_mode = WAL",
		"PRAGMA auto_vacuum = FULL",
		"PRAGMA analysis_limit=400",
		"PRAGMA optimize",
		"PRAGMA busy_timeout = 5000",
	}

	// mang no body wants these lil ahh errors RAHH
	setPragmas(db, exitPragmas, false)

}
