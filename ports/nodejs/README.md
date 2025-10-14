# Passcode for Node.js

Challenge-response OTP library using SHA3-KMAC and BLAKE3 algorithms.

## Installation

```bash
npm install @snowmerak/passcode
```

## Usage

```javascript
const { Passcode, Algorithm } = require('@snowmerak/passcode');
const crypto = require('crypto');

// Create a Passcode instance
const key = crypto.randomBytes(32);
const passcode = new Passcode(Algorithm.Blake3KeyedMode256, key);

// Generate OTP from challenge
const challenge = crypto.randomBytes(16);
const otp = passcode.compute(challenge);

console.log('OTP:', otp); // 12-character hexadecimal string
```

## API

### `Algorithm`

Enum of supported algorithms:
- `Algorithm.Sha3Kmac128` - SHA3-KMAC with 128-bit output
- `Algorithm.Sha3Kmac256` - SHA3-KMAC with 256-bit output
- `Algorithm.Blake3KeyedMode128` - BLAKE3 Keyed Mode with 128-bit output
- `Algorithm.Blake3KeyedMode256` - BLAKE3 Keyed Mode with 256-bit output

### `Passcode`

#### `new Passcode(algorithm, key)`
- `algorithm`: Algorithm enum value
- `key`: Uint8Array or Buffer (32 bytes recommended)

#### `passcode.compute(data)`
- `data`: Uint8Array or Buffer (challenge data)
- Returns: String (12-character hexadecimal OTP)

#### `passcode.algorithmName()`
- Returns: String (name of the algorithm)

### Utility Functions

```javascript
const { blake3KeyedMode256, sha3Kmac256 } = require('@snowmerak/passcode');

// Direct hash computation
const hash1 = blake3KeyedMode256(key, data);
const hash2 = sha3Kmac256(key, customization, data);
```

## Implementation

This is a WebAssembly wrapper around the native Rust implementation, providing high performance and cryptographic safety.

## License

MIT
