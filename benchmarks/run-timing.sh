#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")"

# Defaults
ITERATIONS=1000
FILTER=""
OUTFILE=""

usage() {
    cat <<EOF
Usage: ./run-timing.sh [options]

Runs large-template wall-clock timing benchmarks across all engines.
Unlike run.sh (which uses Go's testing.B micro-benchmarks), this measures
real execution time on production-sized templates.

Options:
  -n, --iterations N   Number of render iterations per engine (default: 1000)
  -f, --filter STR     Only run scenarios containing STR (e.g. "Nested", "Complex")
  -o, --output FILE    Save output to FILE
  -h, --help           Show this help

Examples:
  ./run-timing.sh                          # Run all scenarios, 1000 iterations
  ./run-timing.sh -n 500                   # 500 iterations
  ./run-timing.sh -f "Large Loop"          # Only the Large Loop scenario
  ./run-timing.sh -n 2000 -o timing.txt    # 2000 iterations, save output
EOF
    exit 0
}

while [[ $# -gt 0 ]]; do
    case "$1" in
        -n|--iterations) ITERATIONS="$2"; shift 2 ;;
        -f|--filter)     FILTER="$2"; shift 2 ;;
        -o|--output)     OUTFILE="$2"; shift 2 ;;
        -h|--help)       usage ;;
        *) echo "Unknown option: $1"; usage ;;
    esac
done

ARGS=(-n "$ITERATIONS")
if [[ -n "$FILTER" ]]; then
    ARGS+=(-filter "$FILTER")
fi

if [[ -n "$OUTFILE" ]]; then
    go run ./cmd/timing/ "${ARGS[@]}" | tee "$OUTFILE"
    echo ""
    echo "Output saved to $OUTFILE"
else
    go run ./cmd/timing/ "${ARGS[@]}"
fi
