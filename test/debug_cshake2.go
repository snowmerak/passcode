package main

import (
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/sha3"
)

func main() {
	key, _ := hex.DecodeString("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	
	// Test with empty function name and customization
	hasher := sha3.NewCShake128([]byte(""), []byte(""))
	hasher.Write(key)
	output := make([]byte, 16)
	hasher.Read(output)
	fmt.Printf("Go CShake128(\"\", \"\"): %s\n", hex.EncodeToString(output))
}
