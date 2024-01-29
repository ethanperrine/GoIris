package Hashes

import (
	"crypto"
	"errors"
	"fmt"
	"regexp"
)

var CompiledRegexMap = make(map[string]*regexp.Regexp)

func DetectHashType(data string) (string, error) {
	for hashType, regex := range CompiledRegexMap {
		if regex.MatchString(data) {
			return hashType, nil
		}
	}

	return "", errors.New("unknown hash type")
}

func CompileRegex() {
	regexMap := map[string]string{
		"MD5":    `^[a-f0-9]{32}$`,
		"SHA1":   `^[a-f0-9]{40}$`,
		"SHA256": `^[a-f0-9]{64}$`,
		"SHA512": `^[a-f0-9]{128}$`,
	}

	for hashType, pattern := range regexMap {
		compiled, err := regexp.Compile(pattern)
		if err != nil {
			panic(fmt.Sprintf("Failed to compile regex for %s: %v", hashType, err))
		}
		CompiledRegexMap[hashType] = compiled
	}
}

func ShouldIncludeHash(hashType crypto.Hash, options HashOptions) bool {
	switch hashType {
	case crypto.MD4:
		return options.IncludeMD4
	case crypto.MD5:
		return options.IncludeMD5
	case crypto.SHA1:
		return options.IncludeSHA1
	case crypto.SHA256:
		return options.IncludeSHA256
	case crypto.SHA512:
		return options.IncludeSHA512
	default:
		return false
	}
}
