package main

import (
	"encoding/hex"
	"fmt"

	"github.com/snowmerak/passcode"
)

func main() {
	fmt.Println("=== Go Implementation Test ===\n")

	// Fixed test vectors
	key, _ := hex.DecodeString("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	challenge, _ := hex.DecodeString("fedcba9876543210fedcba9876543210")

	fmt.Printf("Key:       %s\n", hex.EncodeToString(key))
	fmt.Printf("Challenge: %s\n\n", hex.EncodeToString(challenge))

	algorithms := []struct {
		name string
		algo passcode.Algorithm
	}{
		{"SHA3-KMAC-128", passcode.AlgorithmSHA3KMAC128},
		{"SHA3-KMAC-256", passcode.AlgorithmSHA3KMAC256},
		{"BLAKE3-Keyed-128", passcode.AlgorithmBLAKE3KeyedMode128},
		{"BLAKE3-Keyed-256", passcode.AlgorithmBLAKE3KeyedMode256},
	}

	for _, test := range algorithms {
		pc, err := passcode.NewPasscode(test.algo, key)
		if err != nil {
			panic(err)
		}
		otp := pc.Compute(challenge)
		fmt.Printf("%-20s: %s\n", test.name, otp)
	}
}
