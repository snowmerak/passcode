// Node.js wrapper for Passcode WASM
const wasm = require('./wasm/passcode_wasm.js');

// Algorithm enum matching WASM
const Algorithm = {
  Sha3Kmac128: 0,
  Sha3Kmac256: 1,
  Blake3KeyedMode128: 2,
  Blake3KeyedMode256: 3,
};

// Utility functions
function blake3KeyedMode128(key, data) {
  return wasm.blake3KeyedMode128(key, data);
}

function blake3KeyedMode256(key, data) {
  return wasm.blake3KeyedMode256(key, data);
}

function sha3Kmac128(key, customization, data) {
  return wasm.sha3Kmac128(key, customization, data);
}

function sha3Kmac256(key, customization, data) {
  return wasm.sha3Kmac256(key, customization, data);
}

// Main Passcode class
class Passcode {
  constructor(algorithm, key) {
    this.inner = new wasm.Passcode(algorithm, key);
  }

  compute(data) {
    return this.inner.compute(data);
  }

  algorithmName() {
    return this.inner.algorithmName();
  }
}

module.exports = {
  Algorithm,
  Passcode,
  blake3KeyedMode128,
  blake3KeyedMode256,
  sha3Kmac128,
  sha3Kmac256,
};
