import 'dart:io';
import 'dart:typed_data';
import 'algorithm.dart';

/// Challenge-response OTP generator using Node.js WASM bridge
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
  Future<String> compute(Uint8List data) async {
    // Use Node.js WASM implementation via subprocess
    final algorithmIndex = algorithm.index;
    final keyHex = _bytesToHex(key);
    final dataHex = _bytesToHex(data);
    
    // Find the Node.js package
    final projectRoot = _findProjectRoot();
    final nodejsPath = '$projectRoot/ports/nodejs';
    
    // Create Node.js script
    final script = '''
const { Passcode, Algorithm } = require('$nodejsPath');
const key = Buffer.from('$keyHex', 'hex');
const data = Buffer.from('$dataHex', 'hex');
const passcode = new Passcode($algorithmIndex, key);
console.log(passcode.compute(data));
''';
    
    // Run Node.js
    final result = await Process.run(
      'node',
      ['-e', script],
    );
    
    if (result.exitCode != 0) {
      throw Exception('Node.js execution failed: ${result.stderr}');
    }
    
    return result.stdout.toString().trim();
  }
  
  /// Get the name of the algorithm
  String get algorithmName => algorithm.name;
  
  String _bytesToHex(Uint8List bytes) {
    return bytes.map((b) => b.toRadixString(16).padLeft(2, '0')).join();
  }
  
  String _findProjectRoot() {
    // Navigate up from ports/dart to find project root
    var current = Directory.current.path;
    
    // Check if we're already in project root
    if (Directory('$current/ports').existsSync()) {
      return current;
    }
    
    // Navigate up
    while (current != '/') {
      if (Directory('$current/ports').existsSync()) {
        return current;
      }
      current = Directory(current).parent.path;
    }
    
    throw Exception('Could not find project root');
  }
}

/// Utility functions

/// Compute BLAKE3-128 hash
Future<String> blake3KeyedMode128(Uint8List key, Uint8List data) async {
  final passcode = Passcode(Algorithm.blake3KeyedMode128, key);
  return await passcode.compute(data);
}

/// Compute BLAKE3-256 hash
Future<String> blake3KeyedMode256(Uint8List key, Uint8List data) async {
  final passcode = Passcode(Algorithm.blake3KeyedMode256, key);
  return await passcode.compute(data);
}

/// Compute SHA3-KMAC-128 hash
Future<String> sha3Kmac128(Uint8List key, Uint8List customization, Uint8List data) async {
  final passcode = Passcode(Algorithm.sha3Kmac128, key);
  return await passcode.compute(data);
}

/// Compute SHA3-KMAC-256 hash
Future<String> sha3Kmac256(Uint8List key, Uint8List customization, Uint8List data) async {
  final passcode = Passcode(Algorithm.sha3Kmac256, key);
  return await passcode.compute(data);
}
