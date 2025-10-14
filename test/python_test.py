#!/usr/bin/env python3
"""Python test for Passcode library"""

import sys
from pathlib import Path

# Add parent directory to path
sys.path.insert(0, str(Path(__file__).parent.parent / "ports" / "python"))

from passcode_py import Passcode, Algorithm

def main():
    print("=== Python Implementation Test ===\n")
    
    # Fixed test vectors (same as other tests)
    key = bytes.fromhex("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
    challenge = bytes.fromhex("fedcba9876543210fedcba9876543210")
    
    print(f"Key:       {key.hex()}")
    print(f"Challenge: {challenge.hex()}\n")
    
    algorithms = [
        ("SHA3-KMAC-128", Algorithm.SHA3_KMAC_128),
        ("SHA3-KMAC-256", Algorithm.SHA3_KMAC_256),
        ("BLAKE3-Keyed-128", Algorithm.BLAKE3_KEYED_MODE_128),
        ("BLAKE3-Keyed-256", Algorithm.BLAKE3_KEYED_MODE_256),
    ]
    
    for name, algo in algorithms:
        passcode = Passcode(algo, key)
        otp = passcode.compute(challenge)
        print(f"{name:<20}: {otp}")

if __name__ == "__main__":
    main()
