import { cshake128, cshake256 } from '@noble/hashes/sha3';

/**
 * Left encode function for KMAC
 * Encodes an integer as a byte string with length prefix on the left
 */
function leftEncode(x: number): Uint8Array {
  if (x === 0) {
    return new Uint8Array([1, 0]);
  }

  const temp = new Uint8Array(8);
  let n = 0;

  for (let i = 7; i >= 0; i--) {
    temp[i] = x & 0xff;
    x = Math.floor(x / 256);
  }

  let startIdx = 0;
  while (startIdx < 8 && temp[startIdx] === 0) {
    startIdx++;
  }
  n = 8 - startIdx;

  const result = new Uint8Array(n + 1);
  result[0] = n;
  result.set(temp.slice(startIdx), 1);
  return result;
}

/**
 * Right encode function for KMAC
 * Encodes an integer as a byte string with length suffix on the right
 */
function rightEncode(x: number): Uint8Array {
  if (x === 0) {
    return new Uint8Array([0, 1]);
  }

  const temp = new Uint8Array(8);
  let n = 0;

  for (let i = 7; i >= 0; i--) {
    temp[i] = x & 0xff;
    x = Math.floor(x / 256);
  }

  let startIdx = 0;
  while (startIdx < 8 && temp[startIdx] === 0) {
    startIdx++;
  }
  n = 8 - startIdx;

  const result = new Uint8Array(n + 1);
  result.set(temp.slice(startIdx), 0);
  result[n] = n;
  return result;
}

/**
 * Encode a byte string with its bit length
 */
function encodeString(data: Uint8Array): Uint8Array {
  const bitLen = data.length * 8;
  const encoded = leftEncode(bitLen);

  const result = new Uint8Array(encoded.length + data.length);
  result.set(encoded, 0);
  result.set(data, encoded.length);
  return result;
}

/**
 * Bytepad function for KMAC
 * Prepends an encoding of the rate w and pads with zeros
 */
function bytepad(data: Uint8Array, w: number): Uint8Array {
  const wEncoded = leftEncode(w);
  const totalLen = wEncoded.length + data.length;

  let padLen = w - (totalLen % w);
  if (padLen === w) {
    padLen = 0;
  }

  const result = new Uint8Array(totalLen + padLen);
  result.set(wEncoded, 0);
  result.set(data, wEncoded.length);
  return result;
}

/**
 * KMAC implementation using cSHAKE
 */
function kmac(
  key: Uint8Array,
  customization: Uint8Array,
  data: Uint8Array,
  outputLen: number,
  rate: number,
  cshakeFunc: (data: Uint8Array, opts: { personalization: Uint8Array; NISTfn: Uint8Array; dkLen: number }) => Uint8Array
): Uint8Array {
  const encodedKey = encodeString(key);
  const paddedKey = bytepad(encodedKey, rate);

  // Concatenate paddedKey, data, and right_encode(outputLen * 8)
  const rightEncoded = rightEncode(outputLen * 8);
  const inputData = new Uint8Array(paddedKey.length + data.length + rightEncoded.length);
  inputData.set(paddedKey, 0);
  inputData.set(data, paddedKey.length);
  inputData.set(rightEncoded, paddedKey.length + data.length);

  return cshakeFunc(inputData, {
    personalization: customization,
    NISTfn: new TextEncoder().encode('KMAC'),
    dkLen: outputLen,
  });
}

/**
 * SHA3-KMAC128 for passcode (internal use)
 */
export function sha3KMAC128ForPasscode(key: Uint8Array, data: Uint8Array): Uint8Array {
  return kmac(key, new TextEncoder().encode('authorization'), data, 32, 168, cshake128);
}

/**
 * SHA3-KMAC128 with customizable parameters
 */
export function sha3KMAC128(key: Uint8Array, customization: Uint8Array, data: Uint8Array, outputLen: number): Uint8Array {
  return kmac(key, customization, data, outputLen, 168, cshake128);
}

/**
 * SHA3-KMAC256 for passcode (internal use)
 */
export function sha3KMAC256ForPasscode(key: Uint8Array, data: Uint8Array): Uint8Array {
  return kmac(key, new TextEncoder().encode('authorization'), data, 32, 136, cshake256);
}

/**
 * SHA3-KMAC256 with customizable parameters
 */
export function sha3KMAC256(key: Uint8Array, customization: Uint8Array, data: Uint8Array, outputLen: number): Uint8Array {
  return kmac(key, customization, data, outputLen, 136, cshake256);
}
