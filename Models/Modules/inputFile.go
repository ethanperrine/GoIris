package Modules

import (
	"GoIris/Models"
	"GoIris/Models/Database"
	"GoIris/Models/Hashes"
	"bufio"
	"fmt"
	"os"
	"strings"

)

func InsertFile() error {
	file, err := openInputFile()
	if err != nil {
		return err
	}
	defer file.Close()

	if err := checkDiskSpace(file); err != nil {
		return err
	}

	return processInputFile(file)
}

func checkDiskSpace(file *os.File) error {
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
			return fmt.Errorf("error reading input: %w", err)
		}
		response = strings.TrimSpace(response)

		if strings.ToLower(response) != "yes" && strings.ToLower(response) != "y" {
			fmt.Println("Operation aborted.")
			return fmt.Errorf("operation aborted by user")
		}
	}

	fmt.Println("Minimum Size Of DB:", Models.ConvertBytesToPretty(totalSize))
	return nil
}

func processInputFile(file *os.File) error {
	scanner := bufio.NewScanner(file)
	var lines []string
	const batchSize = 512
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
			if err := insertBatch(lines, hashOptions); err != nil {
				return err
			}
			lines = []string{}
		}
	}

	if len(lines) > 0 {
		if err := insertBatch(lines, hashOptions); err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error scanning file: %w", err)
	}

	return nil
}

func insertBatch(lines []string, hashOptions Hashes.HashOptions) error {
	if err := Database.InsertBatchString(Models.DB, lines, hashOptions); err != nil {
		return fmt.Errorf("error inserting batch data: %w", err)
	}
	return nil
}
