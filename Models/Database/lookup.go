package Database

import (
	"GoIris/Models/Hashes"
	"database/sql"
	"errors"
	"fmt"
	"strings"

)

func Lookup(db *sql.DB, query string) (map[string]string, error) {
	hashType, err := Hashes.DetectHashType(query)
	if err != nil {
		return nil, err
	}

	var sqlQuery string
	var columnName string
	switch hashType {
	case "MD4":
		columnName = "MD4"
	case "MD5":
		columnName = "MD5"
	case "SHA1":
		columnName = "SHA1"
	case "SHA256":
		columnName = "SHA256"
	case "SHA512":
		columnName = "SHA512"
	default:
		return nil, errors.New("unsupported hash type for lookup")
	}

	sqlQuery = fmt.Sprintf("SELECT Plaintext, %s FROM GoIris WHERE %s = ?", columnName, columnName)

	row := db.QueryRow(sqlQuery, query)
	var plaintext, hash string
	err = row.Scan(&plaintext, &hash)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("no matching record found")
		}
		return nil, err
	}

	result := map[string]string{
		hashType: plaintext,
	}

	return result, nil
}

func BatchLookup(db *sql.DB, queries []string) (map[string]string, error) {
    if len(queries) == 0 {
        return nil, errors.New("no queries provided")
    }

    categorizedQueries := make(map[string][]string)
    for _, query := range queries {
        hashType, err := Hashes.DetectHashType(query)
        if err != nil {
            return nil, err
        }
        categorizedQueries[hashType] = append(categorizedQueries[hashType], query)
    }

    results := make(map[string]string)
    for hashType, hashes := range categorizedQueries {
        var columnName string
        switch hashType {
        case "MD4":
            columnName = "MD4"
        case "MD5":
            columnName = "MD5"
        case "SHA1":
            columnName = "SHA1"
        case "SHA256":
            columnName = "SHA256"
        case "SHA512":
            columnName = "SHA512"
        default:
            continue
        }

        placeholders := strings.Repeat("?,", len(hashes))
        placeholders = strings.TrimRight(placeholders, ",")
        sqlQuery := fmt.Sprintf("SELECT Plaintext, %s FROM GoIris WHERE %s IN (%s)", columnName, columnName, placeholders)

        args := make([]interface{}, len(hashes))
        for i, q := range hashes {
            args[i] = q
        }

        rows, err := db.Query(sqlQuery, args...)
        if err != nil {
            return nil, err
        }

        for rows.Next() {
            var plaintext, hash string
            if err := rows.Scan(&plaintext, &hash); err != nil {
                rows.Close()
                return nil, err
            }
            results[hash] = plaintext
        }

        if err := rows.Err(); err != nil {
            rows.Close()
            return nil, err
        }
        rows.Close()
    }

    return results, nil
}

