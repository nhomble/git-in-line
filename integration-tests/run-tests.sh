#!/bin/bash
set -ex

cd "$(dirname "$0")"
BINARY="../gil"
REPO="alphonse-bark"

rm -f "$BINARY"
[ -f "$BINARY" ] || (cd .. && go build -o gil ./src)

run_test() {
    local name="$1"
    local expected_file="$2"
    shift 2
    tmp_output=$(mktemp)
    "$BINARY" "$REPO" "$@" > "$tmp_output"

    if diff -q "$tmp_output" "$expected_file" >/dev/null; then
        echo "✓ $name"
    else
        echo "✗ $name"
        echo "Expected:"
        cat "$expected_file"
        echo "Got:"
        cat "$tmp_output"
        exit 1
    fi
}

run_test "default (desc)" expected-desc.txt test-commits-shuffled.txt
run_test "ascending" expected-asc.txt test-commits-shuffled.txt --order asc  
run_test "short flags" expected-asc.txt test-commits-shuffled.txt -o asc
run_test "explicit branch" expected-desc.txt test-commits-shuffled.txt --branch main --order desc
run_test "short hashes" expected-short-desc.txt test-commits-short.txt

echo "all tests passed"
