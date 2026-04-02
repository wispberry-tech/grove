.PHONY: bench compare bench-grove clean

BENCH_FLAGS ?= -benchmem -count=6 -timeout=10m

# Run all benchmarks and save raw output.
bench:
	go test -bench=. -run='^$$' $(BENCH_FLAGS) | tee results.txt

# Format results with benchstat (install: go install golang.org/x/perf/cmd/benchstat@latest).
compare: bench
	@command -v benchstat >/dev/null 2>&1 || { echo "Install benchstat: go install golang.org/x/perf/cmd/benchstat@latest"; exit 1; }
	benchstat results.txt

# Run only Grove benchmarks for quick regression checks.
bench-grove:
	go test -bench=. -run='^$$' $(BENCH_FLAGS) -benchtime=3s 2>&1 | grep -E '(Benchmark|^ok|^FAIL)' | grep -E '(/Grove|^ok|^FAIL)' | tee grove-results.txt

clean:
	rm -f results.txt grove-results.txt
