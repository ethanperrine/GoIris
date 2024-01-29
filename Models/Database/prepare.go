package Database

import (
	"GoIris/Models/Hashes"
	"strings"

)

func prepareBaseInsertQuery(options Hashes.HashOptions) (string, string) {
	queryColumns := []string{"Plaintext"}
	valuePlaceholders := []string{"?"}

	hashOptionsToColumns := map[string]struct {
		option     bool
		columnName string
	}{
		"MD4":    {options.IncludeMD4, "MD4"},
		"MD5":    {options.IncludeMD5, "MD5"},
		"SHA1":   {options.IncludeSHA1, "SHA1"},
		"SHA256": {options.IncludeSHA256, "SHA256"},
		"SHA512": {options.IncludeSHA512, "SHA512"},
	}

	orderedHashTypes := []string{"MD4", "MD5", "SHA1", "SHA256", "SHA512"}

	for _, hashType := range orderedHashTypes {
		if hashOption, ok := hashOptionsToColumns[hashType]; ok && hashOption.option {
			queryColumns = append(queryColumns, hashOption.columnName)
			valuePlaceholders = append(valuePlaceholders, "?")
		}
	}

	query := "INSERT OR IGNORE INTO GoIris (" + strings.Join(queryColumns, ", ") + ")"
	values := "VALUES (" + strings.Join(valuePlaceholders, ", ") + ")"

	return query, values
}

