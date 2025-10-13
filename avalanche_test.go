package passcode_test

import (
	"crypto/rand"
	"testing"

	"github.com/snowmerak/passcode"
)

// hammingDistance calculates the number of differing bits between two byte slices
func hammingDistance(a, b []byte) int {
	if len(a) != len(b) {
		panic("slices must have equal length")
	}

	distance := 0
	for i := 0; i < len(a); i++ {
		xor := a[i] ^ b[i]
		for xor != 0 {
			distance += int(xor & 1)
			xor >>= 1
		}
	}
	return distance
}

// bitFlipRatio calculates the percentage of bits that differ
func bitFlipRatio(a, b []byte) float64 {
	totalBits := len(a) * 8
	flippedBits := hammingDistance(a, b)
	return float64(flippedBits) / float64(totalBits) * 100.0
}

// Test_BLAKE3_AvalancheEffect_SingleBitFlip tests that flipping a single bit in input
// causes approximately 50% of output bits to change (ideal avalanche effect)
func Test_BLAKE3_AvalancheEffect_SingleBitFlip(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)

	data := make([]byte, 64)
	rand.Read(data)

	originalHash := passcode.BLAKE3KeyedMode256(key, data)

	// Test flipping each bit position
	totalFlipRatio := 0.0
	testCount := 0

	for byteIdx := 0; byteIdx < len(data); byteIdx++ {
		for bitIdx := 0; bitIdx < 8; bitIdx++ {
			// Create a copy and flip one bit
			modifiedData := make([]byte, len(data))
			copy(modifiedData, data)
			modifiedData[byteIdx] ^= (1 << bitIdx)

			modifiedHash := passcode.BLAKE3KeyedMode256(key, modifiedData)

			flipRatio := bitFlipRatio(originalHash, modifiedHash)
			totalFlipRatio += flipRatio
			testCount++

			// Each single bit flip should cause significant change (typically 40-60%)
			if flipRatio < 25.0 || flipRatio > 75.0 {
				t.Logf("Warning: Bit flip at byte %d, bit %d caused only %.2f%% output change",
					byteIdx, bitIdx, flipRatio)
			}
		}
	}

	avgFlipRatio := totalFlipRatio / float64(testCount)
	t.Logf("BLAKE3 average bit flip ratio: %.2f%% (ideal: ~50%%)", avgFlipRatio)

	// Overall average should be reasonably close to 50%
	if avgFlipRatio < 40.0 || avgFlipRatio > 60.0 {
		t.Errorf("BLAKE3 avalanche effect out of ideal range: %.2f%%", avgFlipRatio)
	}
}

// Test_SHA3KMAC256_AvalancheEffect_SingleBitFlip tests avalanche effect for SHA3 KMAC256
func Test_SHA3KMAC256_AvalancheEffect_SingleBitFlip(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)
	customization := []byte("avalanche-test")

	data := make([]byte, 64)
	rand.Read(data)

	originalHash := passcode.SHA3KMAC256(key, customization, data, 32)

	totalFlipRatio := 0.0
	testCount := 0

	for byteIdx := 0; byteIdx < len(data); byteIdx++ {
		for bitIdx := 0; bitIdx < 8; bitIdx++ {
			modifiedData := make([]byte, len(data))
			copy(modifiedData, data)
			modifiedData[byteIdx] ^= (1 << bitIdx)

			modifiedHash := passcode.SHA3KMAC256(key, customization, modifiedData, 32)

			flipRatio := bitFlipRatio(originalHash, modifiedHash)
			totalFlipRatio += flipRatio
			testCount++

			if flipRatio < 25.0 || flipRatio > 75.0 {
				t.Logf("Warning: Bit flip at byte %d, bit %d caused only %.2f%% output change",
					byteIdx, bitIdx, flipRatio)
			}
		}
	}

	avgFlipRatio := totalFlipRatio / float64(testCount)
	t.Logf("SHA3 KMAC256 average bit flip ratio: %.2f%% (ideal: ~50%%)", avgFlipRatio)

	if avgFlipRatio < 40.0 || avgFlipRatio > 60.0 {
		t.Errorf("SHA3 KMAC256 avalanche effect out of ideal range: %.2f%%", avgFlipRatio)
	}
}

// Test_BLAKE3_AvalancheEffect_KeyModification tests avalanche effect when modifying key
func Test_BLAKE3_AvalancheEffect_KeyModification(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)
	data := make([]byte, 64)
	rand.Read(data)

	originalHash := passcode.BLAKE3KeyedMode256(key, data)

	totalFlipRatio := 0.0
	testCount := 0

	// Test flipping bits in key
	for byteIdx := 0; byteIdx < len(key); byteIdx++ {
		for bitIdx := 0; bitIdx < 8; bitIdx++ {
			modifiedKey := make([]byte, len(key))
			copy(modifiedKey, key)
			modifiedKey[byteIdx] ^= (1 << bitIdx)

			modifiedHash := passcode.BLAKE3KeyedMode256(modifiedKey, data)

			flipRatio := bitFlipRatio(originalHash, modifiedHash)
			totalFlipRatio += flipRatio
			testCount++
		}
	}

	avgFlipRatio := totalFlipRatio / float64(testCount)
	t.Logf("BLAKE3 key modification average bit flip ratio: %.2f%% (ideal: ~50%%)", avgFlipRatio)

	if avgFlipRatio < 40.0 || avgFlipRatio > 60.0 {
		t.Errorf("BLAKE3 key avalanche effect out of ideal range: %.2f%%", avgFlipRatio)
	}
}

// Test_SHA3KMAC256_AvalancheEffect_KeyModification tests avalanche effect for key changes
func Test_SHA3KMAC256_AvalancheEffect_KeyModification(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)
	customization := []byte("key-avalanche-test")
	data := make([]byte, 64)
	rand.Read(data)

	originalHash := passcode.SHA3KMAC256(key, customization, data, 32)

	totalFlipRatio := 0.0
	testCount := 0

	for byteIdx := 0; byteIdx < len(key); byteIdx++ {
		for bitIdx := 0; bitIdx < 8; bitIdx++ {
			modifiedKey := make([]byte, len(key))
			copy(modifiedKey, key)
			modifiedKey[byteIdx] ^= (1 << bitIdx)

			modifiedHash := passcode.SHA3KMAC256(modifiedKey, customization, data, 32)

			flipRatio := bitFlipRatio(originalHash, modifiedHash)
			totalFlipRatio += flipRatio
			testCount++
		}
	}

	avgFlipRatio := totalFlipRatio / float64(testCount)
	t.Logf("SHA3 KMAC256 key modification average bit flip ratio: %.2f%% (ideal: ~50%%)", avgFlipRatio)

	if avgFlipRatio < 40.0 || avgFlipRatio > 60.0 {
		t.Errorf("SHA3 KMAC256 key avalanche effect out of ideal range: %.2f%%", avgFlipRatio)
	}
}

// Test_SHA3KMAC256_AvalancheEffect_CustomizationModification tests avalanche effect for customization changes
func Test_SHA3KMAC256_AvalancheEffect_CustomizationModification(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)
	customization := []byte("original-customization-string")
	data := make([]byte, 64)
	rand.Read(data)

	originalHash := passcode.SHA3KMAC256(key, customization, data, 32)

	totalFlipRatio := 0.0
	testCount := 0

	for byteIdx := 0; byteIdx < len(customization); byteIdx++ {
		for bitIdx := 0; bitIdx < 8; bitIdx++ {
			modifiedCustomization := make([]byte, len(customization))
			copy(modifiedCustomization, customization)
			modifiedCustomization[byteIdx] ^= (1 << bitIdx)

			modifiedHash := passcode.SHA3KMAC256(key, modifiedCustomization, data, 32)

			flipRatio := bitFlipRatio(originalHash, modifiedHash)
			totalFlipRatio += flipRatio
			testCount++
		}
	}

	avgFlipRatio := totalFlipRatio / float64(testCount)
	t.Logf("SHA3 KMAC256 customization modification average bit flip ratio: %.2f%% (ideal: ~50%%)", avgFlipRatio)

	if avgFlipRatio < 40.0 || avgFlipRatio > 60.0 {
		t.Errorf("SHA3 KMAC256 customization avalanche effect out of ideal range: %.2f%%", avgFlipRatio)
	}
}

// Test_BLAKE3_AvalancheEffect_ByteModification tests changing entire bytes
func Test_BLAKE3_AvalancheEffect_ByteModification(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)
	data := make([]byte, 64)
	rand.Read(data)

	originalHash := passcode.BLAKE3KeyedMode256(key, data)

	totalFlipRatio := 0.0
	testCount := 0

	// Test modifying each byte to a different random value
	for byteIdx := 0; byteIdx < len(data); byteIdx++ {
		modifiedData := make([]byte, len(data))
		copy(modifiedData, data)

		// Change to a different value
		modifiedData[byteIdx] = ^modifiedData[byteIdx]

		modifiedHash := passcode.BLAKE3KeyedMode256(key, modifiedData)

		flipRatio := bitFlipRatio(originalHash, modifiedHash)
		totalFlipRatio += flipRatio
		testCount++
	}

	avgFlipRatio := totalFlipRatio / float64(testCount)
	t.Logf("BLAKE3 byte modification average bit flip ratio: %.2f%% (ideal: ~50%%)", avgFlipRatio)

	if avgFlipRatio < 40.0 || avgFlipRatio > 60.0 {
		t.Errorf("BLAKE3 byte modification avalanche effect out of ideal range: %.2f%%", avgFlipRatio)
	}
}

// Test_SHA3KMAC256_AvalancheEffect_ByteModification tests changing entire bytes
func Test_SHA3KMAC256_AvalancheEffect_ByteModification(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)
	customization := []byte("byte-mod-test")
	data := make([]byte, 64)
	rand.Read(data)

	originalHash := passcode.SHA3KMAC256(key, customization, data, 32)

	totalFlipRatio := 0.0
	testCount := 0

	for byteIdx := 0; byteIdx < len(data); byteIdx++ {
		modifiedData := make([]byte, len(data))
		copy(modifiedData, data)
		modifiedData[byteIdx] = ^modifiedData[byteIdx]

		modifiedHash := passcode.SHA3KMAC256(key, customization, modifiedData, 32)

		flipRatio := bitFlipRatio(originalHash, modifiedHash)
		totalFlipRatio += flipRatio
		testCount++
	}

	avgFlipRatio := totalFlipRatio / float64(testCount)
	t.Logf("SHA3 KMAC256 byte modification average bit flip ratio: %.2f%% (ideal: ~50%%)", avgFlipRatio)

	if avgFlipRatio < 40.0 || avgFlipRatio > 60.0 {
		t.Errorf("SHA3 KMAC256 byte modification avalanche effect out of ideal range: %.2f%%", avgFlipRatio)
	}
}

// Test_BLAKE3_AvalancheEffect_MinimalChange tests with minimal data changes
func Test_BLAKE3_AvalancheEffect_MinimalChange(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)

	testCases := []struct {
		name  string
		data1 []byte
		data2 []byte
	}{
		{
			name:  "Single byte difference",
			data1: []byte("The quick brown fox jumps over the lazy dog"),
			data2: []byte("The quick brown fox jumps over the lazy doh"),
		},
		{
			name:  "Case change",
			data1: []byte("Hello World"),
			data2: []byte("hello World"),
		},
		{
			name:  "Extra space",
			data1: []byte("HelloWorld"),
			data2: []byte("Hello World"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hash1 := passcode.BLAKE3KeyedMode256(key, tc.data1)
			hash2 := passcode.BLAKE3KeyedMode256(key, tc.data2)

			flipRatio := bitFlipRatio(hash1, hash2)
			t.Logf("BLAKE3 %s: %.2f%% bits flipped", tc.name, flipRatio)

			if flipRatio < 25.0 {
				t.Errorf("BLAKE3 %s: insufficient avalanche effect (%.2f%%)", tc.name, flipRatio)
			}
		})
	}
}

// Test_SHA3KMAC256_AvalancheEffect_MinimalChange tests with minimal data changes
func Test_SHA3KMAC256_AvalancheEffect_MinimalChange(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)
	customization := []byte("minimal-change-test")

	testCases := []struct {
		name  string
		data1 []byte
		data2 []byte
	}{
		{
			name:  "Single byte difference",
			data1: []byte("The quick brown fox jumps over the lazy dog"),
			data2: []byte("The quick brown fox jumps over the lazy doh"),
		},
		{
			name:  "Case change",
			data1: []byte("Hello World"),
			data2: []byte("hello World"),
		},
		{
			name:  "Extra space",
			data1: []byte("HelloWorld"),
			data2: []byte("Hello World"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hash1 := passcode.SHA3KMAC256(key, customization, tc.data1, 32)
			hash2 := passcode.SHA3KMAC256(key, customization, tc.data2, 32)

			flipRatio := bitFlipRatio(hash1, hash2)
			t.Logf("SHA3 KMAC256 %s: %.2f%% bits flipped", tc.name, flipRatio)

			if flipRatio < 25.0 {
				t.Errorf("SHA3 KMAC256 %s: insufficient avalanche effect (%.2f%%)", tc.name, flipRatio)
			}
		})
	}
}

// Test_BLAKE3_vs_SHA3_AvalancheComparison compares avalanche effects between algorithms
func Test_BLAKE3_vs_SHA3_AvalancheComparison(t *testing.T) {
	key := make([]byte, 32)
	rand.Read(key)
	customization := []byte("comparison-test")
	data := make([]byte, 128)
	rand.Read(data)

	// Test BLAKE3
	blake3Original := passcode.BLAKE3KeyedMode256(key, data)
	modifiedData := make([]byte, len(data))
	copy(modifiedData, data)
	modifiedData[0] ^= 1
	blake3Modified := passcode.BLAKE3KeyedMode256(key, modifiedData)
	blake3FlipRatio := bitFlipRatio(blake3Original, blake3Modified)

	// Test SHA3
	sha3Original := passcode.SHA3KMAC256(key, customization, data, 32)
	sha3Modified := passcode.SHA3KMAC256(key, customization, modifiedData, 32)
	sha3FlipRatio := bitFlipRatio(sha3Original, sha3Modified)

	t.Logf("Single bit flip avalanche comparison:")
	t.Logf("  BLAKE3:       %.2f%% bits flipped", blake3FlipRatio)
	t.Logf("  SHA3 KMAC256: %.2f%% bits flipped", sha3FlipRatio)

	// Both should have good avalanche effect
	if blake3FlipRatio < 25.0 || sha3FlipRatio < 25.0 {
		t.Error("One or both algorithms show insufficient avalanche effect")
	}
}

// Benchmark_HammingDistance benchmarks the hamming distance calculation
func Benchmark_HammingDistance(b *testing.B) {
	data1 := make([]byte, 32)
	data2 := make([]byte, 32)
	rand.Read(data1)
	rand.Read(data2)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hammingDistance(data1, data2)
	}
}
