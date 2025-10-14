package main

import (
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/sha3"
)

func main() {
	// Test how Go's CShake works
	key, _ := hex.DecodeString("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	
	// Test 1: CShake with function name "KMAC", customization "authorization"
	hasher1 := sha3.NewCShake128([]byte("KMAC"), []byte("authorization"))
	hasher1.Write(key)
	output1 := make([]byte, 16)
	hasher1.Read(output1)
	fmt.Printf("CShake128(\"KMAC\", \"authorization\"): %s\n", hex.EncodeToString(output1))
	
	// Test 2: Check what the Rust implementation expects
	// According to NIST SP 800-185, the input to Keccak should be:
	// bytepad(encode_string(N) || encode_string(S), rate) || ...
	fmt.Println("\nThis tests whether Go's NewCShake handles N and S parameters correctly")
}
