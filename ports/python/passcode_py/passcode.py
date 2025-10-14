"""
Python wrapper for Passcode WASM library
"""

import os
from enum import IntEnum
from pathlib import Path
from wasmer import engine, Store, Module, Instance
from wasmer_compiler_cranelift import Compiler

# Load WASM module
_wasm_path = Path(__file__).parent / "passcode_wasm_bg.wasm"
_store = Store(engine.JIT(Compiler))
_module = Module(_store, open(_wasm_path, 'rb').read())
_instance = Instance(_module)


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
        Compute OTP from challenge data
        
        Args:
            data: Challenge data (bytes)
            
        Returns:
            12-character hexadecimal OTP string
        """
        if not isinstance(data, bytes):
            raise TypeError("data must be bytes")
        
        # Call appropriate WASM function based on algorithm
        if self.algorithm == Algorithm.SHA3_KMAC_128:
            return sha3_kmac_128(self.key, b"authorization", data)
        elif self.algorithm == Algorithm.SHA3_KMAC_256:
            return sha3_kmac_256(self.key, b"authorization", data)
        elif self.algorithm == Algorithm.BLAKE3_KEYED_MODE_128:
            return blake3_keyed_mode_128(self.key, data)
        elif self.algorithm == Algorithm.BLAKE3_KEYED_MODE_256:
            return blake3_keyed_mode_256(self.key, data)
        else:
            raise ValueError(f"Unknown algorithm: {self.algorithm}")
    
    def algorithm_name(self) -> str:
        """Get the name of the algorithm"""
        names = {
            Algorithm.SHA3_KMAC_128: "SHA3-KMAC-128",
            Algorithm.SHA3_KMAC_256: "SHA3-KMAC-256",
            Algorithm.BLAKE3_KEYED_MODE_128: "BLAKE3-Keyed-Mode-128",
            Algorithm.BLAKE3_KEYED_MODE_256: "BLAKE3-Keyed-Mode-256",
        }
        return names.get(self.algorithm, "Unknown")


# Utility functions (placeholders - need WASM bindings)
def blake3_keyed_mode_128(key: bytes, data: bytes) -> str:
    """Compute BLAKE3-128 hash"""
    # TODO: Implement WASM binding
    raise NotImplementedError("WASM bindings not yet implemented")


def blake3_keyed_mode_256(key: bytes, data: bytes) -> str:
    """Compute BLAKE3-256 hash"""
    # TODO: Implement WASM binding
    raise NotImplementedError("WASM bindings not yet implemented")


def sha3_kmac_128(key: bytes, customization: bytes, data: bytes) -> str:
    """Compute SHA3-KMAC-128 hash"""
    # TODO: Implement WASM binding
    raise NotImplementedError("WASM bindings not yet implemented")


def sha3_kmac_256(key: bytes, customization: bytes, data: bytes) -> str:
    """Compute SHA3-KMAC-256 hash"""
    # TODO: Implement WASM binding
    raise NotImplementedError("WASM bindings not yet implemented")
