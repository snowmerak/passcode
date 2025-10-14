import { Passcode, Algorithm } from '../src/passcode';
import { blake3KeyedMode256, blake3KeyedMode512 } from '../src/blake3-keyed';
import { sha3KMAC128, sha3KMAC256 } from '../src/sha3-kmac';

describe('Passcode', () => {
  describe('constructor', () => {
    it('should create instance with SHA3-KMAC-128', () => {
      const key = new Uint8Array(32);
      crypto.getRandomValues(key);
      const pc = new Passcode(Algorithm.SHA3_KMAC_128, key);
      expect(pc).toBeInstanceOf(Passcode);
      expect(pc.getAlgorithm()).toBe('SHA3-KMAC-128');
    });

    it('should create instance with SHA3-KMAC-256', () => {
      const key = new Uint8Array(32);
      crypto.getRandomValues(key);
      const pc = new Passcode(Algorithm.SHA3_KMAC_256, key);
      expect(pc).toBeInstanceOf(Passcode);
      expect(pc.getAlgorithm()).toBe('SHA3-KMAC-256');
    });

    it('should create instance with BLAKE3-Keyed-Mode-128', () => {
      const key = new Uint8Array(32);
      crypto.getRandomValues(key);
      const pc = new Passcode(Algorithm.BLAKE3_KEYED_MODE_128, key);
      expect(pc).toBeInstanceOf(Passcode);
      expect(pc.getAlgorithm()).toBe('BLAKE3-Keyed-Mode-128');
    });

    it('should create instance with BLAKE3-Keyed-Mode-256', () => {
      const key = new Uint8Array(32);
      crypto.getRandomValues(key);
      const pc = new Passcode(Algorithm.BLAKE3_KEYED_MODE_256, key);
      expect(pc).toBeInstanceOf(Passcode);
      expect(pc.getAlgorithm()).toBe('BLAKE3-Keyed-Mode-256');
    });

    it('should throw error for unknown algorithm', () => {
      const key = new Uint8Array(32);
      expect(() => {
        new Passcode('UNKNOWN' as Algorithm, key);
      }).toThrow('Unknown hash algorithm');
    });
  });

  describe('compute', () => {
    it('should generate 12-character hex OTP', () => {
      const key = new Uint8Array(32);
      crypto.getRandomValues(key);
      const challenge = new Uint8Array(16);
      crypto.getRandomValues(challenge);

      const pc = new Passcode(Algorithm.BLAKE3_KEYED_MODE_256, key);
      const otp = pc.compute(challenge);

      expect(otp).toHaveLength(12);
      expect(/^[0-9a-f]{12}$/.test(otp)).toBe(true);
    });

    it('should generate consistent OTP for same key and challenge', () => {
      const key = new Uint8Array(32);
      crypto.getRandomValues(key);
      const challenge = new Uint8Array(16);
      crypto.getRandomValues(challenge);

      const pc1 = new Passcode(Algorithm.BLAKE3_KEYED_MODE_256, key);
      const pc2 = new Passcode(Algorithm.BLAKE3_KEYED_MODE_256, key);

      const otp1 = pc1.compute(challenge);
      const otp2 = pc2.compute(challenge);

      expect(otp1).toBe(otp2);
    });

    it('should generate different OTP for different challenges', () => {
      const key = new Uint8Array(32);
      crypto.getRandomValues(key);
      const challenge1 = new Uint8Array(16);
      const challenge2 = new Uint8Array(16);
      crypto.getRandomValues(challenge1);
      crypto.getRandomValues(challenge2);

      const pc = new Passcode(Algorithm.BLAKE3_KEYED_MODE_256, key);
      const otp1 = pc.compute(challenge1);
      const otp2 = pc.compute(challenge2);

      expect(otp1).not.toBe(otp2);
    });

    it('should generate different OTP for different keys', () => {
      const key1 = new Uint8Array(32);
      const key2 = new Uint8Array(32);
      crypto.getRandomValues(key1);
      crypto.getRandomValues(key2);
      const challenge = new Uint8Array(16);
      crypto.getRandomValues(challenge);

      const pc1 = new Passcode(Algorithm.BLAKE3_KEYED_MODE_256, key1);
      const pc2 = new Passcode(Algorithm.BLAKE3_KEYED_MODE_256, key2);

      const otp1 = pc1.compute(challenge);
      const otp2 = pc2.compute(challenge);

      expect(otp1).not.toBe(otp2);
    });

    it('should handle empty challenge data', () => {
      const key = new Uint8Array(32);
      crypto.getRandomValues(key);
      const challenge = new Uint8Array(0);

      const pc = new Passcode(Algorithm.BLAKE3_KEYED_MODE_256, key);
      const otp = pc.compute(challenge);

      expect(otp).toHaveLength(12);
      expect(/^[0-9a-f]{12}$/.test(otp)).toBe(true);
    });
  });

  describe('algorithms comparison', () => {
    it('should generate different OTPs for different algorithms', () => {
      const key = new Uint8Array(32);
      crypto.getRandomValues(key);
      const challenge = new Uint8Array(16);
      crypto.getRandomValues(challenge);

      const pc1 = new Passcode(Algorithm.SHA3_KMAC_256, key);
      const pc2 = new Passcode(Algorithm.BLAKE3_KEYED_MODE_256, key);

      const otp1 = pc1.compute(challenge);
      const otp2 = pc2.compute(challenge);

      expect(otp1).not.toBe(otp2);
    });
  });
});

describe('BLAKE3 Keyed Mode', () => {
  it('should produce 32 bytes output for blake3KeyedMode256', () => {
    const key = new Uint8Array(32);
    crypto.getRandomValues(key);
    const data = new Uint8Array(64);
    crypto.getRandomValues(data);

    const result = blake3KeyedMode256(key, data);
    expect(result).toHaveLength(32);
  });

  it('should produce 64 bytes output for blake3KeyedMode512', () => {
    const key = new Uint8Array(32);
    crypto.getRandomValues(key);
    const data = new Uint8Array(64);
    crypto.getRandomValues(data);

    const result = blake3KeyedMode512(key, data);
    expect(result).toHaveLength(64);
  });

  it('should produce consistent results', () => {
    const key = new Uint8Array(32);
    crypto.getRandomValues(key);
    const data = new Uint8Array(64);
    crypto.getRandomValues(data);

    const result1 = blake3KeyedMode256(key, data);
    const result2 = blake3KeyedMode256(key, data);

    expect(result1).toEqual(result2);
  });

  it('should produce different results for different keys', () => {
    const key1 = new Uint8Array(32);
    const key2 = new Uint8Array(32);
    crypto.getRandomValues(key1);
    crypto.getRandomValues(key2);
    const data = new Uint8Array(64);
    crypto.getRandomValues(data);

    const result1 = blake3KeyedMode256(key1, data);
    const result2 = blake3KeyedMode256(key2, data);

    expect(result1).not.toEqual(result2);
  });
});

describe('SHA3 KMAC', () => {
  it('should produce specified output length for SHA3KMAC128', () => {
    const key = new Uint8Array(32);
    crypto.getRandomValues(key);
    const customization = new Uint8Array(8);
    const data = new Uint8Array(64);
    crypto.getRandomValues(data);

    const result = sha3KMAC128(key, customization, data, 32);
    expect(result).toHaveLength(32);
  });

  it('should produce specified output length for SHA3KMAC256', () => {
    const key = new Uint8Array(32);
    crypto.getRandomValues(key);
    const customization = new Uint8Array(8);
    const data = new Uint8Array(64);
    crypto.getRandomValues(data);

    const result = sha3KMAC256(key, customization, data, 64);
    expect(result).toHaveLength(64);
  });

  it('should produce consistent results', () => {
    const key = new Uint8Array(32);
    crypto.getRandomValues(key);
    const customization = new Uint8Array(8);
    const data = new Uint8Array(64);
    crypto.getRandomValues(data);

    const result1 = sha3KMAC256(key, customization, data, 32);
    const result2 = sha3KMAC256(key, customization, data, 32);

    expect(result1).toEqual(result2);
  });

  it('should produce different results for different keys', () => {
    const key1 = new Uint8Array(32);
    const key2 = new Uint8Array(32);
    crypto.getRandomValues(key1);
    crypto.getRandomValues(key2);
    const customization = new Uint8Array(8);
    const data = new Uint8Array(64);
    crypto.getRandomValues(data);

    const result1 = sha3KMAC256(key1, customization, data, 32);
    const result2 = sha3KMAC256(key2, customization, data, 32);

    expect(result1).not.toEqual(result2);
  });

  it('should produce different results for different customization', () => {
    const key = new Uint8Array(32);
    crypto.getRandomValues(key);
    const customization1 = new TextEncoder().encode('custom1');
    const customization2 = new TextEncoder().encode('custom2');
    const data = new Uint8Array(64);
    crypto.getRandomValues(data);

    const result1 = sha3KMAC256(key, customization1, data, 32);
    const result2 = sha3KMAC256(key, customization2, data, 32);

    expect(result1).not.toEqual(result2);
  });
});
