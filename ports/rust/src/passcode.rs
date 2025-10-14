use crate::blake3_keyed::{blake3_keyed_mode256, blake3_keyed_mode512};
use crate::sha3_kmac::{sha3_kmac128_for_passcode, sha3_kmac256_for_passcode};

/// Available hash algorithms for OTP generation
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum Algorithm {
    /// SHA3-KMAC with 128-bit security
    Sha3Kmac128,
    /// SHA3-KMAC with 256-bit security
    Sha3Kmac256,
    /// BLAKE3 Keyed Mode with 128-bit security
    Blake3KeyedMode128,
    /// BLAKE3 Keyed Mode with 256-bit security
    Blake3KeyedMode256,
}

impl Algorithm {
    /// Returns the algorithm name as a string
    pub fn as_str(&self) -> &'static str {
        match self {
            Algorithm::Sha3Kmac128 => "SHA3-KMAC-128",
            Algorithm::Sha3Kmac256 => "SHA3-KMAC-256",
            Algorithm::Blake3KeyedMode128 => "BLAKE3-Keyed-Mode-128",
            Algorithm::Blake3KeyedMode256 => "BLAKE3-Keyed-Mode-256",
        }
    }
}

impl std::fmt::Display for Algorithm {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "{}", self.as_str())
    }
}

/// Hasher function type
type Hasher = fn(&[u8], &[u8]) -> Vec<u8>;

/// Passcode struct for Challenge-Response based OTP authentication
pub struct Passcode {
    algorithm: Algorithm,
    key: Vec<u8>,
    hasher: Hasher,
}

impl Passcode {
    /// Creates a new Passcode instance
    ///
    /// # Arguments
    /// * `algorithm` - The hash algorithm to use
    /// * `key` - The secret key (shared between server and client)
    ///
    /// # Example
    /// ```
    /// use passcode::{Passcode, Algorithm};
    ///
    /// let key = vec![0u8; 32];
    /// let passcode = Passcode::new(Algorithm::Blake3KeyedMode256, key);
    /// ```
    pub fn new(algorithm: Algorithm, key: Vec<u8>) -> Self {
        let hasher: Hasher = match algorithm {
            Algorithm::Sha3Kmac128 => sha3_kmac128_for_passcode,
            Algorithm::Sha3Kmac256 => sha3_kmac256_for_passcode,
            Algorithm::Blake3KeyedMode128 => blake3_keyed_mode256, // Using 256-bit output for 128-bit mode
            Algorithm::Blake3KeyedMode256 => blake3_keyed_mode512,
        };

        Self {
            algorithm,
            key,
            hasher,
        }
    }

    /// Computes an OTP from the given challenge data
    ///
    /// # Arguments
    /// * `data` - The challenge data (typically a random value from the server)
    ///
    /// # Returns
    /// A 12-character hexadecimal OTP string
    ///
    /// # Example
    /// ```
    /// use passcode::{Passcode, Algorithm};
    ///
    /// let key = vec![0u8; 32];
    /// let challenge = vec![0u8; 16];
    /// let passcode = Passcode::new(Algorithm::Blake3KeyedMode256, key);
    /// let otp = passcode.compute(&challenge);
    /// assert_eq!(otp.len(), 12);
    /// ```
    pub fn compute(&self, data: &[u8]) -> String {
        let mut hashed = (self.hasher)(&self.key, data);

        // Ensure we have at least 6 bytes
        if hashed.len() < 6 {
            hashed.resize(6, 0);
        }

        // Convert first 6 bytes to hex string
        hex::encode(&hashed[..6])
    }

    /// Gets the algorithm being used
    pub fn algorithm(&self) -> Algorithm {
        self.algorithm
    }

    /// Gets the algorithm name as a string
    pub fn algorithm_name(&self) -> &'static str {
        self.algorithm.as_str()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_new_passcode() {
        let key = vec![0u8; 32];
        let passcode = Passcode::new(Algorithm::Blake3KeyedMode256, key);
        assert_eq!(passcode.algorithm(), Algorithm::Blake3KeyedMode256);
    }

    #[test]
    fn test_compute_generates_12_char_hex() {
        let key = vec![0u8; 32];
        let challenge = vec![0u8; 16];
        let passcode = Passcode::new(Algorithm::Blake3KeyedMode256, key);
        let otp = passcode.compute(&challenge);
        
        assert_eq!(otp.len(), 12);
        assert!(otp.chars().all(|c| c.is_ascii_hexdigit()));
    }

    #[test]
    fn test_consistent_otp() {
        let key = vec![1u8; 32];
        let challenge = vec![2u8; 16];
        
        let passcode1 = Passcode::new(Algorithm::Blake3KeyedMode256, key.clone());
        let passcode2 = Passcode::new(Algorithm::Blake3KeyedMode256, key);
        
        let otp1 = passcode1.compute(&challenge);
        let otp2 = passcode2.compute(&challenge);
        
        assert_eq!(otp1, otp2);
    }

    #[test]
    fn test_different_challenges_different_otps() {
        let key = vec![1u8; 32];
        let challenge1 = vec![2u8; 16];
        let challenge2 = vec![3u8; 16];
        
        let passcode = Passcode::new(Algorithm::Blake3KeyedMode256, key);
        
        let otp1 = passcode.compute(&challenge1);
        let otp2 = passcode.compute(&challenge2);
        
        assert_ne!(otp1, otp2);
    }

    #[test]
    fn test_all_algorithms() {
        let key = vec![1u8; 32];
        let challenge = vec![2u8; 16];
        
        let algorithms = [
            Algorithm::Sha3Kmac128,
            Algorithm::Sha3Kmac256,
            Algorithm::Blake3KeyedMode128,
            Algorithm::Blake3KeyedMode256,
        ];
        
        for algo in &algorithms {
            let passcode = Passcode::new(*algo, key.clone());
            let otp = passcode.compute(&challenge);
            assert_eq!(otp.len(), 12);
        }
    }
}
