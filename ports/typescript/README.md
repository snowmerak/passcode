# Passcode - TypeScript Port

TypeScript implementation of the Passcode Challenge-Response based One-Time Password (OTP) authentication library.

## ✨ Features

- **Challenge-Response Mechanism**: Secure authentication where the server sends a random challenge
- **Multiple Hash Algorithms**: 
  - SHA3-KMAC (128/256)
  - BLAKE3 (Keyed Mode, 128/256)
- **Flexible Security Levels**: Choose between 128-bit and 256-bit security strengths
- **Type-Safe API**: Full TypeScript support with type definitions
- **Zero Dependencies**: Uses well-tested cryptographic libraries (`@noble/hashes`, `blake3`)

## ⚙️ How It Works

1. **Authentication Request**: Client requests authentication with its ID
2. **Challenge Generation**: Server generates a cryptographically secure random challenge
3. **Challenge Transmission**: Server sends the challenge to the client
4. **OTP Computation**: Both server and client compute the OTP using the shared secret key and challenge
5. **Response Submission**: Client submits the computed OTP
6. **Verification**: Server compares the OTPs - if they match, authentication succeeds

The secret key is never transmitted over the network, and each OTP is valid for only one challenge.

## 📦 Installation

```bash
npm install @snowmerak/passcode
# or
yarn add @snowmerak/passcode
# or
pnpm add @snowmerak/passcode
```

## 🚀 Usage

### Basic Example

```typescript
import { Passcode, Algorithm } from '@snowmerak/passcode';

// 1. Shared secret key between server and client (32 bytes recommended)
const secretKey = new Uint8Array(32);
crypto.getRandomValues(secretKey);

// 2. Create a new Passcode instance
const passcode = new Passcode(Algorithm.BLAKE3_KEYED_MODE_256, secretKey);

// 3. Server generates a challenge (a random value)
const challenge = new Uint8Array(16);
crypto.getRandomValues(challenge);

// 4. Compute the OTP based on the challenge
const otp = passcode.compute(challenge);
console.log('Generated OTP:', otp); // 12-character hexadecimal string

// 5. Server verifies by computing the same OTP and comparing
```

### Available Algorithms

```typescript
import { Algorithm } from '@snowmerak/passcode';

// Choose one of:
Algorithm.SHA3_KMAC_128           // SHA3-KMAC with 128-bit security
Algorithm.SHA3_KMAC_256           // SHA3-KMAC with 256-bit security
Algorithm.BLAKE3_KEYED_MODE_128   // BLAKE3 Keyed Mode with 128-bit security
Algorithm.BLAKE3_KEYED_MODE_256   // BLAKE3 Keyed Mode with 256-bit security
```

### Advanced Usage

#### Using SHA3-KMAC directly

```typescript
import { sha3KMAC256 } from '@snowmerak/passcode';

const key = new Uint8Array(32);
const customization = new TextEncoder().encode('my-app');
const data = new TextEncoder().encode('some data');
const outputLength = 64; // bytes

const hash = sha3KMAC256(key, customization, data, outputLength);
```

#### Using BLAKE3 Keyed Mode directly

```typescript
import { blake3KeyedMode256, blake3KeyedMode512 } from '@snowmerak/passcode';

const key = new Uint8Array(32);
const data = new TextEncoder().encode('some data');

const hash256 = blake3KeyedMode256(key, data); // 32 bytes
const hash512 = blake3KeyedMode512(key, data); // 64 bytes
```

## 🧪 Development

### Install dependencies
```bash
npm install
```

### Build
```bash
npm run build
```

### Run tests
```bash
npm test
```

### Run tests with coverage
```bash
npm run test:coverage
```

### Lint
```bash
npm run lint
```

### Format code
```bash
npm run format
```

## 📖 API Documentation

### `Passcode` Class

#### Constructor
```typescript
constructor(algorithm: Algorithm, key: Uint8Array)
```

Creates a new Passcode instance with the specified algorithm and secret key.

- `algorithm`: The hash algorithm to use (see `Algorithm` enum)
- `key`: The shared secret key (Uint8Array, 32 bytes recommended)

#### Methods

##### `compute(data: Uint8Array): string`

Computes a 12-character hexadecimal OTP from the given challenge data.

- `data`: The challenge data (typically a random value from the server)
- Returns: A 12-character hexadecimal string

##### `getAlgorithm(): string`

Returns the algorithm name being used.

### Functions

#### `sha3KMAC128(key, customization, data, outputLen)`
SHA3-KMAC with 128-bit security level and customizable output length.

#### `sha3KMAC256(key, customization, data, outputLen)`
SHA3-KMAC with 256-bit security level and customizable output length.

#### `blake3KeyedMode256(key, data)`
BLAKE3 in keyed mode with 256-bit (32 bytes) output.

#### `blake3KeyedMode512(key, data)`
BLAKE3 in keyed mode with 512-bit (64 bytes) output.

## 🔗 Related

- [Go Implementation](../../) - Reference implementation
- [Dart Port](../dart/) - Dart implementation
- [Rust Port](../rust/) - Rust implementation

## 📄 License

MIT License - see the [LICENSE](../../LICENSE) file for details.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
