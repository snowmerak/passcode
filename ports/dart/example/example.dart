import 'dart:typed_data';
import 'dart:math';
import 'package:passcode/passcode.dart';

void main() {
  print('=== Passcode Dart Example ===\n');
  
  // Generate random key and challenge
  final random = Random.secure();
  final key = Uint8List.fromList(List.generate(32, (_) => random.nextInt(256)));
  final challenge = Uint8List.fromList(List.generate(16, (_) => random.nextInt(256)));
  
  print('Key: ${_toHex(key)}');
  print('Challenge: ${_toHex(challenge)}\n');
  
  // Test all algorithms
  final algorithms = [
    Algorithm.sha3Kmac128,
    Algorithm.sha3Kmac256,
    Algorithm.blake3KeyedMode128,
    Algorithm.blake3KeyedMode256,
  ];
  
  for (final algo in algorithms) {
    try {
      final passcode = Passcode(algo, key);
      print('${passcode.algorithmName}:');
      print('  Algorithm name: ${passcode.algorithmName}');
      
      // This will throw UnimplementedError until WASM bindings are added
      // final otp = passcode.compute(challenge);
      // print('  OTP: $otp');
      print('  Status: API ready, WASM bindings pending\n');
    } catch (e) {
      print('  Error: $e\n');
    }
  }
}

String _toHex(Uint8List bytes) {
  return bytes.map((b) => b.toRadixString(16).padLeft(2, '0')).join();
}
