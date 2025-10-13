package passcode_test

import (
	"crypto/rand"
	"crypto/sha256"
	"testing"

	"github.com/snowmerak/passcode"
	"lukechampine.com/blake3"
)

func Test_BLAKE3Hash(t *testing.T) {
	key := make([]byte, 64)
	rand.Read(key)
	hashedKey := sha256.Sum256(key)
	hasher := blake3.New(4, hashedKey[:])
	data := make([]byte, 1024)
	rand.Read(data)
	hasher.Write(data)
	firstSum := hasher.Sum(nil)
	hasher.Reset()
	hasher.Write(data)
	secondSum := hasher.Sum(nil)
	if string(firstSum) != string(secondSum) {
		t.Fatal("BLAKE3 hash mismatch")
	}

	t.Logf("BLAKE3 hash: %x (%d)", firstSum, len(firstSum))
}

func Test_SHA3Hash(t *testing.T) {
	key := make([]byte, 64)
	rand.Read(key)
	data := make([]byte, 1024)
	rand.Read(data)
	result1 := passcode.SHA3KMAC256(key, []byte("test"), data, 4)
	result2 := passcode.SHA3KMAC256(key, []byte("test"), data, 4)
	if string(result1) != string(result2) {
		t.Fatal("SHA3 hash mismatch")
	}

	t.Logf("SHA3 hash: %x (%d)", result1, len(result1))
}
