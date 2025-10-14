import { sha3KMAC128ForPasscode, sha3KMAC256ForPasscode } from './sha3-kmac';
import { blake3KeyedMode256, blake3KeyedMode512 } from './blake3-keyed';

/**
 * Available hash algorithms for OTP generation
 */
export enum Algorithm {
  SHA3_KMAC_128 = 'SHA3-KMAC-128',
  SHA3_KMAC_256 = 'SHA3-KMAC-256',
  BLAKE3_KEYED_MODE_128 = 'BLAKE3-Keyed-Mode-128',
  BLAKE3_KEYED_MODE_256 = 'BLAKE3-Keyed-Mode-256',
}

/**
 * Hasher function type
 */
export type Hasher = (key: Uint8Array, data: Uint8Array) => Uint8Array;

/**
 * Passcode class for Challenge-Response based OTP authentication
 */
export class Passcode {
  private algorithm: string;
  private key: Uint8Array;
  private hasher: Hasher;

  /**
   * Creates a new Passcode instance
   * @param algorithm - The hash algorithm to use
   * @param key - The secret key (shared between server and client)
   */
  constructor(algorithm: Algorithm, key: Uint8Array) {
    this.algorithm = algorithm;
    this.key = key;

    switch (algorithm) {
      case Algorithm.SHA3_KMAC_128:
        this.hasher = sha3KMAC128ForPasscode;
        break;
      case Algorithm.SHA3_KMAC_256:
        this.hasher = sha3KMAC256ForPasscode;
        break;
      case Algorithm.BLAKE3_KEYED_MODE_128:
        this.hasher = blake3KeyedMode256; // Using 256-bit output for 128-bit mode
        break;
      case Algorithm.BLAKE3_KEYED_MODE_256:
        this.hasher = blake3KeyedMode512;
        break;
      default:
        throw new Error(`Unknown hash algorithm: ${algorithm}`);
    }
  }

  /**
   * Computes an OTP from the given challenge data
   * @param data - The challenge data (typically a random value from the server)
   * @returns A 12-character hexadecimal OTP string
   */
  compute(data: Uint8Array): string {
    let hashed = this.hasher(this.key, data);

    // Ensure we have at least 6 bytes
    if (hashed.length < 6) {
      const padded = new Uint8Array(6);
      padded.set(hashed);
      hashed = padded;
    }

    // Convert first 6 bytes to hex string
    return Array.from(hashed.slice(0, 6))
      .map((b) => b.toString(16).padStart(2, '0'))
      .join('');
  }

  /**
   * Gets the algorithm being used
   */
  getAlgorithm(): string {
    return this.algorithm;
  }
}
