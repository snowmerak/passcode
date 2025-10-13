package passcode

import "fmt"

type Algorithm string

const (
	AlgorithmSHA3KMAC128                  = "SHA3-KMAC-128"
	AlgorithmSHA3KMAC256                  = "SHA3-KMAC-256"
	AlgorithmBLAKE3KeyedMode128 Algorithm = "BLAKE3-Keyed-Mode-128"
	AlgorithmBLAKE3KeyedMode256 Algorithm = "BLAKE3-Keyed-Mode-256"
)

type Hasher func(key []byte, data []byte) []byte

type Passcode struct {
	algorithm string
	hasher    Hasher
}

func NewPasscode(algorithm Algorithm) (*Passcode, error) {
	var hasher Hasher
	switch algorithm {
	case AlgorithmSHA3KMAC128:
		hasher = sha3KMAC256ForPasscode
	case AlgorithmSHA3KMAC256:
		hasher = sha3KMAC256ForPasscode
	case AlgorithmBLAKE3KeyedMode128:
		hasher = BLAKE3KeyedMode256 // Using 256-bit output for 128-bit mode
	case AlgorithmBLAKE3KeyedMode256:
		hasher = BLAKE3KeyedMode512
	default:
		return nil, fmt.Errorf("unknown hash algorithm: %s", algorithm)
	}

	return &Passcode{
		algorithm: string(algorithm),
		hasher:    hasher,
	}, nil
}
