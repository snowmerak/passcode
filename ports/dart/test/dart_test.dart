import 'dart:typed_data';
import 'dart:io';
import 'package:passcode/passcode.dart';

void main() async {
  print('=== Dart Implementation Test ===\n');
  
  // Fixed test vectors (same as other tests)
  final key = _hexToBytes('0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef');
  final challenge = _hexToBytes('fedcba9876543210fedcba9876543210');
  
  print('Key:       ${_bytesToHex(key)}');
  print('Challenge: ${_bytesToHex(challenge)}\n');
  
  final algorithms = [
    ('SHA3-KMAC-128', Algorithm.sha3Kmac128),
    ('SHA3-KMAC-256', Algorithm.sha3Kmac256),
    ('BLAKE3-Keyed-128', Algorithm.blake3KeyedMode128),
    ('BLAKE3-Keyed-256', Algorithm.blake3KeyedMode256),
  ];
  
  for (final (name, algo) in algorithms) {
    try {
      final passcode = Passcode(algo, key);
      final otp = passcode.compute(challenge);
      print('${name.padRight(20)}: $otp');
    } catch (e) {
      print('${name.padRight(20)}: Error - $e');
    }
  }
}

Uint8List _hexToBytes(String hex) {
  final bytes = <int>[];
  for (var i = 0; i < hex.length; i += 2) {
    bytes.add(int.parse(hex.substring(i, i + 2), radix: 16));
  }
  return Uint8List.fromList(bytes);
}

String _bytesToHex(Uint8List bytes) {
  return bytes.map((b) => b.toRadixString(16).padLeft(2, '0')).join();
}
