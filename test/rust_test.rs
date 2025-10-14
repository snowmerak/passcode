use passcode::{Algorithm, Passcode};

fn main() {
    println!("=== Rust Implementation Test ===\n");

    // Fixed test vectors
    let key = hex::decode("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
        .expect("Invalid key hex");
    let challenge = hex::decode("fedcba9876543210fedcba9876543210")
        .expect("Invalid challenge hex");

    println!("Key:       {}", hex::encode(&key));
    println!("Challenge: {}\n", hex::encode(&challenge));

    let algorithms = [
        ("SHA3-KMAC-128", Algorithm::Sha3Kmac128),
        ("SHA3-KMAC-256", Algorithm::Sha3Kmac256),
        ("BLAKE3-Keyed-128", Algorithm::Blake3KeyedMode128),
        ("BLAKE3-Keyed-256", Algorithm::Blake3KeyedMode256),
    ];

    for (name, algo) in algorithms.iter() {
        let passcode = Passcode::new(*algo, key.clone());
        let otp = passcode.compute(&challenge);
        println!("{:<20}: {}", name, otp);
    }
}
