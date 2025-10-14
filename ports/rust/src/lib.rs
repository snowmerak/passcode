//! # Passcode
//!
//! Challenge-Response based One-Time Password (OTP) authentication library.
//!
//! This library provides a secure Challenge-Response authentication mechanism
//! using various hash algorithms including SHA3-KMAC and BLAKE3 Keyed Mode.
//!
//! ## Features
//!
//! - **Challenge-Response Mechanism**: Secure authentication where the server sends a random challenge
//! - **Multiple Hash Algorithms**: SHA3-KMAC (128/256) and BLAKE3 Keyed Mode (128/256)
//! - **Flexible Security Levels**: Choose between 128-bit and 256-bit security strengths
//! - **Type-Safe API**: Leverages Rust's type system for safety
//!
//! ## Example
//!
//! ```
//! use passcode::{Passcode, Algorithm};
//!
//! // 1. Shared secret key between server and client
//! let secret_key = vec![0u8; 32];
//!
//! // 2. Create a new Passcode instance
//! let passcode = Passcode::new(Algorithm::Blake3KeyedMode256, secret_key);
//!
//! // 3. Generate a challenge
//! let challenge = vec![0u8; 16];
//!
//! // 4. Compute the OTP
//! let otp = passcode.compute(&challenge);
//! println!("Generated OTP: {}", otp);
//! ```

mod blake3_keyed;
mod passcode;
mod sha3_kmac;

pub use passcode::{Algorithm, Passcode};
pub use blake3_keyed::{blake3_keyed_mode256, blake3_keyed_mode512};
pub use sha3_kmac::{sha3_kmac128, sha3_kmac256};
