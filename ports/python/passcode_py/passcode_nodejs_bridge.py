#!/usr/bin/env python3
"""
Python wrapper for Passcode using Node.js bridge
This is a simpler approach that uses the Node.js WASM implementation
"""

import subprocess
import json
from enum import IntEnum
from pathlib import Path


class Algorithm(IntEnum):
    """Supported hash algorithms"""
    SHA3_KMAC_128 = 0
    SHA3_KMAC_256 = 1
    BLAKE3_KEYED_MODE_128 = 2
    BLAKE3_KEYED_MODE_256 = 3


class Passcode:
    """Challenge-response OTP generator"""
    
    def __init__(self, algorithm: Algorithm, key: bytes):
        """
        Create a new Passcode instance
        
        Args:
            algorithm: Algorithm to use (Algorithm enum)
            key: Secret key (bytes, 32 bytes recommended)
        """
        if not isinstance(key, bytes):
            raise TypeError("key must be bytes")
        
        self.algorithm = algorithm
        self.key = key
    
    def compute(self, data: bytes) -> str:
        """
        Compute OTP from challenge data via Node.js
        
        Args:
            data: Challenge data (bytes)
            
        Returns:
            12-character hexadecimal OTP string
        """
        if not isinstance(data, bytes):
            raise TypeError("data must be bytes")
        
        # Call Node.js script
        script = f"""
const {{ Passcode, Algorithm }} = require('{Path(__file__).parent.parent.parent / "nodejs"}');
const key = Buffer.from('{self.key.hex()}', 'hex');
const data = Buffer.from('{data.hex()}', 'hex');
const passcode = new Passcode({int(self.algorithm)}, key);
console.log(passcode.compute(data));
"""
        
        result = subprocess.run(
            ['node', '-e', script],
            capture_output=True,
            text=True,
            check=True
        )
        
        return result.stdout.strip()
    
    def algorithm_name(self) -> str:
        """Get the name of the algorithm"""
        names = {
            Algorithm.SHA3_KMAC_128: "SHA3-KMAC-128",
            Algorithm.SHA3_KMAC_256: "SHA3-KMAC-256",
            Algorithm.BLAKE3_KEYED_MODE_128: "BLAKE3-Keyed-Mode-128",
            Algorithm.BLAKE3_KEYED_MODE_256: "BLAKE3-Keyed-Mode-256",
        }
        return names.get(self.algorithm, "Unknown")
