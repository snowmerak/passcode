//! FFI bindings for C/Dart interop

use std::slice;
use crate::{Algorithm, Passcode};

/// Create a new Passcode instance
/// Returns a pointer to the Passcode instance
#[no_mangle]
pub extern "C" fn passcode_new(algorithm: u8, key_ptr: *const u8, key_len: usize) -> *mut Passcode {
    let key = unsafe { slice::from_raw_parts(key_ptr, key_len) }.to_vec();
    let algo = match algorithm {
        0 => Algorithm::Sha3Kmac128,
        1 => Algorithm::Sha3Kmac256,
        2 => Algorithm::Blake3KeyedMode128,
        3 => Algorithm::Blake3KeyedMode256,
        _ => return std::ptr::null_mut(),
    };
    
    Box::into_raw(Box::new(Passcode::new(algo, key)))
}

/// Compute OTP from challenge data
/// Returns a pointer to a null-terminated string (caller must free)
#[no_mangle]
pub extern "C" fn passcode_compute(
    passcode_ptr: *mut Passcode,
    data_ptr: *const u8,
    data_len: usize,
    out_ptr: *mut u8,
    out_len: usize,
) -> i32 {
    if passcode_ptr.is_null() || data_ptr.is_null() || out_ptr.is_null() {
        return -1;
    }
    
    let passcode = unsafe { &*passcode_ptr };
    let data = unsafe { slice::from_raw_parts(data_ptr, data_len) };
    
    let result = passcode.compute(data);
    let result_bytes = result.as_bytes();
    
    if result_bytes.len() >= out_len {
        return -2; // Buffer too small
    }
    
    unsafe {
        std::ptr::copy_nonoverlapping(
            result_bytes.as_ptr(),
            out_ptr,
            result_bytes.len(),
        );
        *out_ptr.add(result_bytes.len()) = 0; // Null terminator
    }
    
    result_bytes.len() as i32
}

/// Free a Passcode instance
#[no_mangle]
pub extern "C" fn passcode_free(passcode_ptr: *mut Passcode) {
    if !passcode_ptr.is_null() {
        unsafe {
            let _ = Box::from_raw(passcode_ptr);
        }
    }
}

/// Get the last error message
#[no_mangle]
pub extern "C" fn passcode_get_error() -> *const u8 {
    b"No error\0".as_ptr()
}
