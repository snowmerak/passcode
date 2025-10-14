/// Supported hash algorithms
enum Algorithm {
  /// SHA3-KMAC with 128-bit output
  sha3Kmac128(0),
  
  /// SHA3-KMAC with 256-bit output
  sha3Kmac256(1),
  
  /// BLAKE3 Keyed Mode with 128-bit output
  blake3KeyedMode128(2),
  
  /// BLAKE3 Keyed Mode with 256-bit output
  blake3KeyedMode256(3);

  const Algorithm(this.value);
  
  /// Numeric value for FFI
  final int value;
  
  /// Get algorithm name
  String get name {
    switch (this) {
      case Algorithm.sha3Kmac128:
        return 'SHA3-KMAC-128';
      case Algorithm.sha3Kmac256:
        return 'SHA3-KMAC-256';
      case Algorithm.blake3KeyedMode128:
        return 'BLAKE3-Keyed-Mode-128';
      case Algorithm.blake3KeyedMode256:
        return 'BLAKE3-Keyed-Mode-256';
    }
  }
}
