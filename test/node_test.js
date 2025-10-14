const { Passcode, Algorithm } = require('../ports/wasm/pkg-node/passcode_wasm');

console.log('=== Node.js (WASM) Implementation Test ===\n');

// Fixed test vectors (must match Go and Rust tests)
const keyHex = '0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef';
const challengeHex = 'fedcba9876543210fedcba9876543210';

const key = Buffer.from(keyHex, 'hex');
const challenge = Buffer.from(challengeHex, 'hex');

console.log(`Key:       ${keyHex}`);
console.log(`Challenge: ${challengeHex}\n`);

const algorithms = [
    { name: 'SHA3-KMAC-128', algo: Algorithm.Sha3Kmac128 },
    { name: 'SHA3-KMAC-256', algo: Algorithm.Sha3Kmac256 },
    { name: 'BLAKE3-Keyed-128', algo: Algorithm.Blake3KeyedMode128 },
    { name: 'BLAKE3-Keyed-256', algo: Algorithm.Blake3KeyedMode256 },
];

for (const { name, algo } of algorithms) {
    const passcode = new Passcode(algo, key);
    const otp = passcode.compute(challenge);
    console.log(`${name.padEnd(20)}: ${otp}`);
}
