use sha3::digest::{ExtendableOutput, Update, XofReader};
use sha3::CShake128;
use hex;

fn encode_string(data: &[u8]) -> Vec<u8> {
    let bit_len = (data.len() * 8) as u64;
    let mut encoded = left_encode(bit_len);
    encoded.extend_from_slice(data);
    encoded
}

fn left_encode(x: u64) -> Vec<u8> {
    if x == 0 {
        return vec![1, 0];
    }
    let mut temp = [0u8; 8];
    let mut val = x;
    for i in (0..8).rev() {
        temp[i] = (val & 0xff) as u8;
        val >>= 8;
    }
    let mut start_idx = 0;
    while start_idx < 8 && temp[start_idx] == 0 {
        start_idx += 1;
    }
    let n = 8 - start_idx;
    let mut result = vec![n as u8];
    result.extend_from_slice(&temp[start_idx..]);
    result
}

fn main() {
    let key = hex::decode("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
        .expect("Invalid hex");
    
    println!("Test 1: Just raw bytes");
    let mut hasher1 = CShake128::from_core(
        sha3::CShake128Core::new(b"KMACauthorization"),
    );
    hasher1.update(&key);
    let mut output1 = vec![0u8; 16];
    hasher1.finalize_xof().read(&mut output1);
    println!("  Raw 'KMACauthorization': {}", hex::encode(&output1));
    
    println!("
Test 2: With encode_string");
    let mut domain_sep = Vec::new();
    domain_sep.extend_from_slice(&encode_string(b"KMAC"));
    domain_sep.extend_from_slice(&encode_string(b"authorization"));
    let mut hasher2 = CShake128::from_core(
        sha3::CShake128Core::new(&domain_sep),
    );
    hasher2.update(&key);
    let mut output2 = vec![0u8; 16];
    hasher2.finalize_xof().read(&mut output2);
    println!("  With encode_string: {}", hex::encode(&output2));
}", hex::encode(&output));
}
