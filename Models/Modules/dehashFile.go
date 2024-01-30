package Modules

import (
	"GoIris/Models"
	"GoIris/Models/Database"
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sqweek/dialog"
)

type ItemHashPair struct {
	Item string
	Hash string
}

func DehashFile() error {
	file, err := openInputFile()
	if err != nil {
		return err
	}
	defer file.Close()

	crackedOutFile, nonCrackedOutFile, err := createOutputFiles()
	if err != nil {
		return err
	}
	defer crackedOutFile.Close()
	defer nonCrackedOutFile.Close()

	return processDehashInputFile(file, crackedOutFile, nonCrackedOutFile)
}

func openInputFile() (*os.File, error) {
	filePath, err := dialog.File().Load()
	if err != nil {
		return nil, fmt.Errorf("error opening file dialog: %w", err)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}

	return file, nil
}

func createOutputFiles() (*os.File, *os.File, error) {
	outputDir := fmt.Sprintf("Output/%s", time.Now().Format("2006-01-02 [15 04 05]"))

	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating output directory: %w", err)
	}

	crackedFile := fmt.Sprintf("%s/Cracked.txt", outputDir)
	crackedOutFile, err := os.Create(crackedFile)
	if err != nil {
		return nil, nil, fmt.Errorf("error creating cracked output file: %w", err)
	}

	nonCrackedFile := fmt.Sprintf("%s/NonCracked.txt", outputDir)
	nonCrackedOutFile, err := os.Create(nonCrackedFile)
	if err != nil {
		crackedOutFile.Close()
		return nil, nil, fmt.Errorf("error creating non-cracked output file: %w", err)
	}

	return crackedOutFile, nonCrackedOutFile, nil
}

func processDehashInputFile(file *os.File, crackedOutFile, nonCrackedOutFile *os.File) error {
	scanner := bufio.NewScanner(file)
	const batchSize = 512
	var batch []ItemHashPair

	for scanner.Scan() {
		line := scanner.Text()
		pair, err := parseLine(line)
		if err != nil {
			fmt.Println(err)
			continue
		}

		batch = append(batch, pair)
		if len(batch) >= batchSize {
			if err := processBatch(batch, crackedOutFile, nonCrackedOutFile); err != nil {
				return err
			}
			batch = []ItemHashPair{}
		}
	}

	if len(batch) > 0 {
		if err := processBatch(batch, crackedOutFile, nonCrackedOutFile); err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	return nil
}

func parseLine(line string) (ItemHashPair, error) {
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return ItemHashPair{}, fmt.Errorf("invalid line format: %s", line)
	}
	return ItemHashPair{Item: parts[0], Hash: parts[1]}, nil
}

func processBatch(batch []ItemHashPair, crackedOutFile, nonCrackedOutFile *os.File) error {
	var hashBatch []string
	for _, pair := range batch {
		hashBatch = append(hashBatch, pair.Hash)
	}

	dehashedResults, _ := Database.BatchLookup(Models.DB, hashBatch)

	for _, pair := range batch {
		plaintext, found := dehashedResults[pair.Hash]
		if found && plaintext != "" {
			if _, err := fmt.Fprintf(crackedOutFile, "%s:%s:%s\n", pair.Item, pair.Hash, plaintext); err != nil {
				fmt.Printf("Error writing to cracked output file: %v\n", err)
			}
		} else {
			if _, err := fmt.Fprintf(nonCrackedOutFile, "%s:%s\n", pair.Item, pair.Hash); err != nil {
				fmt.Printf("Error writing to non-cracked output file: %v\n", err)
			}
		}
	}

	return nil
}
