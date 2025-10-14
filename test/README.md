# Cross-Implementation Test

This directory contains test files to verify that all language implementations produce identical outputs for the same inputs.

## Test Vectors

All tests use the same fixed values:
- **Key**: `0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef` (32 bytes)
- **Challenge**: `fedcba9876543210fedcba9876543210` (16 bytes)

## Running Individual Tests

### Go Test
```bash
cd test
go run go_main.go
```

### Rust Test
```bash
cd test
cargo run --bin rust_test
```

### Node.js Test (WASM Direct)
```bash
cd test
node node_test.js
```

### Node.js Test (Library Package)
```bash
cd test
node nodejs_lib_test.js
```

### Python Test
```bash
cd test
pip install wasmtime  # Required dependency
python3 python_test.py
```

### Dart Test
```bash
cd test
dart run dart_test.dart
```

## Automated Test Suite

Run all tests and compare outputs:
```bash
cd test
./run_all_tests.sh
```

This script will:
1. Run Go, Rust, Node.js (WASM), Node.js (Library) tests
2. Attempt to run Python and Dart tests (if dependencies are installed)
3. Compare all outputs and verify they match
4. Display a color-coded summary of results

### Prerequisites

**Core (Required for full test suite):**
- Go 1.25+
- Rust with cargo
- Node.js 16+

**Optional (for additional language tests):**
- Python 3.8+ with `wasmtime` package: `pip install wasmtime`
- Dart SDK 3.0+

## Current Test Status

| Implementation | Status | Notes |
|---------------|--------|-------|
| Go | ✅ Production | Reference implementation |
| Rust | ✅ Production | Native performance |
| Node.js (WASM) | ✅ Production | Direct WASM binding |
| Node.js (Library) | ✅ Production | npm package wrapper |
| Python | 🚧 Partial | API complete, needs `wasmtime` |
| Dart | 🚧 Partial | API complete, needs setup |

## Expected Output

All working implementations produce identical OTP values:

```
SHA3-KMAC-128       : 2ce05573dd4e
SHA3-KMAC-256       : f391e239e588
BLAKE3-Keyed-128    : 2ce4568631de
BLAKE3-Keyed-256    : 2ce4568631de
```

## Test Files

- `go_main.go` - Go implementation test
- `rust_test.rs` - Rust implementation test
- `node_test.js` - Node.js WASM direct binding test
- `nodejs_lib_test.js` - Node.js library package test
- `python_test.py` - Python implementation test
- `dart_test.dart` - Dart implementation test
- `run_all_tests.sh` - Automated test runner script

## Troubleshooting

### Python "ModuleNotFoundError: No module named 'wasmtime'"
Install wasmtime: `pip install wasmtime` or `python3 -m pip install wasmtime`

### Dart "Error: Not found: 'package:passcode/passcode.dart'"
The Dart package needs to be published or symlinked. For local testing, copy the test file into the `ports/dart/example/` directory.

### Cargo not found
Ensure Rust is installed and `~/.cargo/bin` is in your PATH, or use the full path: `~/.cargo/bin/cargo`
