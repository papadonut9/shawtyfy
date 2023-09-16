package generator

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"os"

	"github.com/itchyny/base58-go"
)

// Hashing initial input
func sha256Of(input string) []byte{
	algorithm := sha256.New()
	algorithm.Write([]byte(input))
	return algorithm.Sum(nil)
}

// Base58 encoding implemented
// Using bitcoin style encoding to make the text less ambiguous
func base58Encode(bytes []byte) string{
	encoding := base58.BitcoinEncoding
	encoded, err := encoding.Encode(bytes)
	if err != nil{
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return string(encoded)
}

// shortening logic
/*
	Logic:
		- Generate a SHA256 hash of the url and user UUID
		- Encode the hash into Base58 string
		- output the first 8 characters as the shortened link
*/
func generateShortUrl(originalUrl string, userId string) string{
	urlHashBytes := sha256Of(originalUrl + userId)
	generatedNumber := new(big.Int).SetBytes(urlHashBytes).Uint64()
	finalString := base58Encode([]byte(fmt.Sprintf("%d", generatedNumber)))
	return finalString[:8]
}