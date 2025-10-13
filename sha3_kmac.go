package passcode

import "golang.org/x/crypto/sha3"

func leftEncode(x uint64) []byte {
	if x == 0 {
		return []byte{1, 0}
	}

	var temp [8]byte
	n := 0

	for i := 7; i >= 0; i-- {
		temp[i] = byte(x & 0xff)
		x >>= 8
	}

	startIdx := 0
	for startIdx < 8 && temp[startIdx] == 0 {
		startIdx++
	}
	n = 8 - startIdx

	result := make([]byte, n+1)
	result[0] = byte(n)
	copy(result[1:], temp[startIdx:])
	return result
}

func rightEncode(x uint64) []byte {
	if x == 0 {
		return []byte{0, 1}
	}

	var temp [8]byte
	n := 0

	for i := 7; i >= 0; i-- {
		temp[i] = byte(x & 0xff)
		x >>= 8
	}

	startIdx := 0
	for startIdx < 8 && temp[startIdx] == 0 {
		startIdx++
	}
	n = 8 - startIdx

	result := make([]byte, n+1)
	copy(result, temp[startIdx:])
	result[n] = byte(n)
	return result
}

func encodeString(data []byte) []byte {
	bitLen := uint64(len(data) * 8)
	encoded := leftEncode(bitLen)

	result := make([]byte, len(encoded)+len(data))
	copy(result, encoded)
	copy(result[len(encoded):], data)
	return result
}

func bytepad(data []byte, w int) []byte {
	// left_encode(w)
	wEncoded := leftEncode(uint64(w))

	totalLen := len(wEncoded) + len(data)

	padLen := w - (totalLen % w)
	if padLen == w {
		padLen = 0
	}

	result := make([]byte, totalLen+padLen)
	copy(result, wEncoded)
	copy(result[len(wEncoded):], data)

	return result
}

func SHA3KMAC128(key, customization, data []byte, outputLen int) []byte {
	return kmac(key, customization, data, outputLen, 168, sha3.NewCShake128)
}

func SHA3KMAC256(key, customization, data []byte, outputLen int) []byte {
	return kmac(key, customization, data, outputLen, 136, sha3.NewCShake256)
}

func kmac(key, customization, data []byte, outputLen int, rate int, newCShake func([]byte, []byte) sha3.ShakeHash) []byte {
	encodedKey := encodeString(key)
	paddedKey := bytepad(encodedKey, rate)

	hasher := newCShake([]byte("KMAC"), customization)

	hasher.Write(paddedKey)
	hasher.Write(data)
	hasher.Write(rightEncode(uint64(outputLen * 8)))

	output := make([]byte, outputLen)
	hasher.Read(output)
	return output
}
