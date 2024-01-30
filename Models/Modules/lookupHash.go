package Modules

import (
	"GoIris/Models"
	"GoIris/Models/Database"
	"bufio"
	"fmt"
	"strings"
	"time"

)

func LookupHash(reader *bufio.Reader) (map[string]string, time.Duration, error) {
	fmt.Print("Enter the hash to lookup: ")
	hash, err := reader.ReadString('\n')
	if err != nil {
		return nil, 0, fmt.Errorf("error reading input: %w", err)
	}
	hash = strings.TrimSpace(hash)

	startTime := time.Now()
	result, err := Database.Lookup(Models.DB, hash)
	elapsedTime := time.Since(startTime)

	if err != nil {
		return nil, elapsedTime, fmt.Errorf("error in lookup: %w", err)
	}

	fmt.Println("Lookup Result:", result)
	fmt.Printf("Elapsed Time: %v\n", elapsedTime)

	return result, elapsedTime, nil
}
