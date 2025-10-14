module test

go 1.25.2

replace github.com/snowmerak/passcode => ../

require github.com/snowmerak/passcode v0.0.0-00010101000000-000000000000

require (
	github.com/klauspost/cpuid/v2 v2.0.9 // indirect
	golang.org/x/crypto v0.43.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
	lukechampine.com/blake3 v1.4.1 // indirect
)
