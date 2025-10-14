use sha3::digest::{ExtendableOutput, Update, XofReader};
use sha3::{CShake128, CShake256};

/// Left encode function for KMAC
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

    let mut result = Vec::with_capacity(n + 1);
    result.push(n as u8);
    result.extend_from_slice(&temp[start_idx..]);
    result
}

/// Right encode function for KMAC
fn right_encode(x: u64) -> Vec<u8> {
    if x == 0 {
        return vec![0, 1];
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

    let mut result = Vec::with_capacity(n + 1);
    result.extend_from_slice(&temp[start_idx..8]);
    result.push(n as u8);
    result
}

/// Encode a byte string with its bit length
fn encode_string(data: &[u8]) -> Vec<u8> {
    let bit_len = (data.len() * 8) as u64;
    let encoded = left_encode(bit_len);

    let mut result = Vec::with_capacity(encoded.len() + data.len());
    result.extend_from_slice(&encoded);
    result.extend_from_slice(data);
    result
}

/// Bytepad function for KMAC
fn bytepad(data: &[u8], w: usize) -> Vec<u8> {
    let w_encoded = left_encode(w as u64);
    let total_len = w_encoded.len() + data.len();

    let mut pad_len = w - (total_len % w);
    if pad_len == w {
        pad_len = 0;
    }

    let mut result = Vec::with_capacity(total_len + pad_len);
    result.extend_from_slice(&w_encoded);
    result.extend_from_slice(data);
    result.resize(total_len + pad_len, 0);
    result
}

/// KMAC implementation using CShake128
fn kmac128(
    key: &[u8],
    customization: &[u8],
    data: &[u8],
    output_len: usize,
) -> Vec<u8> {
    let encoded_key = encode_string(key);
    let padded_key = bytepad(&encoded_key, 168); // rate for SHA3-128

    // NIST SP 800-185: KMAC uses cSHAKE with function name "KMAC" and customization
    let mut hasher = CShake128::from_core(
        sha3::CShake128Core::new_with_function_name(b"KMAC", customization),
    );
    
    hasher.update(&padded_key);
    hasher.update(data);
    hasher.update(&right_encode((output_len * 8) as u64));

    let mut output = vec![0u8; output_len];
    hasher.finalize_xof().read(&mut output);
    output
}

/// KMAC implementation using CShake256
fn kmac256(
    key: &[u8],
    customization: &[u8],
    data: &[u8],
    output_len: usize,
) -> Vec<u8> {
    let encoded_key = encode_string(key);
    let padded_key = bytepad(&encoded_key, 136); // rate for SHA3-256

    // NIST SP 800-185: KMAC uses cSHAKE with function name "KMAC" and customization
    let mut hasher = CShake256::from_core(
        sha3::CShake256Core::new_with_function_name(b"KMAC", customization),
    );
    
    hasher.update(&padded_key);
    hasher.update(data);
    hasher.update(&right_encode((output_len * 8) as u64));

    let mut output = vec![0u8; output_len];
    hasher.finalize_xof().read(&mut output);
    output
}

/// SHA3-KMAC128 for passcode (internal use)
pub fn sha3_kmac128_for_passcode(key: &[u8], data: &[u8]) -> Vec<u8> {
    kmac128(key, b"authorization", data, 32)
}

/// SHA3-KMAC128 with customizable parameters
pub fn sha3_kmac128(
    key: &[u8],
    customization: &[u8],
    data: &[u8],
    output_len: usize,
) -> Vec<u8> {
    kmac128(key, customization, data, output_len)
}

/// SHA3-KMAC256 for passcode (internal use)
pub fn sha3_kmac256_for_passcode(key: &[u8], data: &[u8]) -> Vec<u8> {
    kmac256(key, b"authorization", data, 32)
}

/// SHA3-KMAC256 with customizable parameters
pub fn sha3_kmac256(
    key: &[u8],
    customization: &[u8],
    data: &[u8],
    output_len: usize,
) -> Vec<u8> {
    kmac256(key, customization, data, output_len)
}
