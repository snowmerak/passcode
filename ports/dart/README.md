# Passcode for Dart

Challenge-response OTP library using SHA3-KMAC and BLAKE3 algorithms.

## Installation

Add this to your `pubspec.yaml`:

```yaml
dependencies:
  passcode: ^1.0.0
```

## Usage

```dart
import 'dart:typed_data';
import 'dart:math';
import 'package:passcode/passcode.dart';

void main() {
  // Generate random key
  final random = Random.secure();
  final key = Uint8List.fromList(
    List.generate(32, (_) => random.nextInt(256))
  );
  
  // Create Passcode instance
  final passcode = Passcode(Algorithm.blake3KeyedMode256, key);
  
  // Generate OTP from challenge
  final challenge = Uint8List.fromList(
    List.generate(16, (_) => random.nextInt(256))
  );
  final otp = passcode.compute(challenge);
  
  print('OTP: $otp'); // 12-character hexadecimal string
}
```

## API

### `Algorithm`

Enum of supported algorithms:
- `Algorithm.sha3Kmac128` - SHA3-KMAC with 128-bit output
- `Algorithm.sha3Kmac256` - SHA3-KMAC with 256-bit output
- `Algorithm.blake3KeyedMode128` - BLAKE3 Keyed Mode with 128-bit output
- `Algorithm.blake3KeyedMode256` - BLAKE3 Keyed Mode with 256-bit output

### `Passcode`

#### Constructor
```dart
Passcode(Algorithm algorithm, Uint8List key)
```
- `algorithm`: Algorithm enum value
- `key`: Secret key (32 bytes recommended)

#### Methods

```dart
String compute(Uint8List data)
```
Compute OTP from challenge data.
- `data`: Challenge data (Uint8List)
- Returns: 12-character hexadecimal OTP string

```dart
String get algorithmName
```
Get the name of the algorithm.

### Utility Functions

```dart
String blake3KeyedMode128(Uint8List key, Uint8List data);
String blake3KeyedMode256(Uint8List key, Uint8List data);
String sha3Kmac128(Uint8List key, Uint8List customization, Uint8List data);
String sha3Kmac256(Uint8List key, Uint8List customization, Uint8List data);
```

## Platform Support

- ✅ Flutter Web (via WASM)
- ✅ Flutter Desktop (via FFI, planned)
- ✅ Flutter Mobile (via FFI, planned)
- ✅ Dart CLI (via FFI, planned)

## Implementation

This library uses WebAssembly for web platforms and FFI for native platforms, wrapping the high-performance Rust implementation.

**Note:** WASM/FFI bindings are currently under development. The API is stable but the implementation is pending.

## License

MIT
