# Passcode - Rust Port

Rust implementation of the Passcode Challenge-Response based One-Time Password (OTP) authentication library.

## ✨ Features

- **Challenge-Response Mechanism**: Secure authentication where the server sends a random challenge
- **Multiple Hash Algorithms**: 
  - SHA3-KMAC (128/256)
  - BLAKE3 (Keyed Mode, 128/256)
- **Flexible Security Levels**: Choose between 128-bit and 256-bit security strengths
- **Type-Safe API**: Leverages Rust's type system for safety and performance
- **Zero-Cost Abstractions**: No runtime overhead
- **Memory Safety**: Rust's ownership system prevents common security vulnerabilities

## ⚙️ How It Works

1. **Authentication Request**: Client requests authentication with its ID
2. **Challenge Generation**: Server generates a cryptographically secure random challenge
3. **Challenge Transmission**: Server sends the challenge to the client
4. **OTP Computation**: Both server and client compute the OTP using the shared secret key and challenge
5. **Response Submission**: Client submits the computed OTP
6. **Verification**: Server compares the OTPs - if they match, authentication succeeds

The secret key is never transmitted over the network, and each OTP is valid for only one challenge.

## 📦 Installation

Add this to your `Cargo.toml`:

```toml
[dependencies]
passcode = { git = "https://github.com/snowmerak/passcode", branch = "main", subdirectory = "ports/rust" }
```

Or if published to crates.io:

```toml
[dependencies]
passcode = "1.0"
```

## 🚀 Usage

### Basic Example

```rust
use passcode::{Passcode, Algorithm};
use rand::RngCore;

fn main() {
    let mut rng = rand::thread_rng();

    // 1. Shared secret key between server and client (32 bytes recommended)
    let mut secret_key = vec![0u8; 32];
    rng.fill_bytes(&mut secret_key);

    // 2. Create a new Passcode instance
    let passcode = Passcode::new(Algorithm::Blake3KeyedMode256, secret_key);

    // 3. Server generates a challenge (a random value)
    let mut challenge = vec![0u8; 16];
    rng.fill_bytes(&mut challenge);

    // 4. Compute the OTP based on the challenge
    let otp = passcode.compute(&challenge);
    println!("Generated OTP: {}", otp); // 12-character hexadecimal string

    // 5. Server verifies by computing the same OTP and comparing
}
```

### Available Algorithms

```rust
use passcode::Algorithm;

// Choose one of:
Algorithm::Sha3Kmac128           // SHA3-KMAC with 128-bit security
Algorithm::Sha3Kmac256           // SHA3-KMAC with 256-bit security
Algorithm::Blake3KeyedMode128    // BLAKE3 Keyed Mode with 128-bit security
Algorithm::Blake3KeyedMode256    // BLAKE3 Keyed Mode with 256-bit security
```

### Advanced Usage

#### Using SHA3-KMAC directly

```rust
use passcode::sha3_kmac256;

let key = vec![0u8; 32];
let customization = b"my-app";
let data = b"some data";
let output_length = 64; // bytes

let hash = sha3_kmac256(&key, customization, data, output_length);
```

#### Using BLAKE3 Keyed Mode directly

```rust
use passcode::{blake3_keyed_mode256, blake3_keyed_mode512};

let key = vec![0u8; 32];
let data = b"some data";

let hash256 = blake3_keyed_mode256(&key, data); // 32 bytes
let hash512 = blake3_keyed_mode512(&key, data); // 64 bytes
```

## 🧪 Development

### Build
```bash
cargo build
```

### Run tests
```bash
cargo test
```

### Run tests with output
```bash
cargo test -- --nocapture
```

### Run example
```bash
cargo run --example basic
```

### Generate documentation
```bash
cargo doc --open
```

### Check formatting
```bash
cargo fmt --check
```

### Run clippy linter
```bash
cargo clippy -- -D warnings
```

## 📖 API Documentation

### `Passcode` Struct

#### Constructor
```rust
pub fn new(algorithm: Algorithm, key: Vec<u8>) -> Self
```

Creates a new Passcode instance with the specified algorithm and secret key.

- `algorithm`: The hash algorithm to use (see `Algorithm` enum)
- `key`: The shared secret key (32 bytes recommended)

#### Methods

##### `pub fn compute(&self, data: &[u8]) -> String`

Computes a 12-character hexadecimal OTP from the given challenge data.

- `data`: The challenge data (typically a random value from the server)
- Returns: A 12-character hexadecimal string

##### `pub fn algorithm(&self) -> Algorithm`

Returns the algorithm enum value being used.

##### `pub fn algorithm_name(&self) -> &'static str`

Returns the algorithm name as a string.

### `Algorithm` Enum

```rust
pub enum Algorithm {
    Sha3Kmac128,
    Sha3Kmac256,
    Blake3KeyedMode128,
    Blake3KeyedMode256,
}
```

Implements `Display`, `Debug`, `Clone`, `Copy`, `PartialEq`, `Eq`.

### Functions

#### `sha3_kmac128(key, customization, data, output_len) -> Vec<u8>`
SHA3-KMAC with 128-bit security level and customizable output length.

#### `sha3_kmac256(key, customization, data, output_len) -> Vec<u8>`
SHA3-KMAC with 256-bit security level and customizable output length.

#### `blake3_keyed_mode256(key, data) -> Vec<u8>`
BLAKE3 in keyed mode with 256-bit (32 bytes) output.

#### `blake3_keyed_mode512(key, data) -> Vec<u8>`
BLAKE3 in keyed mode with 512-bit (64 bytes) output.

## 🔗 Related

- [Go Implementation](../../) - Reference implementation
- [TypeScript Port](../typescript/) - TypeScript/JavaScript implementation
- [Dart Port](../dart/) - Dart/Flutter implementation

## 🔒 Security Notes

- Uses well-tested cryptographic libraries: `sha3` and `blake3` crates
- Memory safety guaranteed by Rust's ownership system
- No unsafe code in the implementation
- Constant-time operations where applicable

## 📄 License

MIT License - see the [LICENSE](../../LICENSE) file for details.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📊 Benchmarks

To run benchmarks:

```bash
cargo bench
```

(Benchmarks to be added in future releases)
