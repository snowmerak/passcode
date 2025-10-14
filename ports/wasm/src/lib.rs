use wasm_bindgen::prelude::*;
use passcode::{Algorithm as RustAlgorithm, Passcode as RustPasscode};

/// Algorithm enum for WASM
#[wasm_bindgen]
#[derive(Debug, Clone, Copy)]
pub enum Algorithm {
    Sha3Kmac128,
    Sha3Kmac256,
    Blake3KeyedMode128,
    Blake3KeyedMode256,
}

impl From<Algorithm> for RustAlgorithm {
    fn from(algo: Algorithm) -> Self {
        match algo {
            Algorithm::Sha3Kmac128 => RustAlgorithm::Sha3Kmac128,
            Algorithm::Sha3Kmac256 => RustAlgorithm::Sha3Kmac256,
            Algorithm::Blake3KeyedMode128 => RustAlgorithm::Blake3KeyedMode128,
            Algorithm::Blake3KeyedMode256 => RustAlgorithm::Blake3KeyedMode256,
        }
    }
}

/// Passcode struct for WASM
#[wasm_bindgen]
pub struct Passcode {
    inner: RustPasscode,
}

#[wasm_bindgen]
impl Passcode {
    /// Creates a new Passcode instance
    ///
    /// # Arguments
    /// * `algorithm` - The hash algorithm to use
    /// * `key` - The secret key as a Uint8Array
    #[wasm_bindgen(constructor)]
    pub fn new(algorithm: Algorithm, key: &[u8]) -> Result<Passcode, JsValue> {
        let rust_algo = algorithm.into();
        let inner = RustPasscode::new(rust_algo, key.to_vec());
        
        Ok(Passcode { inner })
    }

    /// Computes an OTP from the given challenge data
    ///
    /// # Arguments
    /// * `data` - The challenge data as a Uint8Array
    ///
    /// # Returns
    /// A 12-character hexadecimal OTP string
    #[wasm_bindgen]
    pub fn compute(&self, data: &[u8]) -> String {
        self.inner.compute(data)
    }

    /// Gets the algorithm name as a string
    #[wasm_bindgen(getter, js_name = algorithmName)]
    pub fn algorithm_name(&self) -> String {
        self.inner.algorithm_name().to_string()
    }
}

/// Utility function: BLAKE3 keyed mode with 256-bit output
#[wasm_bindgen(js_name = blake3KeyedMode256)]
pub fn blake3_keyed_mode256(key: &[u8], data: &[u8]) -> Vec<u8> {
    passcode::blake3_keyed_mode256(key, data)
}

/// Utility function: BLAKE3 keyed mode with 512-bit output
#[wasm_bindgen(js_name = blake3KeyedMode512)]
pub fn blake3_keyed_mode512(key: &[u8], data: &[u8]) -> Vec<u8> {
    passcode::blake3_keyed_mode512(key, data)
}

/// Utility function: SHA3-KMAC128
#[wasm_bindgen(js_name = sha3Kmac128)]
pub fn sha3_kmac128(key: &[u8], customization: &[u8], data: &[u8], output_len: usize) -> Vec<u8> {
    passcode::sha3_kmac128(key, customization, data, output_len)
}

/// Utility function: SHA3-KMAC256
#[wasm_bindgen(js_name = sha3Kmac256)]
pub fn sha3_kmac256(key: &[u8], customization: &[u8], data: &[u8], output_len: usize) -> Vec<u8> {
    passcode::sha3_kmac256(key, customization, data, output_len)
}

// Module initialization for better error messages
#[wasm_bindgen(start)]
pub fn main() {
    #[cfg(feature = "console_error_panic_hook")]
    console_error_panic_hook::set_once();
}
