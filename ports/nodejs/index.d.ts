export declare enum Algorithm {
  Sha3Kmac128 = 0,
  Sha3Kmac256 = 1,
  Blake3KeyedMode128 = 2,
  Blake3KeyedMode256 = 3,
}

export declare class Passcode {
  constructor(algorithm: Algorithm, key: Uint8Array);
  compute(data: Uint8Array): string;
  algorithmName(): string;
}

export declare function blake3KeyedMode128(key: Uint8Array, data: Uint8Array): string;
export declare function blake3KeyedMode256(key: Uint8Array, data: Uint8Array): string;
export declare function sha3Kmac128(key: Uint8Array, customization: Uint8Array, data: Uint8Array): string;
export declare function sha3Kmac256(key: Uint8Array, customization: Uint8Array, data: Uint8Array): string;
