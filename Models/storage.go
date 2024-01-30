package Models

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"os"
)

func CalculateFileHashSizes(file *os.File) (int, int, int, int, int) {
	scanner := bufio.NewScanner(file)
	lineCount := 0

	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return 0, 0, 0, 0, 0
	}

	file.Seek(0, 0)

	md4Size := lineCount * md5.Size // MD4 is roughly the same size as MD5
	md5Size := lineCount * md5.Size
	sha1Size := lineCount * sha1.Size
	sha256Size := lineCount * sha256.Size
	sha512Size := lineCount * sha512.Size

	return md4Size, md5Size, sha1Size, sha256Size, sha512Size
}
