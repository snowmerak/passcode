package passcode

import (
	"crypto/sha256"

	"lukechampine.com/blake3"
)

func blake3KeyedMode(key []byte, data []byte, outLen int) []byte {
	hasehdKey := blake3.Sum256(key)
	hasher := blake3.New(outLen, hasehdKey[:])
	hasher.Write(data)
	return hasher.Sum(nil)
}

func BLAKE3KeyedMode256(key, data []byte) []byte {
	return blake3KeyedMode(key, data, 32)
}

func BLAKE3KeyedMode512(key, data []byte) []byte {
	return blake3KeyedMode(key, data, 64)
}
