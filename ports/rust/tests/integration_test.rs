use passcode::{Algorithm, Passcode};
use passcode::{blake3_keyed_mode256, blake3_keyed_mode512};
use passcode::{sha3_kmac128, sha3_kmac256};
use rand::RngCore;

fn random_bytes(len: usize) -> Vec<u8> {
    let mut bytes = vec![0u8; len];
    rand::thread_rng().fill_bytes(&mut bytes);
    bytes
}

#[test]
fn test_passcode_creation() {
    let algorithms = [
        Algorithm::Sha3Kmac128,
        Algorithm::Sha3Kmac256,
        Algorithm::Blake3KeyedMode128,
        Algorithm::Blake3KeyedMode256,
    ];

    for algo in &algorithms {
        let key = random_bytes(32);
        let passcode = Passcode::new(*algo, key);
        assert_eq!(passcode.algorithm(), *algo);
    }
}

#[test]
fn test_otp_format() {
    let key = random_bytes(32);
    let challenge = random_bytes(16);
    let passcode = Passcode::new(Algorithm::Blake3KeyedMode256, key);
    
    let otp = passcode.compute(&challenge);
    
    // Should be 12 characters
    assert_eq!(otp.len(), 12);
    
    // Should be valid hex
    assert!(otp.chars().all(|c| c.is_ascii_hexdigit()));
}

#[test]
fn test_otp_consistency() {
    let key = random_bytes(32);
    let challenge = random_bytes(16);
    
    let passcode1 = Passcode::new(Algorithm::Blake3KeyedMode256, key.clone());
    let passcode2 = Passcode::new(Algorithm::Blake3KeyedMode256, key);
    
    let otp1 = passcode1.compute(&challenge);
    let otp2 = passcode2.compute(&challenge);
    
    assert_eq!(otp1, otp2);
}

#[test]
fn test_different_challenges() {
    let key = random_bytes(32);
    let challenge1 = random_bytes(16);
    let challenge2 = random_bytes(16);
    
    let passcode = Passcode::new(Algorithm::Blake3KeyedMode256, key);
    
    let otp1 = passcode.compute(&challenge1);
    let otp2 = passcode.compute(&challenge2);
    
    assert_ne!(otp1, otp2);
}

#[test]
fn test_different_keys() {
    let key1 = random_bytes(32);
    let key2 = random_bytes(32);
    let challenge = random_bytes(16);
    
    let passcode1 = Passcode::new(Algorithm::Blake3KeyedMode256, key1);
    let passcode2 = Passcode::new(Algorithm::Blake3KeyedMode256, key2);
    
    let otp1 = passcode1.compute(&challenge);
    let otp2 = passcode2.compute(&challenge);
    
    assert_ne!(otp1, otp2);
}

#[test]
fn test_empty_challenge() {
    let key = random_bytes(32);
    let challenge = vec![];
    
    let passcode = Passcode::new(Algorithm::Blake3KeyedMode256, key);
    let otp = passcode.compute(&challenge);
    
    assert_eq!(otp.len(), 12);
    assert!(otp.chars().all(|c| c.is_ascii_hexdigit()));
}

#[test]
fn test_different_algorithms() {
    let key = random_bytes(32);
    let challenge = random_bytes(16);
    
    let passcode1 = Passcode::new(Algorithm::Sha3Kmac256, key.clone());
    let passcode2 = Passcode::new(Algorithm::Blake3KeyedMode256, key);
    
    let otp1 = passcode1.compute(&challenge);
    let otp2 = passcode2.compute(&challenge);
    
    assert_ne!(otp1, otp2);
}

#[test]
fn test_blake3_output_lengths() {
    let key = random_bytes(32);
    let data = random_bytes(64);
    
    let hash256 = blake3_keyed_mode256(&key, &data);
    assert_eq!(hash256.len(), 32);
    
    let hash512 = blake3_keyed_mode512(&key, &data);
    assert_eq!(hash512.len(), 64);
}

#[test]
fn test_blake3_consistency() {
    let key = random_bytes(32);
    let data = random_bytes(64);
    
    let hash1 = blake3_keyed_mode256(&key, &data);
    let hash2 = blake3_keyed_mode256(&key, &data);
    
    assert_eq!(hash1, hash2);
}

#[test]
fn test_blake3_different_keys() {
    let key1 = random_bytes(32);
    let key2 = random_bytes(32);
    let data = random_bytes(64);
    
    let hash1 = blake3_keyed_mode256(&key1, &data);
    let hash2 = blake3_keyed_mode256(&key2, &data);
    
    assert_ne!(hash1, hash2);
}

#[test]
fn test_sha3_kmac_output_lengths() {
    let key = random_bytes(32);
    let customization = b"test";
    let data = random_bytes(64);
    
    let hash128 = sha3_kmac128(&key, customization, &data, 32);
    assert_eq!(hash128.len(), 32);
    
    let hash256 = sha3_kmac256(&key, customization, &data, 64);
    assert_eq!(hash256.len(), 64);
}

#[test]
fn test_sha3_kmac_consistency() {
    let key = random_bytes(32);
    let customization = b"test";
    let data = random_bytes(64);
    
    let hash1 = sha3_kmac256(&key, customization, &data, 32);
    let hash2 = sha3_kmac256(&key, customization, &data, 32);
    
    assert_eq!(hash1, hash2);
}

#[test]
fn test_sha3_kmac_different_keys() {
    let key1 = random_bytes(32);
    let key2 = random_bytes(32);
    let customization = b"test";
    let data = random_bytes(64);
    
    let hash1 = sha3_kmac256(&key1, customization, &data, 32);
    let hash2 = sha3_kmac256(&key2, customization, &data, 32);
    
    assert_ne!(hash1, hash2);
}

#[test]
fn test_sha3_kmac_different_customization() {
    let key = random_bytes(32);
    let customization1 = b"custom1";
    let customization2 = b"custom2";
    let data = random_bytes(64);
    
    let hash1 = sha3_kmac256(&key, customization1, &data, 32);
    let hash2 = sha3_kmac256(&key, customization2, &data, 32);
    
    assert_ne!(hash1, hash2);
}

#[test]
fn test_algorithm_display() {
    assert_eq!(Algorithm::Sha3Kmac128.to_string(), "SHA3-KMAC-128");
    assert_eq!(Algorithm::Sha3Kmac256.to_string(), "SHA3-KMAC-256");
    assert_eq!(Algorithm::Blake3KeyedMode128.to_string(), "BLAKE3-Keyed-Mode-128");
    assert_eq!(Algorithm::Blake3KeyedMode256.to_string(), "BLAKE3-Keyed-Mode-256");
}
