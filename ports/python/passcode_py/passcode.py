"""
Python wrapper for Passcode WASM library using wasmtime
"""

from enum import IntEnum
from pathlib import Path
from wasmtime import Store, Module, Instance


class Algorithm(IntEnum):
    """Supported hash algorithms"""
    SHA3_KMAC_128 = 0
    SHA3_KMAC_256 = 1
    BLAKE3_KEYED_MODE_128 = 2
    BLAKE3_KEYED_MODE_256 = 3


class WasmRuntime:
    """Manages WASM instance and memory operations"""
    
    def __init__(self):
        wasm_path = Path(__file__).parent / "passcode_wasm_bg.wasm"
        self.store = Store()
        module = Module.from_file(self.store.engine, str(wasm_path))
        self.instance = Instance(self.store, module, [])
        
        # Cache exports
        exports = self.instance.exports(self.store)
        self.memory = exports["memory"]
        self._alloc = exports["__wbindgen_malloc"]
        self._dealloc = exports["__wbindgen_free"]
        self._passcode_new = exports["passcode_new"]
        self._passcode_compute = exports["passcode_compute"]
        self._blake3_128 = exports["blake3KeyedMode128"]
        self._blake3_256 = exports["blake3KeyedMode256"]
        self._sha3_128 = exports["sha3Kmac128"]
        self._sha3_256 = exports["sha3Kmac256"]
    
    def write_bytes(self, data: bytes) -> tuple[int, int]:
        """Write bytes to WASM memory, return (ptr, len)"""
        length = len(data)
        ptr = self._alloc(self.store, length, 1)
        memory_data = self.memory.data_ptr(self.store)
        for i, byte in enumerate(data):
            memory_data[ptr + i] = byte
        return ptr, length
    
    def free_bytes(self, ptr: int, length: int):
        """Free WASM memory"""
        if ptr > 0:
            self._dealloc(self.store, ptr, length, 1)
    
    def read_string(self, ptr: int) -> str:
        """Read wasm-bindgen encoded string (ptr to [data_ptr, length])"""
        memory_data = self.memory.data_ptr(self.store)
        
        # Read string pointer and length (both i32, little endian)
        str_ptr = int.from_bytes(memory_data[ptr:ptr+4], 'little')
        str_len = int.from_bytes(memory_data[ptr+4:ptr+8], 'little')
        
        # Read actual string bytes
        string_bytes = bytes(memory_data[str_ptr:str_ptr + str_len])
        return string_bytes.decode('utf-8')


# Global WASM runtime instance
_runtime = WasmRuntime()


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
        
        # Create instance in WASM
        key_ptr, key_len = _runtime.write_bytes(key)
        try:
            self._handle = _runtime._passcode_new(
                _runtime.store, 
                int(algorithm), 
                key_ptr, 
                key_len
            )
        finally:
            _runtime.free_bytes(key_ptr, key_len)
    
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
        
        data_ptr, data_len = _runtime.write_bytes(data)
        try:
            result_ptr = _runtime._passcode_compute(
                _runtime.store, 
                self._handle, 
                data_ptr, 
                data_len
            )
            return _runtime.read_string(result_ptr)
        finally:
            _runtime.free_bytes(data_ptr, data_len)
    
    def algorithm_name(self) -> str:
        """Get the name of the algorithm"""
        names = {
            Algorithm.SHA3_KMAC_128: "SHA3-KMAC-128",
            Algorithm.SHA3_KMAC_256: "SHA3-KMAC-256",
            Algorithm.BLAKE3_KEYED_MODE_128: "BLAKE3-Keyed-Mode-128",
            Algorithm.BLAKE3_KEYED_MODE_256: "BLAKE3-Keyed-Mode-256",
        }
        return names.get(self.algorithm, "Unknown")


# Utility functions
def blake3_keyed_mode_128(key: bytes, data: bytes) -> str:
    """Compute BLAKE3-128 hash"""
    if not isinstance(key, bytes) or not isinstance(data, bytes):
        raise TypeError("key and data must be bytes")
    
    key_ptr, key_len = _runtime.write_bytes(key)
    data_ptr, data_len = _runtime.write_bytes(data)
    try:
        result_ptr = _runtime._blake3_128(
            _runtime.store, 
            key_ptr, key_len, 
            data_ptr, data_len
        )
        return _runtime.read_string(result_ptr)
    finally:
        _runtime.free_bytes(key_ptr, key_len)
        _runtime.free_bytes(data_ptr, data_len)


def blake3_keyed_mode_256(key: bytes, data: bytes) -> str:
    """Compute BLAKE3-256 hash"""
    if not isinstance(key, bytes) or not isinstance(data, bytes):
        raise TypeError("key and data must be bytes")
    
    key_ptr, key_len = _runtime.write_bytes(key)
    data_ptr, data_len = _runtime.write_bytes(data)
    try:
        result_ptr = _runtime._blake3_256(
            _runtime.store, 
            key_ptr, key_len, 
            data_ptr, data_len
        )
        return _runtime.read_string(result_ptr)
    finally:
        _runtime.free_bytes(key_ptr, key_len)
        _runtime.free_bytes(data_ptr, data_len)


def sha3_kmac_128(key: bytes, customization: bytes, data: bytes) -> str:
    """Compute SHA3-KMAC-128 hash"""
    if not all(isinstance(x, bytes) for x in [key, customization, data]):
        raise TypeError("key, customization, and data must be bytes")
    
    key_ptr, key_len = _runtime.write_bytes(key)
    cust_ptr, cust_len = _runtime.write_bytes(customization)
    data_ptr, data_len = _runtime.write_bytes(data)
    try:
        result_ptr = _runtime._sha3_128(
            _runtime.store,
            key_ptr, key_len,
            cust_ptr, cust_len,
            data_ptr, data_len
        )
        return _runtime.read_string(result_ptr)
    finally:
        _runtime.free_bytes(key_ptr, key_len)
        _runtime.free_bytes(cust_ptr, cust_len)
        _runtime.free_bytes(data_ptr, data_len)


def sha3_kmac_256(key: bytes, customization: bytes, data: bytes) -> str:
    """Compute SHA3-KMAC-256 hash"""
    if not all(isinstance(x, bytes) for x in [key, customization, data]):
        raise TypeError("key, customization, and data must be bytes")
    
    key_ptr, key_len = _runtime.write_bytes(key)
    cust_ptr, cust_len = _runtime.write_bytes(customization)
    data_ptr, data_len = _runtime.write_bytes(data)
    try:
        result_ptr = _runtime._sha3_256(
            _runtime.store,
            key_ptr, key_len,
            cust_ptr, cust_len,
            data_ptr, data_len
        )
        return _runtime.read_string(result_ptr)
    finally:
        _runtime.free_bytes(key_ptr, key_len)
        _runtime.free_bytes(cust_ptr, cust_len)
        _runtime.free_bytes(data_ptr, data_len)
