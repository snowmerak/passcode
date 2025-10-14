use blake3::Hasher;

/// BLAKE3 keyed mode implementation
fn blake3_keyed_mode(key: &[u8], data: &[u8], out_len: usize) -> Vec<u8> {
    // Hash the key first to get a 32-byte key
    let hashed_key = blake3::hash(key);
    
    // Use BLAKE3 keyed hash with the hashed key
    let mut hasher = Hasher::new_keyed(hashed_key.as_bytes());
    hasher.update(data);
    
    // Get the output with specified length
    let mut output = vec![0u8; out_len];
    let mut reader = hasher.finalize_xof();
    reader.fill(&mut output);
    output
}

/// BLAKE3 keyed mode with 256-bit (32 bytes) output
pub fn blake3_keyed_mode256(key: &[u8], data: &[u8]) -> Vec<u8> {
    blake3_keyed_mode(key, data, 32)
}

/// BLAKE3 keyed mode with 512-bit (64 bytes) output
pub fn blake3_keyed_mode512(key: &[u8], data: &[u8]) -> Vec<u8> {
    blake3_keyed_mode(key, data, 64)
}
