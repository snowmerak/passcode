import { Passcode, Algorithm } from '../src';

// Example: Challenge-Response OTP Authentication

async function main() {
  console.log('=== Passcode TypeScript Example ===\n');

  // 1. Shared secret key between server and client (32 bytes recommended)
  const secretKey = new Uint8Array(32);
  crypto.getRandomValues(secretKey);
  console.log('Secret Key:', Buffer.from(secretKey).toString('hex'), '\n');

  // 2. Create a new Passcode instance
  // Available algorithms:
  // - Algorithm.SHA3_KMAC_128
  // - Algorithm.SHA3_KMAC_256
  // - Algorithm.BLAKE3_KEYED_MODE_128
  // - Algorithm.BLAKE3_KEYED_MODE_256
  const passcode = new Passcode(Algorithm.BLAKE3_KEYED_MODE_256, secretKey);
  console.log('Using algorithm:', passcode.getAlgorithm(), '\n');

  // 3. Server generates a challenge (a random value)
  const challenge = new Uint8Array(16);
  crypto.getRandomValues(challenge);
  console.log('Challenge:', Buffer.from(challenge).toString('hex'));

  // 4. Both server and client compute the OTP based on the challenge
  const serverOTP = passcode.compute(challenge);
  console.log('Server computed OTP:', serverOTP);

  // Client receives the challenge and computes the OTP
  const clientPasscode = new Passcode(Algorithm.BLAKE3_KEYED_MODE_256, secretKey);
  const clientOTP = clientPasscode.compute(challenge);
  console.log('Client computed OTP:', clientOTP);

  // 5. Server verifies the OTP submitted by the client
  if (serverOTP === clientOTP) {
    console.log('\n✅ Authentication successful!');
  } else {
    console.log('\n❌ Authentication failed!');
  }

  // Demonstrate different algorithms
  console.log('\n=== Testing Different Algorithms ===\n');

  const algorithms = [
    Algorithm.SHA3_KMAC_128,
    Algorithm.SHA3_KMAC_256,
    Algorithm.BLAKE3_KEYED_MODE_128,
    Algorithm.BLAKE3_KEYED_MODE_256,
  ];

  for (const algo of algorithms) {
    const pc = new Passcode(algo, secretKey);
    const otp = pc.compute(challenge);
    console.log(`${algo}: ${otp}`);
  }

  // Demonstrate challenge-response workflow
  console.log('\n=== Challenge-Response Workflow ===\n');

  const workflows = [
    { name: 'Attempt 1', challengeData: new Uint8Array(16) },
    { name: 'Attempt 2', challengeData: new Uint8Array(16) },
    { name: 'Attempt 3', challengeData: new Uint8Array(16) },
  ];

  workflows.forEach((workflow) => {
    crypto.getRandomValues(workflow.challengeData);
    const otp = passcode.compute(workflow.challengeData);
    console.log(`${workflow.name}:`, {
      challenge: Buffer.from(workflow.challengeData).toString('hex').substring(0, 16) + '...',
      otp: otp,
    });
  });
}

main().catch(console.error);
