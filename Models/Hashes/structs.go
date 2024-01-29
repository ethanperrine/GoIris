package Hashes

import (
	"crypto"
)

type HashFunction func(string) string

type HashOptions struct {
	IncludeMD4    bool
	IncludeMD5    bool
	IncludeSHA1   bool
	IncludeSHA256 bool
	IncludeSHA512 bool
}

var HashFunctions = map[crypto.Hash]HashFunction{
	crypto.MD4:    GetMD4Hash,
	crypto.MD5:    GetMD5Hash,
	crypto.SHA1:   GetSHA1Hash,
	crypto.SHA256: GetSHA256Hash,
	crypto.SHA512: GetSHA512Hash,
}
