use passcode::{Algorithm, Passcode};
use rand::RngCore;

fn main() {
    println!("=== Passcode Rust Example ===\n");

    let mut rng = rand::thread_rng();

    // 1. Shared secret key between server and client (32 bytes recommended)
    let mut secret_key = vec![0u8; 32];
    rng.fill_bytes(&mut secret_key);
    println!("Secret Key: {}\n", hex::encode(&secret_key));

    // 2. Create a new Passcode instance
    // Available algorithms:
    // - Algorithm::Sha3Kmac128
    // - Algorithm::Sha3Kmac256
    // - Algorithm::Blake3KeyedMode128
    // - Algorithm::Blake3KeyedMode256
    let passcode = Passcode::new(Algorithm::Blake3KeyedMode256, secret_key.clone());
    println!("Using algorithm: {}\n", passcode.algorithm_name());

    // 3. Server generates a challenge (a random value)
    let mut challenge = vec![0u8; 16];
    rng.fill_bytes(&mut challenge);
    println!("Challenge: {}", hex::encode(&challenge));

    // 4. Both server and client compute the OTP based on the challenge
    let server_otp = passcode.compute(&challenge);
    println!("Server computed OTP: {}", server_otp);

    // Client receives the challenge and computes the OTP
    let client_passcode = Passcode::new(Algorithm::Blake3KeyedMode256, secret_key.clone());
    let client_otp = client_passcode.compute(&challenge);
    println!("Client computed OTP: {}", client_otp);

    // 5. Server verifies the OTP submitted by the client
    if server_otp == client_otp {
        println!("\n✅ Authentication successful!");
    } else {
        println!("\n❌ Authentication failed!");
    }

    // Demonstrate different algorithms
    println!("\n=== Testing Different Algorithms ===\n");

    let algorithms = [
        Algorithm::Sha3Kmac128,
        Algorithm::Sha3Kmac256,
        Algorithm::Blake3KeyedMode128,
        Algorithm::Blake3KeyedMode256,
    ];

    for algo in &algorithms {
        let pc = Passcode::new(*algo, secret_key.clone());
        let otp = pc.compute(&challenge);
        println!("{}: {}", algo, otp);
    }

    // Demonstrate challenge-response workflow
    println!("\n=== Challenge-Response Workflow ===\n");

    for i in 1..=3 {
        let mut challenge_data = vec![0u8; 16];
        rng.fill_bytes(&mut challenge_data);
        let otp = passcode.compute(&challenge_data);
        println!(
            "Attempt {}: challenge={}..., otp={}",
            i,
            &hex::encode(&challenge_data)[..16],
            otp
        );
    }

    // Demonstrate using hash functions directly
    println!("\n=== Using Hash Functions Directly ===\n");

    let mut key = vec![0u8; 32];
    rng.fill_bytes(&mut key);
    let mut data = vec![0u8; 64];
    rng.fill_bytes(&mut data);

    // BLAKE3
    let blake3_hash = passcode::blake3_keyed_mode256(&key, &data);
    println!("BLAKE3-256: {}...", &hex::encode(&blake3_hash)[..32]);

    // SHA3-KMAC
    let customization = b"my-app";
    let kmac_hash = passcode::sha3_kmac256(&key, customization, &data, 32);
    println!("SHA3-KMAC256: {}...", &hex::encode(&kmac_hash)[..32]);
}
