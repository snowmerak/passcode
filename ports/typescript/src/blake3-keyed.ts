import { hash } from 'blake3';

/**
 * BLAKE3 keyed mode implementation
 * Uses BLAKE3 in keyed hashing mode with a hashed key
 */
function blake3KeyedMode(key: Uint8Array, data: Uint8Array, outLen: number): Uint8Array {
  // Hash the key first to get a 32-byte key
  const hashedKey = hash(key);
  
  // Use BLAKE3 keyed hash with the hashed key
  const hasher = hash(data, { key: hashedKey });
  
  // Return the specified output length
  return hasher.slice(0, outLen);
}

/**
 * BLAKE3 keyed mode with 256-bit (32 bytes) output
 */
export function blake3KeyedMode256(key: Uint8Array, data: Uint8Array): Uint8Array {
  return blake3KeyedMode(key, data, 32);
}

/**
 * BLAKE3 keyed mode with 512-bit (64 bytes) output
 */
export function blake3KeyedMode512(key: Uint8Array, data: Uint8Array): Uint8Array {
  return blake3KeyedMode(key, data, 64);
}
