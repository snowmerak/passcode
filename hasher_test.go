package passcode_test

import (
	"bytes"
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

// Test BLAKE3 KeyedMode with 256-bit output
func Test_BLAKE3KeyedMode256(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)
	data := []byte("test data for blake3 keyed mode 256")

	// Verify consistency by hashing twice with same key and data
	result1 := passcode.BLAKE3KeyedMode256(key, data)
	result2 := passcode.BLAKE3KeyedMode256(key, data)

	if !bytes.Equal(result1, result2) {
		t.Fatal("BLAKE3KeyedMode256: inconsistent results")
	}

	if len(result1) != 32 {
		t.Fatalf("BLAKE3KeyedMode256: expected 32 bytes, got %d", len(result1))
	}

	t.Logf("BLAKE3KeyedMode256: %x", result1)
}

// Test BLAKE3 KeyedMode with 512-bit output
func Test_BLAKE3KeyedMode512(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)
	data := []byte("test data for blake3 keyed mode 512")

	result1 := passcode.BLAKE3KeyedMode512(key, data)
	result2 := passcode.BLAKE3KeyedMode512(key, data)

	if !bytes.Equal(result1, result2) {
		t.Fatal("BLAKE3KeyedMode512: inconsistent results")
	}

	if len(result1) != 64 {
		t.Fatalf("BLAKE3KeyedMode512: expected 64 bytes, got %d", len(result1))
	}

	t.Logf("BLAKE3KeyedMode512: %x", result1)
}

// Test that different keys produce different results
func Test_BLAKE3KeyedMode_DifferentKeys(t *testing.T) {
	key1 := make([]byte, 32)
	key2 := make([]byte, 32)
	rand.Read(key1)
	rand.Read(key2)
	data := []byte("same data")

	result1 := passcode.BLAKE3KeyedMode256(key1, data)
	result2 := passcode.BLAKE3KeyedMode256(key2, data)

	if bytes.Equal(result1, result2) {
		t.Fatal("BLAKE3KeyedMode256: different keys produced same hash")
	}

	t.Logf("Key1 result: %x", result1)
	t.Logf("Key2 result: %x", result2)
}

// Test that different data produces different results
func Test_BLAKE3KeyedMode_DifferentData(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)
	data1 := []byte("data 1")
	data2 := []byte("data 2")

	result1 := passcode.BLAKE3KeyedMode256(key, data1)
	result2 := passcode.BLAKE3KeyedMode256(key, data2)

	if bytes.Equal(result1, result2) {
		t.Fatal("BLAKE3KeyedMode256: different data produced same hash")
	}

	t.Logf("Data1 result: %x", result1)
	t.Logf("Data2 result: %x", result2)
}

// Test empty data handling
func Test_BLAKE3KeyedMode_EmptyData(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)
	data := []byte{}

	result := passcode.BLAKE3KeyedMode256(key, data)

	if len(result) != 32 {
		t.Fatalf("BLAKE3KeyedMode256: expected 32 bytes for empty data, got %d", len(result))
	}

	t.Logf("Empty data result: %x", result)
}

// Test SHA3 KMAC128
func Test_SHA3KMAC128(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)
	customization := []byte("test-customization")
	data := []byte("test data for sha3 kmac128")

	result1 := passcode.SHA3KMAC128(key, customization, data, 32)
	result2 := passcode.SHA3KMAC128(key, customization, data, 32)

	if !bytes.Equal(result1, result2) {
		t.Fatal("SHA3KMAC128: inconsistent results")
	}

	if len(result1) != 32 {
		t.Fatalf("SHA3KMAC128: expected 32 bytes, got %d", len(result1))
	}

	t.Logf("SHA3KMAC128: %x", result1)
}

// Test SHA3 KMAC256 with various output lengths
func Test_SHA3KMAC256_Extended(t *testing.T) {
	key := make([]byte, 64)
	rand.Read(key)
	customization := []byte("my-app-v1")
	data := []byte("test data for sha3 kmac256")

	// Test various output lengths
	lengths := []int{16, 32, 64, 128}
	for _, length := range lengths {
		result := passcode.SHA3KMAC256(key, customization, data, length)
		if len(result) != length {
			t.Fatalf("SHA3KMAC256: expected %d bytes, got %d", length, len(result))
		}
		t.Logf("SHA3KMAC256 (%d bytes): %x", length, result)
	}
}

// Test that different keys produce different results
func Test_SHA3KMAC256_DifferentKeys(t *testing.T) {
	key1 := make([]byte, 32)
	key2 := make([]byte, 32)
	rand.Read(key1)
	rand.Read(key2)
	customization := []byte("test")
	data := []byte("same data")

	result1 := passcode.SHA3KMAC256(key1, customization, data, 32)
	result2 := passcode.SHA3KMAC256(key2, customization, data, 32)

	if bytes.Equal(result1, result2) {
		t.Fatal("SHA3KMAC256: different keys produced same hash")
	}

	t.Logf("Key1 result: %x", result1)
	t.Logf("Key2 result: %x", result2)
}

// Test that different customizations produce different results
func Test_SHA3KMAC256_DifferentCustomization(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)
	customization1 := []byte("app-v1")
	customization2 := []byte("app-v2")
	data := []byte("same data")

	result1 := passcode.SHA3KMAC256(key, customization1, data, 32)
	result2 := passcode.SHA3KMAC256(key, customization2, data, 32)

	if bytes.Equal(result1, result2) {
		t.Fatal("SHA3KMAC256: different customizations produced same hash")
	}

	t.Logf("Customization1 result: %x", result1)
	t.Logf("Customization2 result: %x", result2)
}

// Test that different data produces different results
func Test_SHA3KMAC256_DifferentData(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)
	customization := []byte("test")
	data1 := []byte("data 1")
	data2 := []byte("data 2")

	result1 := passcode.SHA3KMAC256(key, customization, data1, 32)
	result2 := passcode.SHA3KMAC256(key, customization, data2, 32)

	if bytes.Equal(result1, result2) {
		t.Fatal("SHA3KMAC256: different data produced same hash")
	}

	t.Logf("Data1 result: %x", result1)
	t.Logf("Data2 result: %x", result2)
}

// Test empty data handling
func Test_SHA3KMAC256_EmptyData(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)
	customization := []byte("test")
	data := []byte{}

	result := passcode.SHA3KMAC256(key, customization, data, 32)

	if len(result) != 32 {
		t.Fatalf("SHA3KMAC256: expected 32 bytes for empty data, got %d", len(result))
	}

	t.Logf("Empty data result: %x", result)
}

// Test empty customization handling
func Test_SHA3KMAC256_EmptyCustomization(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)
	customization := []byte{}
	data := []byte("test data")

	result := passcode.SHA3KMAC256(key, customization, data, 32)

	if len(result) != 32 {
		t.Fatalf("SHA3KMAC256: expected 32 bytes with empty customization, got %d", len(result))
	}

	t.Logf("Empty customization result: %x", result)
}

// Test large data (1MB)
func Test_BLAKE3_SHA3_LargeData(t *testing.T) {
	key := make([]byte, 64)
	rand.Read(key)

	// 1MB data
	largeData := make([]byte, 1024*1024)
	rand.Read(largeData)

	// BLAKE3
	blake3Result := passcode.BLAKE3KeyedMode256(key, largeData)
	if len(blake3Result) != 32 {
		t.Fatalf("BLAKE3: expected 32 bytes for large data, got %d", len(blake3Result))
	}
	t.Logf("BLAKE3 large data result: %x", blake3Result[:16])

	// SHA3 KMAC256
	sha3Result := passcode.SHA3KMAC256(key, []byte("large-data-test"), largeData, 32)
	if len(sha3Result) != 32 {
		t.Fatalf("SHA3KMAC256: expected 32 bytes for large data, got %d", len(sha3Result))
	}
	t.Logf("SHA3KMAC256 large data result: %x", sha3Result[:16])
}

// Benchmark: BLAKE3 256-bit
func Benchmark_BLAKE3KeyedMode256(b *testing.B) {
	key := make([]byte, 32)
	rand.Read(key)
	data := make([]byte, 1024)
	rand.Read(data)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		passcode.BLAKE3KeyedMode256(key, data)
	}
}

// Benchmark: BLAKE3 512-bit
func Benchmark_BLAKE3KeyedMode512(b *testing.B) {
	key := make([]byte, 32)
	rand.Read(key)
	data := make([]byte, 1024)
	rand.Read(data)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		passcode.BLAKE3KeyedMode512(key, data)
	}
}

// Benchmark: SHA3 KMAC128
func Benchmark_SHA3KMAC128(b *testing.B) {
	key := make([]byte, 32)
	rand.Read(key)
	customization := []byte("bench")
	data := make([]byte, 1024)
	rand.Read(data)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		passcode.SHA3KMAC128(key, customization, data, 32)
	}
}

// Benchmark: SHA3 KMAC256
func Benchmark_SHA3KMAC256(b *testing.B) {
	key := make([]byte, 32)
	rand.Read(key)
	customization := []byte("bench")
	data := make([]byte, 1024)
	rand.Read(data)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		passcode.SHA3KMAC256(key, customization, data, 32)
	}
}

// ========== Additional Test Cases ==========

// Test BLAKE3 with standard 32-byte key
func Test_BLAKE3KeyedMode_StandardKey(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)
	data := []byte("test data with standard 32-byte key")

	result := passcode.BLAKE3KeyedMode256(key, data)
	if len(result) != 32 {
		t.Fatalf("BLAKE3KeyedMode256: expected 32 bytes with 32-byte key, got %d", len(result))
	}

	t.Logf("Standard 32-byte key result: %x", result)
}

// Test comparison between 256-bit and 512-bit outputs
func Test_BLAKE3KeyedMode_256vs512(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)
	data := []byte("same data for both")

	result256 := passcode.BLAKE3KeyedMode256(key, data)
	result512 := passcode.BLAKE3KeyedMode512(key, data)

	// The first 32 bytes of 512-bit result may differ from 256-bit result
	if len(result256) != 32 {
		t.Fatalf("BLAKE3KeyedMode256: expected 32 bytes, got %d", len(result256))
	}
	if len(result512) != 64 {
		t.Fatalf("BLAKE3KeyedMode512: expected 64 bytes, got %d", len(result512))
	}

	t.Logf("256-bit result: %x", result256)
	t.Logf("512-bit result: %x", result512)
}

// Test multiple rounds of hashing
func Test_BLAKE3KeyedMode_MultipleRounds(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)
	data := []byte("initial data")

	// First round
	result1 := passcode.BLAKE3KeyedMode256(key, data)

	// Use result as key
	result2 := passcode.BLAKE3KeyedMode256(result1, data)

	// Use result as key again
	result3 := passcode.BLAKE3KeyedMode256(result2, data)

	// All results should be different
	if bytes.Equal(result1, result2) {
		t.Fatal("BLAKE3KeyedMode256: round 1 and 2 produced same hash")
	}
	if bytes.Equal(result2, result3) {
		t.Fatal("BLAKE3KeyedMode256: round 2 and 3 produced same hash")
	}
	if bytes.Equal(result1, result3) {
		t.Fatal("BLAKE3KeyedMode256: round 1 and 3 produced same hash")
	}

	t.Logf("Round 1: %x", result1)
	t.Logf("Round 2: %x", result2)
	t.Logf("Round 3: %x", result3)
}

// Test SHA3 KMAC with short key
func Test_SHA3KMAC_ShortKey(t *testing.T) {
	key := []byte("short")
	customization := []byte("test")
	data := []byte("test data")

	result128 := passcode.SHA3KMAC128(key, customization, data, 32)
	result256 := passcode.SHA3KMAC256(key, customization, data, 32)

	if len(result128) != 32 {
		t.Fatalf("SHA3KMAC128: expected 32 bytes with short key, got %d", len(result128))
	}
	if len(result256) != 32 {
		t.Fatalf("SHA3KMAC256: expected 32 bytes with short key, got %d", len(result256))
	}

	t.Logf("KMAC128 short key: %x", result128)
	t.Logf("KMAC256 short key: %x", result256)
}

// Test SHA3 KMAC with long key
func Test_SHA3KMAC_LongKey(t *testing.T) {
	key := make([]byte, 512)
	rand.Read(key)
	customization := []byte("test")
	data := []byte("test data")

	result128 := passcode.SHA3KMAC128(key, customization, data, 32)
	result256 := passcode.SHA3KMAC256(key, customization, data, 32)

	if len(result128) != 32 {
		t.Fatalf("SHA3KMAC128: expected 32 bytes with long key, got %d", len(result128))
	}
	if len(result256) != 32 {
		t.Fatalf("SHA3KMAC256: expected 32 bytes with long key, got %d", len(result256))
	}

	t.Logf("KMAC128 long key: %x", result128)
	t.Logf("KMAC256 long key: %x", result256)
}

// Test KMAC128 vs KMAC256 with same input produces different results
func Test_SHA3KMAC_128vs256(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)
	customization := []byte("test")
	data := []byte("same data for both")

	result128 := passcode.SHA3KMAC128(key, customization, data, 32)
	result256 := passcode.SHA3KMAC256(key, customization, data, 32)

	if bytes.Equal(result128, result256) {
		t.Fatal("SHA3KMAC: KMAC128 and KMAC256 produced same hash")
	}

	t.Logf("KMAC128 result: %x", result128)
	t.Logf("KMAC256 result: %x", result256)
}

// Test multiple rounds of SHA3 KMAC hashing
func Test_SHA3KMAC_MultipleRounds(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)
	customization := []byte("multi-round")
	data := []byte("initial data")

	// First round
	result1 := passcode.SHA3KMAC256(key, customization, data, 32)

	// Use result as key
	result2 := passcode.SHA3KMAC256(result1, customization, data, 32)

	// Use result as key again
	result3 := passcode.SHA3KMAC256(result2, customization, data, 32)

	// All results should be different
	if bytes.Equal(result1, result2) {
		t.Fatal("SHA3KMAC256: round 1 and 2 produced same hash")
	}
	if bytes.Equal(result2, result3) {
		t.Fatal("SHA3KMAC256: round 2 and 3 produced same hash")
	}
	if bytes.Equal(result1, result3) {
		t.Fatal("SHA3KMAC256: round 1 and 3 produced same hash")
	}

	t.Logf("Round 1: %x", result1)
	t.Logf("Round 2: %x", result2)
	t.Logf("Round 3: %x", result3)
}

// Test SHA3 KMAC128 with variable output lengths
func Test_SHA3KMAC128_VariableOutputLength(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)
	customization := []byte("var-len")
	data := []byte("test data")

	lengths := []int{8, 16, 32, 64, 128, 256}
	for _, length := range lengths {
		result := passcode.SHA3KMAC128(key, customization, data, length)
		if len(result) != length {
			t.Fatalf("SHA3KMAC128: expected %d bytes, got %d", length, len(result))
		}
		t.Logf("KMAC128 (%d bytes): %x", length, result[:min(16, length)])
	}
}

// Test SHA3 KMAC256 with very long customization string
func Test_SHA3KMAC256_LongCustomization(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)
	customization := make([]byte, 1024)
	rand.Read(customization)
	data := []byte("test data")

	result := passcode.SHA3KMAC256(key, customization, data, 32)
	if len(result) != 32 {
		t.Fatalf("SHA3KMAC256: expected 32 bytes with long customization, got %d", len(result))
	}

	t.Logf("Long customization result: %x", result)
}

// Test comparison between BLAKE3 and SHA3 algorithms
func Test_BLAKE3_vs_SHA3_Comparison(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)
	data := []byte("comparison test data")

	blake3Result := passcode.BLAKE3KeyedMode256(key, data)
	sha3Result := passcode.SHA3KMAC256(key, []byte(""), data, 32)

	// Different algorithms should produce different results
	if bytes.Equal(blake3Result, sha3Result) {
		t.Fatal("BLAKE3 and SHA3 produced same hash (extremely unlikely)")
	}

	t.Logf("BLAKE3 result: %x", blake3Result)
	t.Logf("SHA3 result: %x", sha3Result)
}

// Test BLAKE3 with various data sizes
func Test_BLAKE3KeyedMode_VariousDataSizes(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)

	sizes := []int{1, 16, 64, 256, 1024, 4096, 16384}
	for _, size := range sizes {
		data := make([]byte, size)
		rand.Read(data)

		result := passcode.BLAKE3KeyedMode256(key, data)
		if len(result) != 32 {
			t.Fatalf("BLAKE3KeyedMode256: expected 32 bytes for %d byte data, got %d", size, len(result))
		}
		t.Logf("Data size %d bytes: %x", size, result[:16])
	}
}

// Test SHA3 KMAC256 with various data sizes
func Test_SHA3KMAC256_VariousDataSizes(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)
	customization := []byte("var-size")

	sizes := []int{1, 16, 64, 256, 1024, 4096, 16384}
	for _, size := range sizes {
		data := make([]byte, size)
		rand.Read(data)

		result := passcode.SHA3KMAC256(key, customization, data, 32)
		if len(result) != 32 {
			t.Fatalf("SHA3KMAC256: expected 32 bytes for %d byte data, got %d", size, len(result))
		}
		t.Logf("Data size %d bytes: %x", size, result[:16])
	}
}

// Benchmark: BLAKE3 with small data (64 bytes)
func Benchmark_BLAKE3KeyedMode256_SmallData(b *testing.B) {
	key := make([]byte, 32)
	rand.Read(key)
	data := make([]byte, 64)
	rand.Read(data)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		passcode.BLAKE3KeyedMode256(key, data)
	}
}

// Benchmark: BLAKE3 with medium data (1KB)
func Benchmark_BLAKE3KeyedMode256_MediumData(b *testing.B) {
	key := make([]byte, 32)
	rand.Read(key)
	data := make([]byte, 1024)
	rand.Read(data)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		passcode.BLAKE3KeyedMode256(key, data)
	}
}

// Benchmark: BLAKE3 with large data (1MB)
func Benchmark_BLAKE3KeyedMode256_LargeData(b *testing.B) {
	key := make([]byte, 32)
	rand.Read(key)
	data := make([]byte, 1024*1024)
	rand.Read(data)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		passcode.BLAKE3KeyedMode256(key, data)
	}
}

// Benchmark: SHA3 KMAC256 with small data (64 bytes)
func Benchmark_SHA3KMAC256_SmallData(b *testing.B) {
	key := make([]byte, 32)
	rand.Read(key)
	customization := []byte("bench")
	data := make([]byte, 64)
	rand.Read(data)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		passcode.SHA3KMAC256(key, customization, data, 32)
	}
}

// Benchmark: SHA3 KMAC256 with medium data (1KB)
func Benchmark_SHA3KMAC256_MediumData(b *testing.B) {
	key := make([]byte, 32)
	rand.Read(key)
	customization := []byte("bench")
	data := make([]byte, 1024)
	rand.Read(data)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		passcode.SHA3KMAC256(key, customization, data, 32)
	}
}

// Benchmark: SHA3 KMAC256 with large data (1MB)
func Benchmark_SHA3KMAC256_LargeData(b *testing.B) {
	key := make([]byte, 32)
	rand.Read(key)
	customization := []byte("bench")
	data := make([]byte, 1024*1024)
	rand.Read(data)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		passcode.SHA3KMAC256(key, customization, data, 32)
	}
}

// Benchmark: Direct comparison between BLAKE3 and SHA3
func Benchmark_Comparison_BLAKE3_vs_SHA3(b *testing.B) {
	key := make([]byte, 32)
	rand.Read(key)
	data := make([]byte, 1024)
	rand.Read(data)

	b.Run("BLAKE3-256", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			passcode.BLAKE3KeyedMode256(key, data)
		}
	})

	b.Run("SHA3-KMAC256", func(b *testing.B) {
		customization := []byte("bench")
		for i := 0; i < b.N; i++ {
			passcode.SHA3KMAC256(key, customization, data, 32)
		}
	})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
