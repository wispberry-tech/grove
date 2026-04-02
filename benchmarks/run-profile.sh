#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "$0")"

# Defaults
ITERATIONS=5000
SCENARIO="all"
OUTDIR="profiles"

usage() {
    cat <<EOF
Usage: ./run-profile.sh [options]

Profiles Grove rendering with CPU/memory pprof and per-opcode timing.
Builds with -tags groveprofile for opcode-level instrumentation.

Options:
  -n, --iterations N   Iterations per scenario (default: 5000)
  -s, --scenario S     Scenario: simple|loop|conditional|complex|large|all (default: all)
  -o, --outdir DIR     Output directory for profiles (default: profiles)
  -h, --help           Show this help

Examples:
  ./run-profile.sh                        # Profile all scenarios
  ./run-profile.sh -s loop -n 10000       # Focus on loop performance
  ./run-profile.sh -s large               # Profile large templates only

After profiling, explore interactively:
  go tool pprof -http=:8080 profiles/cpu.prof
EOF
    exit 0
}

while [[ $# -gt 0 ]]; do
    case "$1" in
        -n|--iterations) ITERATIONS="$2"; shift 2 ;;
        -s|--scenario)   SCENARIO="$2"; shift 2 ;;
        -o|--outdir)     OUTDIR="$2"; shift 2 ;;
        -h|--help)       usage ;;
        *) echo "Unknown option: $1"; usage ;;
    esac
done

mkdir -p "$OUTDIR"

echo "Building with opcode instrumentation..."
echo ""

go run -tags groveprofile ./cmd/profile/ \
    -n "$ITERATIONS" \
    -scenario "$SCENARIO" \
    -cpuprofile "$OUTDIR/cpu.prof" \
    -memprofile "$OUTDIR/mem.prof"

echo ""
echo "════════════════════════════════════════════════════════"
echo "  Top CPU consumers"
echo "════════════════════════════════════════════════════════"
echo ""
go tool pprof -top -nodecount=15 "$OUTDIR/cpu.prof" 2>/dev/null || true

echo ""
echo "════════════════════════════════════════════════════════"
echo "  Top memory allocators"
echo "════════════════════════════════════════════════════════"
echo ""
go tool pprof -top -nodecount=15 -sample_index=alloc_space "$OUTDIR/mem.prof" 2>/dev/null || true

echo ""
echo "Interactive exploration:"
echo "  go tool pprof -http=:8080 $OUTDIR/cpu.prof"
echo "  go tool pprof -http=:8080 $OUTDIR/mem.prof"
