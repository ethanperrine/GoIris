package Hashes

import (
    "crypto/md5"
    "crypto/sha1"
    "crypto/sha256"
    "crypto/sha512"
    "encoding/hex"

    "golang.org/x/crypto/md4"

)

func GetMD4Hash(text string) string {
    hasher := md4.New()
    hasher.Write([]byte(text))
    return hex.EncodeToString(hasher.Sum(nil))
}

func GetMD5Hash(text string) string {
    hash := md5.Sum([]byte(text))
    return hex.EncodeToString(hash[:])
}

func GetSHA1Hash(text string) string {
    hash := sha1.Sum([]byte(text))
    return hex.EncodeToString(hash[:])
}

func GetSHA256Hash(text string) string {
    hash := sha256.Sum256([]byte(text))
    return hex.EncodeToString(hash[:])
}

func GetSHA512Hash(text string) string {
    hash := sha512.Sum512([]byte(text))
    return hex.EncodeToString(hash[:])
}
