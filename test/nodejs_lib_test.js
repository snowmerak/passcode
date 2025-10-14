const { Passcode, Algorithm } = require('../ports/nodejs');
const crypto = require('crypto');

console.log('=== Node.js Library Test ===\n');

// Fixed test vectors (same as other tests)
const key = Buffer.from('0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef', 'hex');
const challenge = Buffer.from('fedcba9876543210fedcba9876543210', 'hex');

console.log(`Key:       ${key.toString('hex')}`);
console.log(`Challenge: ${challenge.toString('hex')}\n`);

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
