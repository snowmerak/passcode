# Passcode for Python

Challenge-response OTP library using SHA3-KMAC and BLAKE3 algorithms.

## Installation

```bash
pip install passcode-py
```

## Usage

```python
from passcode_py import Passcode, Algorithm
import os

# Create a Passcode instance
key = os.urandom(32)
passcode = Passcode(Algorithm.BLAKE3_KEYED_MODE_256, key)

# Generate OTP from challenge
challenge = os.urandom(16)
otp = passcode.compute(challenge)

print(f"OTP: {otp}")  # 12-character hexadecimal string
```

## API

### `Algorithm`

Enum of supported algorithms:
- `Algorithm.SHA3_KMAC_128` - SHA3-KMAC with 128-bit output
- `Algorithm.SHA3_KMAC_256` - SHA3-KMAC with 256-bit output
- `Algorithm.BLAKE3_KEYED_MODE_128` - BLAKE3 Keyed Mode with 128-bit output
- `Algorithm.BLAKE3_KEYED_MODE_256` - BLAKE3 Keyed Mode with 256-bit output

### `Passcode`

#### `Passcode(algorithm, key)`
- `algorithm`: Algorithm enum value
- `key`: bytes (32 bytes recommended)

#### `passcode.compute(data)`
- `data`: bytes (challenge data)
- Returns: str (12-character hexadecimal OTP)

#### `passcode.algorithm_name()`
- Returns: str (name of the algorithm)

## Implementation

This is a WebAssembly wrapper using the Wasmer runtime, providing high performance and cryptographic safety.

**Note:** WASM bindings are currently under development. The API is stable but some functions may not be fully implemented yet.

## License

MIT
