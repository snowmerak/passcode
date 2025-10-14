// Node.js example for Passcode WASM
const { Passcode, Algorithm } = require('../pkg-node/passcode_wasm');
const crypto = require('crypto');

console.log('=== Passcode WASM Node.js Example ===\n');

// 1. Create secret key
const secretKey = crypto.randomBytes(32);
console.log('Secret Key:', secretKey.toString('hex'), '\n');

// 2. Create Passcode instance
const passcode = new Passcode(Algorithm.Blake3KeyedMode256, secretKey);
console.log('Algorithm:', passcode.algorithmName, '\n');

// 3. Generate challenge
const challenge = crypto.randomBytes(16);
console.log('Challenge:', challenge.toString('hex'));

// 4. Compute OTP
const otp = passcode.compute(challenge);
console.log('Generated OTP:', otp);

// 5. Verify (server side would do the same)
const serverPasscode = new Passcode(Algorithm.Blake3KeyedMode256, secretKey);
const serverOtp = serverPasscode.compute(challenge);

if (otp === serverOtp) {
    console.log('\n✅ Authentication successful!');
} else {
    console.log('\n❌ Authentication failed!');
}

// Test all algorithms
console.log('\n=== Testing All Algorithms ===\n');

const algorithms = [
    { name: 'SHA3-KMAC-128', algo: Algorithm.Sha3Kmac128 },
    { name: 'SHA3-KMAC-256', algo: Algorithm.Sha3Kmac256 },
    { name: 'BLAKE3-Keyed-Mode-128', algo: Algorithm.Blake3KeyedMode128 },
    { name: 'BLAKE3-Keyed-Mode-256', algo: Algorithm.Blake3KeyedMode256 },
];

for (const { name, algo } of algorithms) {
    const pc = new Passcode(algo, secretKey);
    const result = pc.compute(challenge);
    console.log(`${name}: ${result}`);
}
