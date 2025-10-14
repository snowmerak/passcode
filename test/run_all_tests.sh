#!/bin/bash

set -e

echo "========================================"
echo "Running All Cross-Implementation Tests"
echo "========================================"
echo ""

cd "$(dirname "$0")"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Track results
declare -A results

# 1. Go test
echo "1. Running Go implementation..."
if go run go_main.go > go_output.txt 2>&1; then
    cat go_output.txt
    results["Go"]="✅"
else
    echo -e "${RED}❌ Go test failed${NC}"
    cat go_output.txt
    results["Go"]="❌"
fi
echo ""

# 2. Rust test
echo "2. Running Rust implementation..."
if ~/.cargo/bin/cargo run --bin rust_test --quiet 2>/dev/null > rust_output.txt; then
    cat rust_output.txt
    results["Rust"]="✅"
else
    echo -e "${RED}❌ Rust test failed${NC}"
    cat rust_output.txt 2>/dev/null || echo "Build failed"
    results["Rust"]="❌"
fi
echo ""

# 3. Node.js (WASM) test
echo "3. Running Node.js (WASM) implementation..."
if node node_test.js > node_output.txt 2>&1; then
    cat node_output.txt
    results["Node.js-WASM"]="✅"
else
    echo -e "${RED}❌ Node.js test failed${NC}"
    cat node_output.txt
    results["Node.js-WASM"]="❌"
fi
echo ""

# 4. Node.js (Library) test
echo "4. Running Node.js (Library) implementation..."
if node nodejs_lib_test.js > nodejs_lib_output.txt 2>&1; then
    cat nodejs_lib_output.txt
    results["Node.js-Lib"]="✅"
else
    echo -e "${RED}❌ Node.js library test failed${NC}"
    cat nodejs_lib_output.txt
    results["Node.js-Lib"]="❌"
fi
echo ""

# 5. Python test
echo "5. Running Python implementation..."
if python3 python_test.py > python_output.txt 2>&1; then
    cat python_output.txt
    results["Python"]="✅"
else
    echo -e "${YELLOW}⚠️  Python test failed (may need wasmtime installation)${NC}"
    cat python_output.txt
    results["Python"]="⚠️"
fi
echo ""

# 6. Dart test
echo "6. Running Dart implementation..."
if command -v dart &> /dev/null; then
    if (cd ../ports/dart && dart run test/dart_test.dart) > dart_output.txt 2>&1; then
        cat dart_output.txt
        results["Dart"]="✅"
    else
        echo -e "${YELLOW}⚠️  Dart test failed (WASM bindings not fully implemented)${NC}"
        cat dart_output.txt
        results["Dart"]="⚠️"
    fi
else
    echo -e "${YELLOW}⚠️  Dart not installed, skipping...${NC}"
    results["Dart"]="⏭️"
fi
echo ""

# Compare outputs from working implementations
echo "========================================"
echo "Comparing Results"
echo "========================================"
echo ""

# Extract OTP values only (lines containing colons)
grep -E "^(SHA3|BLAKE3)" go_output.txt 2>/dev/null | sort > go_otps.txt || touch go_otps.txt
grep -E "^(SHA3|BLAKE3)" rust_output.txt 2>/dev/null | sort > rust_otps.txt || touch rust_otps.txt
grep -E "^(SHA3|BLAKE3)" node_output.txt 2>/dev/null | sort > node_otps.txt || touch node_otps.txt
grep -E "^(SHA3|BLAKE3)" nodejs_lib_output.txt 2>/dev/null | sort > nodejs_lib_otps.txt || touch nodejs_lib_otps.txt
grep -E "^(SHA3|BLAKE3)" python_output.txt 2>/dev/null | sort > python_otps.txt || touch python_otps.txt

all_match=true

if [ -s go_otps.txt ] && [ -s rust_otps.txt ] && [ -s node_otps.txt ] && [ -s nodejs_lib_otps.txt ]; then
    if diff -q go_otps.txt rust_otps.txt > /dev/null && \
       diff -q rust_otps.txt node_otps.txt > /dev/null && \
       diff -q node_otps.txt nodejs_lib_otps.txt > /dev/null; then
        echo -e "${GREEN}✅ SUCCESS: Core implementations (Go, Rust, Node.js) produce identical results!${NC}"
        echo ""
        echo "OTP Values:"
        cat go_otps.txt
        
        # Check Python if available
        if [ -s python_otps.txt ]; then
            if diff -q go_otps.txt python_otps.txt > /dev/null; then
                echo -e "\n${GREEN}✅ Python implementation also matches!${NC}"
            else
                echo -e "\n${YELLOW}⚠️  Python implementation produces different results${NC}"
                all_match=false
            fi
        fi
    else
        echo -e "${RED}❌ FAILURE: Core implementations produce different results!${NC}"
        all_match=false
    fi
else
    echo -e "${RED}❌ FAILURE: Not all core implementations ran successfully${NC}"
    all_match=false
fi

# Cleanup
rm -f go_output.txt rust_output.txt node_output.txt nodejs_lib_output.txt python_output.txt dart_output.txt
rm -f go_otps.txt rust_otps.txt node_otps.txt nodejs_lib_otps.txt python_otps.txt

echo ""
echo "========================================"
echo "Summary"
echo "========================================"
for impl in "Go" "Rust" "Node.js-WASM" "Node.js-Lib" "Python" "Dart"; do
    echo -e "$impl: ${results[$impl]}"
done
echo "========================================"

if [ "$all_match" = true ]; then
    exit 0
else
    exit 1
fi
