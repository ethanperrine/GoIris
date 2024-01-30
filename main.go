package main

import (
	"GoIris/Models"
	"GoIris/Models/Database"
	"GoIris/Models/Hashes"
	"bufio"
	"context"
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
			ctx, cancel := context.WithCancel(context.Background())
			Models.ClearConsole()
			Database.SetPragmasForInsert(Models.DB)
			go Models.UpdateInsertConsoleTitle(ctx)

			err := insertFile()
			if err != nil {
				fmt.Println("Error opening file dialog:", err)
			}

			cancel()
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

	freeSpace := int64(Models.CheckDiskSpace())
	md4size, md5Size, sha1Size, sha256Size, sha512Size := Models.CalculateFileHashSizes(file)
	totalSize := int64(md4size + md5Size + sha1Size + sha256Size + sha512Size)

	if freeSpace < totalSize {
		fmt.Println("Warning: Not enough free space to insert file.")
		fmt.Println("Continuing may cause database corruption due to running out of space.")
		fmt.Print("Do you want to continue? (yes/no): ")

		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			return err
		}
		response = strings.TrimSpace(response)

		if strings.ToLower(response) != "yes" || strings.ToLower(response) != "y" {
			fmt.Println("Operation aborted.")
			return err
		}
	}
	fmt.Println("Minimum Size Of DB:", Models.ConvertBytesToPretty(totalSize))

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

type ItemHashPair struct {
	Item string
	Hash string
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
	var batch []ItemHashPair

	outputDir := fmt.Sprintf("Output/%s", time.Now().Format("[2006-01-02]"))
	os.MkdirAll(outputDir, os.ModePerm)
	crackedFile := fmt.Sprintf("%s/Cracked.txt", outputDir)
	nonCrackedFile := fmt.Sprintf("%s/NonCracked.txt", outputDir)

	crackedOutFile, err := os.Create(crackedFile)
	if err != nil {
		fmt.Println("Error creating cracked output file:", err)
		return err
	}
	defer crackedOutFile.Close()

	nonCrackedOutFile, err := os.Create(nonCrackedFile)
	if err != nil {
		fmt.Println("Error creating non-cracked output file:", err)
		return err
	}
	defer nonCrackedOutFile.Close()

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			fmt.Println("Invalid line format:", line)
			continue
		}
		batch = append(batch, ItemHashPair{Item: parts[0], Hash: parts[1]})
		if len(batch) >= batchSize {
			var hashBatch []string
			for _, pair := range batch {
				hashBatch = append(hashBatch, pair.Hash)
			}

			dehashedResults, err := Database.BatchLookup(Models.DB, hashBatch)
			if err != nil {
				fmt.Println("Error dehashing batch:", err)
				return err
			}

			for _, pair := range batch {
				plaintext := dehashedResults[pair.Hash]
				if plaintext != "" {
					fmt.Fprintf(crackedOutFile, "%s:%s:%s\n", pair.Item, pair.Hash, plaintext)
				} else {
					fmt.Fprintf(nonCrackedOutFile, "%s:%s\n", pair.Item, pair.Hash)
				}
			}
			batch = []ItemHashPair{}
		}
	}

	if len(batch) > 0 {
		var hashBatch []string
		for _, pair := range batch {
			hashBatch = append(hashBatch, pair.Hash)
		}

		dehashedResults, err := Database.BatchLookup(Models.DB, hashBatch)
		if err != nil {
			fmt.Println("Error dehashing final batch:", err)
			return err
		}

		for _, pair := range batch {
			plaintext := dehashedResults[pair.Hash]
			if plaintext != "" {
				fmt.Fprintf(crackedOutFile, "%s:%s:%s\n", pair.Item, pair.Hash, plaintext)
			} else {
				fmt.Fprintf(nonCrackedOutFile, "%s:%s\n", pair.Item, pair.Hash)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return err
	}
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
