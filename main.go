package main

import (
	"GoIris/Models"
	"GoIris/Models/Database"
	"GoIris/Models/Hashes"
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sqweek/dialog"
	_ "modernc.org/sqlite"

)

var TotalInserts int
var StartTime time.Time

func main() {
	Models.SetupSignalHandler()
	Models.ClearConsole()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("\033]0;%s\007", "GoIris ~ GoLang Rainbow Table")
		choice := promptUser(reader)
		switch choice {
		case "1":
			Models.ClearConsole()
			Database.SetPragmasForInsert(Models.DB)
			go Models.UpdateInsertConsoleTitle()
			insertFile()
		case "2":
			Models.ClearConsole()
			Database.SetPragmasForRead(Models.DB)
			dehashFile()
		case "3":
			Models.ClearConsole()
			Database.SetPragmasForRead(Models.DB)
			lookupHash(reader)
		case "4":
			fmt.Println("Exiting...")
			return
		}
	}
}
func promptUser(reader *bufio.Reader) string {
	fmt.Println("What would you like to do?")
	fmt.Println("1. Insert a file")
	fmt.Println("2. Dehash a File")
	fmt.Println("3. Lookup a hash")
	fmt.Println("4. Exit")

	fmt.Print("Enter choice: ")
	choice, _ := reader.ReadString('\n')
	return strings.TrimSpace(choice)
}

func insertFile() error {
	filePath, err := dialog.File().Load()
	if err != nil {
		fmt.Println("Error opening file dialog:", err)
		return err
	}

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return err
	}
	defer file.Close()

	Database.StartTime = time.Now()

	scanner := bufio.NewScanner(file)
	var lines []string
	batchSize := 512
	hashOptions := Hashes.HashOptions{
		IncludeMD4:    true,
		IncludeMD5:    true,
		IncludeSHA1:   true,
		IncludeSHA256: true,
		IncludeSHA512: true,
	}

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		if len(lines) == batchSize {
			if err := Database.InsertBatchString(Models.DB, lines, hashOptions); err != nil {
				fmt.Println("Error inserting batch data:", err)
				return err
			}
			lines = []string{}
		}
	}

	if len(lines) > 0 {
		if err := Database.InsertBatchString(Models.DB, lines, hashOptions); err != nil {
			fmt.Println("Error inserting final batch data:", err)
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning file:", err)
		return err
	}

	return nil
}

func dehashFile() error {
	filePath, err := dialog.File().Load()
	if err != nil {
		fmt.Println("Error opening file dialog:", err)
		return err
	}

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	batchSize := 512
	var batch []string

	outputDir := "Output"
	os.MkdirAll(outputDir, os.ModePerm)
	outputFile := fmt.Sprintf("%s/File%s.txt", outputDir, time.Now().Format("20060102"))
	outFile, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return err
	}
	defer outFile.Close()

	item := ""
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			fmt.Println("Invalid line format:", line)
			continue
		}
		item, hash := parts[0], parts[1]

		batch = append(batch, hash)
		if len(batch) >= batchSize {
			dehashedResults, err := Database.BatchLookup(Models.DB, batch)
			if err != nil {
				fmt.Println("Error dehashing batch:", err)
				return err
			}

			for _, h := range batch {
				plaintext := dehashedResults[h]
				fmt.Fprintf(outFile, "%s:%s: %s\n", item, h, plaintext)
			}
			batch = []string{}
		}
	}

	if len(batch) > 0 {
		dehashedResults, err := Database.BatchLookup(Models.DB, batch)
		if err != nil {
			fmt.Println("Error dehashing final batch:", err)
			return err
		}
		for _, h := range batch {
			plaintext := dehashedResults[h]
			fmt.Fprintf(outFile, "%s:%s: %s\n", item, h, plaintext)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return err
	}

	fmt.Println("Dehashing completed. Results saved to", outputFile)
	return nil
}

func lookupHash(reader *bufio.Reader) (map[string]string, time.Duration, error) {
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
