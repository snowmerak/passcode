import 'dart:typed_data';
import 'algorithm.dart';

/// Challenge-response OTP generator
class Passcode {
  /// The algorithm to use
  final Algorithm algorithm;
  
  /// The secret key
  final Uint8List key;
  
  /// Create a new Passcode instance
  /// 
  /// [algorithm] - The hash algorithm to use
  /// [key] - Secret key (32 bytes recommended)
  Passcode(this.algorithm, this.key) {
    if (key.isEmpty) {
      throw ArgumentError('Key cannot be empty');
    }
  }
  
  /// Compute OTP from challenge data
  /// 
  /// [data] - Challenge data
  /// Returns a 12-character hexadecimal OTP string
  String compute(Uint8List data) {
    // Simple hash-based implementation for testing
    // In production, this would use WASM or FFI
    switch (algorithm) {
      case Algorithm.sha3Kmac128:
        return _computeHash(key, data, 'sha3-128');
      case Algorithm.sha3Kmac256:
        return _computeHash(key, data, 'sha3-256');
      case Algorithm.blake3KeyedMode128:
        return _computeHash(key, data, 'blake3-128');
      case Algorithm.blake3KeyedMode256:
        return _computeHash(key, data, 'blake3-256');
    }
  }
  
  String _computeHash(Uint8List key, Uint8List data, String mode) {
    // Placeholder: Use the Rust/Go reference implementation via FFI in production
    // For now, just demonstrate the API works
    throw UnimplementedError(
      'WASM/FFI bindings not yet implemented. '
      'This is a placeholder showing the API structure. '
      'In production, this would call the Rust implementation via FFI or WASM.'
    );
  }
  
  /// Get the name of the algorithm
  String get algorithmName => algorithm.name;
}

/// Utility functions

/// Compute BLAKE3-128 hash
String blake3KeyedMode128(Uint8List key, Uint8List data) {
  throw UnimplementedError('WASM/FFI bindings not yet implemented');
}

/// Compute BLAKE3-256 hash
String blake3KeyedMode256(Uint8List key, Uint8List data) {
  throw UnimplementedError('WASM/FFI bindings not yet implemented');
}

/// Compute SHA3-KMAC-128 hash
String sha3Kmac128(Uint8List key, Uint8List customization, Uint8List data) {
  throw UnimplementedError('WASM/FFI bindings not yet implemented');
}

/// Compute SHA3-KMAC-256 hash
String sha3Kmac256(Uint8List key, Uint8List customization, Uint8List data) {
  throw UnimplementedError('WASM/FFI bindings not yet implemented');
}
