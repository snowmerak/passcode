import 'dart:io';
import 'dart:typed_data';
import 'dart:convert';
import 'algorithm.dart';

/// Challenge-response OTP generator using Go implementation via subprocess
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
    // Use Go implementation via subprocess
    final algorithmIndex = algorithm.index;
    final keyHex = _bytesToHex(key);
    final dataHex = _bytesToHex(data);
    
    // Create a temporary Go program
    final goCode = '''
package main

import (
\t"encoding/hex"
\t"fmt"
\t"os"
)

// Algorithm represents the hash algorithm to use
type Algorithm int

const (
\tAlgorithmSHA3KMAC128 Algorithm = iota
\tAlgorithmSHA3KMAC256
\tAlgorithmBLAKE3KeyedMode128
\tAlgorithmBLAKE3KeyedMode256
)

func main() {
\tif len(os.Args) != 4 {
\t\tfmt.Fprintf(os.Stderr, "Usage: %s <algorithm> <key_hex> <data_hex>\\n", os.Args[0])
\t\tos.Exit(1)
\t}
\t
\tvar algorithm Algorithm
\tswitch os.Args[1] {
\tcase "0":
\t\talgorithm = AlgorithmSHA3KMAC128
\tcase "1":
\t\talgorithm = AlgorithmSHA3KMAC256
\tcase "2":
\t\talgorithm = AlgorithmBLAKE3KeyedMode128
\tcase "3":
\t\talgorithm = AlgorithmBLAKE3KeyedMode256
\tdefault:
\t\tfmt.Fprintf(os.Stderr, "Invalid algorithm: %s\\n", os.Args[1])
\t\tos.Exit(1)
\t}
\t
\tkey, err := hex.DecodeString(os.Args[2])
\tif err != nil {
\t\tfmt.Fprintf(os.Stderr, "Invalid key hex: %v\\n", err)
\t\tos.Exit(1)
\t}
\t
\tdata, err := hex.DecodeString(os.Args[3])
\tif err != nil {
\t\tfmt.Fprintf(os.Stderr, "Invalid data hex: %v\\n", err)
\t\tos.Exit(1)
\t}
\t
\tpasscode := NewPasscode(algorithm, key)
\totp := passcode.Compute(data)
\tfmt.Print(otp)
}

// Passcode is the main struct for OTP generation
type Passcode struct {
\talgorithm Algorithm
\tkey       []byte
}

// NewPasscode creates a new Passcode instance
func NewPasscode(algorithm Algorithm, key []byte) *Passcode {
\treturn &Passcode{
\t\talgorithm: algorithm,
\t\tkey:       key,
\t}
}

// Compute generates an OTP from challenge data
func (p *Passcode) Compute(data []byte) string {
\tvar hash []byte
\tswitch p.algorithm {
\tcase AlgorithmSHA3KMAC128:
\t\thash = SHA3KMAC128(p.key, []byte("OTP"), data)
\tcase AlgorithmSHA3KMAC256:
\t\thash = SHA3KMAC256(p.key, []byte("OTP"), data)
\tcase AlgorithmBLAKE3KeyedMode128:
\t\thash = BLAKE3KeyedMode128(p.key, data)
\tcase AlgorithmBLAKE3KeyedMode256:
\t\thash = BLAKE3KeyedMode256(p.key, data)
\t}
\treturn hex.EncodeToString(hash)
}
''';

    // Write temporary files
    final tempDir = await Directory.systemTemp.createTemp('passcode_');
    final goFile = File('${tempDir.path}/main.go');
    final goModFile = File('${tempDir.path}/go.mod');
    
    // Copy the Go implementation files
    final projectRoot = _findProjectRoot();
    final passcodeGo = File('$projectRoot/passcode.go');
    final sha3KmacGo = File('$projectRoot/sha3_kmac.go');
    final blake3KeyedGo = File('$projectRoot/blake3_keyedmode.go');
    
    await goFile.writeAsString(goCode);
    await goModFile.writeAsString('''
module temppasscode

go 1.25.2

require golang.org/x/crypto v0.32.0
''');
    
    // Copy implementation files
    await File('${tempDir.path}/passcode.go').writeAsBytes(await passcodeGo.readAsBytes());
    await File('${tempDir.path}/sha3_kmac.go').writeAsBytes(await sha3KmacGo.readAsBytes());
    await File('${tempDir.path}/blake3_keyedmode.go').writeAsBytes(await blake3KeyedGo.readAsBytes());
    
    try {
      // Run go mod download
      await Process.run('go', ['mod', 'download'], workingDirectory: tempDir.path);
      
      // Run the Go program
      final result = await Process.run(
        'go',
        ['run', '.', algorithmIndex.toString(), keyHex, dataHex],
        workingDirectory: tempDir.path,
      );
      
      if (result.exitCode != 0) {
        throw Exception('Go execution failed: ${result.stderr}');
      }
      
      return result.stdout.toString().trim();
    } finally {
      // Cleanup
      await tempDir.delete(recursive: true);
    }
  }
  
  /// Get the name of the algorithm
  String get algorithmName => algorithm.name;
  
  String _bytesToHex(Uint8List bytes) {
    return bytes.map((b) => b.toRadixString(16).padLeft(2, '0')).join();
  }
  
  String _findProjectRoot() {
    // Navigate up from ports/dart to find project root
    var current = Directory.current.path;
    while (current != '/') {
      if (File('$current/go.mod').existsSync()) {
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
