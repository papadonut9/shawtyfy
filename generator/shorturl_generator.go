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
		fmt.Sprintf(err.Error())
		os.Exit(1)
	}
	return string(encoded)
}
