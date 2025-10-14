# Passcode - WebAssembly Port

WebAssembly bindings for the Passcode Challenge-Response based One-Time Password (OTP) authentication library. This allows you to use the high-performance Rust implementation in JavaScript, TypeScript, Python, and other environments that support WebAssembly.

## ✨ Features

- **High Performance**: Compiled from Rust to WebAssembly for near-native speed
- **Cross-Platform**: Works in browsers, Node.js, Deno, and any WASM runtime
- **Multiple Hash Algorithms**: SHA3-KMAC (128/256) and BLAKE3 Keyed Mode (128/256)
- **Type-Safe**: Full TypeScript definitions included
- **Small Bundle Size**: Optimized for size with LTO and stripping

## 📦 Installation

### For JavaScript/TypeScript (Bundler)

```bash
npm install @snowmerak/passcode-wasm
# or
yarn add @snowmerak/passcode-wasm
```

### For Node.js

```bash
npm install @snowmerak/passcode-wasm
```

## 🚀 Usage

### Browser (with bundler like webpack, vite, etc.)

```javascript
import init, { Passcode, Algorithm } from '@snowmerak/passcode-wasm';

async function main() {
  // Initialize the WASM module
  await init();

  // Create secret key (32 bytes recommended)
  const secretKey = new Uint8Array(32);
  crypto.getRandomValues(secretKey);

  // Create Passcode instance
  const passcode = new Passcode(Algorithm.Blake3KeyedMode256, secretKey);

  // Generate challenge
  const challenge = new Uint8Array(16);
  crypto.getRandomValues(challenge);

  // Compute OTP
  const otp = passcode.compute(challenge);
  console.log('Generated OTP:', otp); // 12-character hex string
}

main();
```

### Node.js

```javascript
const { Passcode, Algorithm } = require('@snowmerak/passcode-wasm');
const crypto = require('crypto');

// Create secret key
const secretKey = crypto.randomBytes(32);

// Create Passcode instance
const passcode = new Passcode(Algorithm.Blake3KeyedMode256, secretKey);

// Generate challenge
const challenge = crypto.randomBytes(16);

// Compute OTP
const otp = passcode.compute(challenge);
console.log('Generated OTP:', otp);
```

### TypeScript

```typescript
import init, { Passcode, Algorithm } from '@snowmerak/passcode-wasm';

async function authenticate(): Promise<void> {
  await init();

  const secretKey = new Uint8Array(32);
  crypto.getRandomValues(secretKey);

  const passcode = new Passcode(Algorithm.Blake3KeyedMode256, secretKey);
  const challenge = new Uint8Array(16);
  crypto.getRandomValues(challenge);

  const otp: string = passcode.compute(challenge);
  console.log('OTP:', otp);
}
```

## 📖 API Documentation

### `Algorithm` Enum

Available algorithms:

```typescript
enum Algorithm {
  Sha3Kmac128,
  Sha3Kmac256,
  Blake3KeyedMode128,
  Blake3KeyedMode256,
}
```

### `Passcode` Class

#### Constructor

```typescript
constructor(algorithm: Algorithm, key: Uint8Array): Passcode
```

Creates a new Passcode instance.

- `algorithm`: The hash algorithm to use
- `key`: The shared secret key (32 bytes recommended)

#### Methods

##### `compute(data: Uint8Array): string`

Computes a 12-character hexadecimal OTP from the given challenge data.

- `data`: The challenge data (typically a random value from the server)
- Returns: A 12-character hexadecimal string

##### `algorithmName: string` (getter)

Returns the algorithm name as a string.

### Utility Functions

#### `blake3KeyedMode256(key: Uint8Array, data: Uint8Array): Uint8Array`

BLAKE3 in keyed mode with 256-bit (32 bytes) output.

#### `blake3KeyedMode512(key: Uint8Array, data: Uint8Array): Uint8Array`

BLAKE3 in keyed mode with 512-bit (64 bytes) output.

#### `sha3Kmac128(key: Uint8Array, customization: Uint8Array, data: Uint8Array, outputLen: number): Uint8Array`

SHA3-KMAC with 128-bit security level.

#### `sha3Kmac256(key: Uint8Array, customization: Uint8Array, data: Uint8Array, outputLen: number): Uint8Array`

SHA3-KMAC with 256-bit security level.

## 🔧 Building from Source

### Prerequisites

- Rust toolchain (rustup)
- wasm-pack: `cargo install wasm-pack`

### Build Commands

```bash
# For bundler (webpack, vite, etc.)
npm run build

# For Node.js
npm run build:nodejs

# For web (no bundler)
npm run build:web

# Build all targets
npm run build:all
```

### Output

- `pkg/` - Bundler target
- `pkg-node/` - Node.js target
- `pkg-web/` - Web target (no bundler)

## 📊 Bundle Size

The WASM module is highly optimized:

- ~50-100KB (gzipped) depending on the target
- Optimized with LTO and size optimization flags
- Symbols stripped for smaller binary

## 🌐 Platform Support

- ✅ Modern browsers (Chrome, Firefox, Safari, Edge)
- ✅ Node.js 12+
- ✅ Deno
- ✅ Cloudflare Workers
- ✅ Any environment with WebAssembly support

## 🔗 Related

- [Go Implementation](../../) - Reference implementation
- [Rust Implementation](../rust/) - Source for WASM bindings
- [TypeScript Port](../typescript/) - Pure TypeScript implementation
- [Dart Port](../dart/) - Dart/Flutter implementation

## 💡 Use Cases

- **Web Applications**: Fast OTP generation in the browser
- **Serverless Functions**: Use in Edge workers and Lambda functions
- **Cross-Platform Apps**: Same library across all platforms
- **Python**: Use via `wasmtime` or similar WASM runtimes

## 📄 License

MIT License - see the [LICENSE](../../LICENSE) file for details.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
