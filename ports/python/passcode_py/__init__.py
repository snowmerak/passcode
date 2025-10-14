"""
Passcode - Challenge-response OTP library using SHA3-KMAC and BLAKE3
"""

try:
    from .passcode import Algorithm, Passcode, blake3_keyed_mode_128, blake3_keyed_mode_256, sha3_kmac_128, sha3_kmac_256
except Exception:
    # Fallback to Node.js bridge if wasmtime is not available
    from .passcode_nodejs_bridge import Algorithm, Passcode
    blake3_keyed_mode_128 = None
    blake3_keyed_mode_256 = None
    sha3_kmac_128 = None
    sha3_kmac_256 = None

__version__ = "1.0.0"
__all__ = [
    "Algorithm",
    "Passcode",
    "blake3_keyed_mode_128",
    "blake3_keyed_mode_256",
    "sha3_kmac_128",
    "sha3_kmac_256",
]
