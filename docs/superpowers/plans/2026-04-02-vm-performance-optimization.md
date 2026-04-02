# Grove VM Performance Optimization Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make Grove faster than Pongo2 on all benchmark scenarios by eliminating allocation hotspots and reducing per-opcode overhead in the VM.

**Architecture:** Six targeted optimizations to the VM execution loop and supporting types. Each addresses a specific bottleneck identified by profiling. Changes are confined to `internal/vm/`, `internal/scope/`, and `internal/compiler/` — the public API is unchanged.

**Tech Stack:** Go 1.24, no new dependencies

---

## Performance Baseline

```
BenchmarkRender_Loop/Grove      4186 ns/op   6757 B/op   60 allocs/op
BenchmarkRender_Loop/Pongo2     2539 ns/op   3194 B/op   69 allocs/op

BenchmarkRender_Complex/Grove   29373 ns/op  45460 B/op  284 allocs/op
BenchmarkRender_Complex/Pongo2  20915 ns/op  28077 B/op  541 allocs/op
```

**Profiler bottlenecks:**
- `makeLoopMap` — 35% of all memory (1555MB/5000 iters). Allocates `map[string]any` per loop iteration.
- `strings.Builder.WriteString` — 51% of memory. Output buffer regrowth.
- `FromAny` — 4.5% of memory + CPU. Type switch on every `scope.Get` result.
- Var Lookup — 31.8% of opcode time. Scope chain walk per `OP_LOAD`.
- `scope.New()` — allocates `map[string]any` per `OP_PUSH_SCOPE` (every loop body).
- Context check — `select { case <-ctx.Done() }` on every opcode.
- `fromConst` — type switch on every `OP_PUSH_CONST`.

---

## File Structure

| File | Responsibility | Action |
|------|---------------|--------|
| `internal/vm/value.go` | VM value types and constructors | Modify: add `TypeLoopVar`, `LoopVarVal` constructor |
| `internal/vm/vm.go` | VM execution loop | Modify: loop var, context batching, output pre-sizing, scope changes |
| `internal/scope/scope.go` | Variable scope chain | Modify: change storage from `map[string]any` to `map[string]Value` |
| `internal/compiler/bytecode.go` | Bytecode structures | Modify: add `CompiledConsts []Value` field |
| `internal/compiler/compiler.go` | AST→bytecode compiler | Modify: pre-compile constants to `Value` |
| `internal/vm/vm_profile.go` | Profiling instrumentation | Modify: adapt to scope type changes |
| `internal/vm/vm_noprofile.go` | No-op profiling stubs | No change needed |

---

### Task 1: Struct-Based Loop Variable (eliminate makeLoopMap — 35% of memory)

**Files:**
- Modify: `internal/vm/value.go` (add TypeLoopVar and attribute resolution)
- Modify: `internal/vm/vm.go:984-1010` (replace `makeLoopMap` with `makeLoopVal`)
- Test: `pkg/grove/engine_test.go` (existing loop tests cover this)

The current `makeLoopMap()` allocates a `map[string]any` with 6-8 entries on **every** loop iteration, then wraps it through `FromAny`. Replace with a value type that resolves attributes on demand.

- [ ] **Step 1: Add TypeLoopVar to value.go**

In `internal/vm/value.go`, add a new value type and struct:

```go
// Add to the const block after TypeMacro:
TypeLoopVar // oval: *loopVarData

// Add after the Resolvable interface:

// loopVarData holds loop metadata without map allocation.
type loopVarData struct {
	index  int
	length int
	depth  int
	parent *loopVarData // nil if depth == 1
}

func LoopVarVal(d *loopVarData) Value { return Value{typ: TypeLoopVar, oval: d} }
```

- [ ] **Step 2: Handle TypeLoopVar in GetAttr**

In `internal/vm/value.go`, add a case to `GetAttr` (before the TypeNil case):

```go
case TypeLoopVar:
    ld := obj.oval.(*loopVarData)
    switch name {
    case "index":
        return IntVal(int64(ld.index + 1)), nil
    case "index0":
        return IntVal(int64(ld.index)), nil
    case "first":
        return BoolVal(ld.index == 0), nil
    case "last":
        return BoolVal(ld.index == ld.length-1), nil
    case "length":
        return IntVal(int64(ld.length)), nil
    case "depth":
        return IntVal(int64(ld.depth)), nil
    case "parent":
        if ld.parent != nil {
            return LoopVarVal(ld.parent), nil
        }
        return Nil, nil
    }
    if strict {
        return Nil, fmt.Errorf("undefined loop attribute %q", name)
    }
    return Nil, nil
```

- [ ] **Step 3: Handle TypeLoopVar in String() and Truthy()**

In `value.go`, add to `String()`:
```go
case TypeLoopVar:
    return "[loop]"
```

In `Truthy()` (if it exists), add:
```go
case TypeLoopVar:
    return true
```

- [ ] **Step 4: Pool loopVarData to eliminate allocation**

In `internal/vm/vm.go`, add a pool of loopVarData objects embedded in the VM struct. Since max loop depth is 32, pre-allocate:

```go
// In the VM struct, add:
loopVars [32]loopVarData
```

- [ ] **Step 5: Replace makeLoopMap with makeLoopVal**

In `internal/vm/vm.go`, replace `makeLoopMap` (lines 984-1010):

```go
// makeLoopVal constructs the `loop` magic variable without allocation.
func (v *VM) makeLoopVal() Value {
	ls := &v.loops[v.ldepth-1]
	ld := &v.loopVars[v.ldepth-1]
	ld.index = ls.idx
	ld.length = len(ls.items)
	ld.depth = v.ldepth
	if v.ldepth > 1 {
		ld.parent = &v.loopVars[v.ldepth-2]
	} else {
		ld.parent = nil
	}
	return LoopVarVal(ld)
}
```

- [ ] **Step 6: Update OP_FOR_BIND_1 and OP_FOR_BIND_KV to use makeLoopVal**

In `internal/vm/vm.go`, change both call sites:

```go
// OP_FOR_BIND_1 (line ~407):
v.sc.Set("loop", v.makeLoopVal())

// OP_FOR_BIND_KV (line ~420):
v.sc.Set("loop", v.makeLoopVal())
```

- [ ] **Step 7: Run tests to verify**

Run: `cd /home/theo/Work/grove && go clean -testcache && go test ./... -v`
Expected: All tests pass. Loop template features (loop.index, loop.first, loop.last, loop.parent, etc.) work identically.

- [ ] **Step 8: Run benchmarks to measure improvement**

Run: `cd /home/theo/Work/grove/benchmarks && go test -bench=BenchmarkRender_Loop -benchmem -count=3`
Expected: Significant reduction in B/op and allocs/op for Grove loop benchmarks.

- [ ] **Step 9: Commit**

```bash
git add internal/vm/value.go internal/vm/vm.go
git commit -m "$(cat <<'EOF'
perf: replace loop map allocation with struct-based loop variable

makeLoopMap allocated a map[string]any on every loop iteration (35% of
total memory). Replace with a pre-allocated loopVarData struct on the VM
that resolves attributes on demand via GetAttr.
EOF
)"
```

---

### Task 2: Store Value in Scope (eliminate FromAny on every OP_LOAD)

**Files:**
- Modify: `internal/scope/scope.go` (change `map[string]any` → `map[string]Value`)
- Modify: `internal/vm/vm.go:130-167` (adapt Execute to store Values in scope)
- Modify: `internal/vm/vm.go:211-221` (simplify OP_LOAD — no more FromAny)
- Modify: `internal/vm/vm.go:378-380` (OP_STORE_VAR stores Value directly)
- Test: `internal/scope/scope_test.go` (if exists), `pkg/grove/engine_test.go`

Currently, `scope.Scope` stores `map[string]any`. Every `OP_LOAD` calls `scope.Get` which returns `any`, then wraps it with `FromAny` (type switch). By storing `Value` directly, we eliminate this conversion on every variable access.

**Important:** This is a cross-cutting change. The scope package is imported by both `vm` and potentially other packages. We need to introduce the `Value` type dependency carefully. Since `scope` is internal and only used by `vm`, we can make scope generic or have it accept an interface. The simplest approach: move scope storage to use `any` but have vm always store `Value` in it, then type-assert to `Value` (single assertion, not a type switch). Even simpler: since scope is tiny and internal, just make it store `vm.Value`.

**Circular dependency concern:** `scope` imports nothing, `vm` imports `scope`. If scope imports `vm.Value`, we get a cycle. Solutions:
1. Move `Value` type to its own package (e.g., `internal/value/`)
2. Make scope generic with Go generics
3. Keep scope storing `any` but always store `Value` in it — single type assertion

Option 3 is simplest and still eliminates the full `FromAny` type switch:

- [ ] **Step 1: Write a test verifying current OP_LOAD behavior**

Existing tests already cover this. Run them to establish baseline:

Run: `cd /home/theo/Work/grove && go test ./pkg/grove/ -v -run TestVariable`
Expected: PASS

- [ ] **Step 2: Modify Execute to store Values in scope**

In `internal/vm/vm.go`, change the scope setup in `Execute` (lines 159-167):

```go
globalSc := scope.New(nil)
for k, val := range eng.GlobalData() {
    globalSc.Set(k, FromAny(val))
}
renderSc := scope.New(globalSc)
for k, val := range data {
    renderSc.Set(k, FromAny(val))
}
v.sc = scope.New(renderSc)
```

This converts all data to `Value` once at scope creation time, not on every access.

- [ ] **Step 3: Simplify OP_LOAD to skip FromAny**

In `internal/vm/vm.go`, change OP_LOAD (lines 211-221):

```go
case compiler.OP_LOAD:
    name := bc.Names[instr.A]
    val, found := v.sc.Get(name)
    if !found {
        if v.eng.StrictVariables() {
            return "", &runtimeErr{msg: fmt.Sprintf("undefined variable %q", name)}
        }
        v.push(Nil)
    } else {
        v.push(val.(Value))
    }
```

- [ ] **Step 4: Ensure OP_STORE_VAR stores Value directly**

In `internal/vm/vm.go`, OP_STORE_VAR (line 380) already does `v.sc.Set(bc.Names[instr.A], val)` where `val` is a `Value`. Since scope stores `any`, this already works — the `Value` is stored as-is.

- [ ] **Step 5: Update all other scope.Set call sites to store Values**

Search for all `v.sc.Set(` calls in vm.go. Each must store a `Value`:

- `OP_FOR_BIND_1`: `v.sc.Set(varName, ls.items[ls.idx])` — already a `Value` ✓
- `OP_FOR_BIND_KV`: `v.sc.Set(name1, StringVal(...))` / `v.sc.Set(name2, ls.items[ls.idx])` — already Values ✓
- `OP_FOR_BIND_KV`: `v.sc.Set(name1, IntVal(...))` — already a Value ✓
- `OP_CAPTURE_END`: `v.sc.Set(varName, StringVal(content))` — already a Value ✓
- `makeLoopVal`: `v.sc.Set("loop", ...)` — already a Value ✓
- Any macro/component scope setup — verify these also store Values

- [ ] **Step 6: Update GetAttr for TypeMap to handle Value-stored maps**

When maps are stored in scope as `Value` (via `FromAny`), `GetAttr` on TypeMap calls `FromAny(v)` on the map value. This still works but is a remaining `FromAny` call. For map access specifically, the values inside the map are still `any` (since `MapVal` wraps `map[string]any`). This is fine — map attribute access is less frequent than plain variable loads. No change needed here.

- [ ] **Step 7: Run tests**

Run: `cd /home/theo/Work/grove && go clean -testcache && go test ./... -v`
Expected: All tests pass.

- [ ] **Step 8: Run benchmarks**

Run: `cd /home/theo/Work/grove/benchmarks && go test -bench=BenchmarkRender -benchmem -count=3`
Expected: Reduction in allocs/op across all scenarios, especially Simple and Conditional where variable loads dominate.

- [ ] **Step 9: Commit**

```bash
git add internal/vm/vm.go
git commit -m "$(cat <<'EOF'
perf: convert data to Value at scope creation, not on every load

Scope now always stores Value objects. FromAny conversion happens once
when data enters the scope, eliminating the type switch on every OP_LOAD.
EOF
)"
```

---

### Task 3: Pre-Compile Constants to Value (eliminate fromConst per OP_PUSH_CONST)

**Files:**
- Modify: `internal/compiler/bytecode.go:142-147` (add `ValueConsts` field)
- Modify: `internal/vm/vm.go:208-209` (use pre-compiled values)
- Modify: `internal/vm/vm.go` (add `PrecompileConstants` function)
- Test: `pkg/grove/engine_test.go` (existing tests)

Every `OP_PUSH_CONST` calls `fromConst(bc.Consts[instr.A])` which is a 4-way type switch. Pre-compile the constant pool to `[]Value` once at compile time.

**Circular dependency concern:** `compiler` cannot import `vm.Value`. Solutions:
1. Add a `ValueConsts []any` field to Bytecode and have the VM populate it on first use
2. Store the pre-compiled values on the VM side

Option 2 avoids touching the compiler package. The VM can pre-compile constants once per template execution.

- [ ] **Step 1: Add constant pre-compilation in the VM**

In `internal/vm/vm.go`, add a helper and a cache field. Since bytecode is immutable and shared, we can pre-compile constants into a `[]Value` slice stored alongside the bytecode. But Bytecode is in the compiler package. Instead, use a sync.Map cache in the vm package:

```go
// At package level in vm.go:
var constCache sync.Map // map[*compiler.Bytecode][]Value
```

- [ ] **Step 2: Add precompileConsts function**

In `internal/vm/vm.go`:

```go
func precompileConsts(bc *compiler.Bytecode) []Value {
	if cached, ok := constCache.Load(bc); ok {
		return cached.([]Value)
	}
	vals := make([]Value, len(bc.Consts))
	for i, c := range bc.Consts {
		vals[i] = fromConst(c)
	}
	constCache.Store(bc, vals)
	return vals
}
```

- [ ] **Step 3: Use pre-compiled constants in the run loop**

In `internal/vm/vm.go`, at the start of `run()`, add:

```go
func (v *VM) run(ctx context.Context, bc *compiler.Bytecode) (string, error) {
	ip := 0
	instrs := bc.Instrs
	consts := precompileConsts(bc)
	ps := profileInit()
```

Then change OP_PUSH_CONST:

```go
case compiler.OP_PUSH_CONST:
    v.push(consts[instr.A])
```

- [ ] **Step 4: Run tests**

Run: `cd /home/theo/Work/grove && go clean -testcache && go test ./... -v`
Expected: All tests pass.

- [ ] **Step 5: Commit**

```bash
git add internal/vm/vm.go
git commit -m "$(cat <<'EOF'
perf: pre-compile constant pool to Value slice

Cache the fromConst type switch result per bytecode object, eliminating
the per-opcode type switch on every OP_PUSH_CONST execution.
EOF
)"
```

---

### Task 4: Batch Context Cancellation Checks

**Files:**
- Modify: `internal/vm/vm.go:189-194` (replace per-opcode check with batched check)

The `select { case <-ctx.Done(): ... default: }` on every opcode adds overhead even when the channel is never closed. Check every 64 opcodes instead.

- [ ] **Step 1: Add an opcode counter and batch the check**

In `internal/vm/vm.go`, modify the `run()` loop:

```go
func (v *VM) run(ctx context.Context, bc *compiler.Bytecode) (string, error) {
	ip := 0
	instrs := bc.Instrs
	consts := precompileConsts(bc)
	ps := profileInit()
	done := ctx.Done()
	opcCount := 0
	for ip < len(instrs) {
		opcCount++
		if opcCount&63 == 0 { // check every 64 opcodes
			select {
			case <-done:
				return "", ctx.Err()
			default:
			}
		}

		instr := instrs[ip]
```

Note: We cache `ctx.Done()` as `done` to avoid the method call overhead per check.

- [ ] **Step 2: Run tests**

Run: `cd /home/theo/Work/grove && go clean -testcache && go test ./... -v`
Expected: All tests pass. Context cancellation still works (just up to 64 opcodes delayed).

- [ ] **Step 3: Run benchmarks**

Run: `cd /home/theo/Work/grove/benchmarks && go test -bench=BenchmarkRender_Simple -benchmem -count=3`
Expected: Measurable improvement on simple/short templates where the overhead was proportionally larger.

- [ ] **Step 4: Commit**

```bash
git add internal/vm/vm.go
git commit -m "$(cat <<'EOF'
perf: batch context cancellation check every 64 opcodes

Reduces per-opcode overhead of the select on ctx.Done() channel.
Cancellation is still detected within 64 instructions.
EOF
)"
```

---

### Task 5: Pre-Size Output Buffer

**Files:**
- Modify: `internal/vm/vm.go:130-183` (grow output buffer based on template size)
- Modify: `internal/compiler/bytecode.go` (add `EstimatedOutputSize` field)
- Modify: `internal/compiler/compiler.go` (calculate estimated output size)

`strings.Builder.WriteString` accounts for 51% of memory. The builder starts empty and grows by doubling, causing multiple allocations. Pre-sizing based on the template's static content eliminates most regrowth.

- [ ] **Step 1: Add EstimatedOutputSize to Bytecode**

In `internal/compiler/bytecode.go`, add to the `Bytecode` struct:

```go
EstimatedOutputSize int // sum of static string constant lengths (hint for output buffer)
```

- [ ] **Step 2: Calculate EstimatedOutputSize during compilation**

In `internal/compiler/compiler.go`, after compilation is complete, calculate the estimate. Find the function that produces the final `Bytecode` and add:

```go
// After building the bytecode, estimate output size from string constants
est := 0
for _, c := range bc.Consts {
    if s, ok := c.(string); ok {
        est += len(s)
    }
}
bc.EstimatedOutputSize = est
```

- [ ] **Step 3: Pre-grow the output buffer in Execute**

In `internal/vm/vm.go`, in the `Execute` function, before calling `v.run()`:

```go
if bc.EstimatedOutputSize > 0 {
    v.out.Grow(bc.EstimatedOutputSize)
}
```

- [ ] **Step 4: Run tests**

Run: `cd /home/theo/Work/grove && go clean -testcache && go test ./... -v`
Expected: All tests pass.

- [ ] **Step 5: Commit**

```bash
git add internal/compiler/bytecode.go internal/compiler/compiler.go internal/vm/vm.go
git commit -m "$(cat <<'EOF'
perf: pre-size output buffer from template static content

Calculate estimated output size during compilation and use Builder.Grow
to avoid repeated buffer doubling during rendering.
EOF
)"
```

---

### Task 6: Scope Allocation Reduction

**Files:**
- Modify: `internal/scope/scope.go` (add pre-sized constructor, Reset method)
- Modify: `internal/vm/vm.go` (reuse scope objects for loop bodies)

Every `OP_PUSH_SCOPE` calls `scope.New(parent)` which allocates `make(map[string]any)`. In loop bodies, this happens every iteration. We can pre-size the map (loop bodies typically set 2-3 variables: loop var + "loop" meta) and reuse scope objects.

- [ ] **Step 1: Add NewWithSize and Reset to scope**

In `internal/scope/scope.go`:

```go
// NewWithSize creates a scope with a pre-sized variable map.
func NewWithSize(parent *Scope, size int) *Scope {
	return &Scope{vars: make(map[string]any, size), parent: parent}
}

// Reset clears all variables and sets a new parent, reusing the map memory.
func (s *Scope) Reset(parent *Scope) {
	for k := range s.vars {
		delete(s.vars, k)
	}
	s.parent = parent
}
```

- [ ] **Step 2: Add reusable loop scopes to VM**

In `internal/vm/vm.go`, add to the VM struct:

```go
loopScopes [32]*scope.Scope // pre-allocated scopes for loop bodies
```

- [ ] **Step 3: Initialize loop scopes in Execute**

In `internal/vm/vm.go`, in the `Execute` function, initialize loop scopes lazily. Change `OP_PUSH_SCOPE` to reuse when inside a loop:

```go
case compiler.OP_PUSH_SCOPE:
    if v.ldepth > 0 && v.loopScopes[v.ldepth-1] != nil {
        // Reuse existing loop scope
        v.loopScopes[v.ldepth-1].Reset(v.sc)
        v.sc = v.loopScopes[v.ldepth-1]
    } else if v.ldepth > 0 {
        // First iteration: create scope, cache for reuse
        s := scope.NewWithSize(v.sc, 3) // loop var + "loop" + maybe key
        v.loopScopes[v.ldepth-1] = s
        v.sc = s
    } else {
        v.sc = scope.New(v.sc)
    }
```

- [ ] **Step 4: Clear loop scopes on VM reset**

In the `Execute` defer block, add:

```go
for i := range v.loopScopes {
    v.loopScopes[i] = nil
}
```

- [ ] **Step 5: Run tests**

Run: `cd /home/theo/Work/grove && go clean -testcache && go test ./... -v`
Expected: All tests pass. Variable scoping in loops still works correctly (scope push/pop semantics preserved).

- [ ] **Step 6: Run benchmarks**

Run: `cd /home/theo/Work/grove/benchmarks && go test -bench=BenchmarkRender -benchmem -count=3`
Expected: Reduction in allocs/op for Loop and Complex scenarios.

- [ ] **Step 7: Commit**

```bash
git add internal/scope/scope.go internal/vm/vm.go
git commit -m "$(cat <<'EOF'
perf: reuse scope objects in loop bodies

Pre-allocate and reuse scope objects for loop iterations instead of
allocating a new map on every OP_PUSH_SCOPE inside a loop.
EOF
)"
```

---

### Task 7: Final Validation and Benchmark Comparison

**Files:**
- No code changes — validation only

- [ ] **Step 1: Run full test suite**

Run: `cd /home/theo/Work/grove && go clean -testcache && go test ./... -v`
Expected: All tests pass.

- [ ] **Step 2: Run full benchmark comparison**

Run: `cd /home/theo/Work/grove/benchmarks && go test -bench=BenchmarkRender -benchmem -count=6 -timeout=10m`
Expected: Grove beats Pongo2 on all scenarios.

- [ ] **Step 3: Run timing benchmarks**

Run: `cd /home/theo/Work/grove/benchmarks && bash run-timing.sh -n 1000`
Expected: Grove faster than Pongo2 on Large Loop, Nested Loops, and Complex Page.

- [ ] **Step 4: Run profiler to verify bottleneck reduction**

Run: `cd /home/theo/Work/grove/benchmarks && bash run-profile.sh -s all -n 5000`
Expected: `makeLoopMap` no longer appears in memory profile. Loop category time reduced significantly.

- [ ] **Step 5: Commit any remaining changes**

If any adjustments were needed during validation, commit them.

---

## Expected Impact Summary

| Optimization | Target Bottleneck | Expected Improvement |
|---|---|---|
| Task 1: Struct loop var | makeLoopMap (35% memory) | ~30% reduction in loop allocs |
| Task 2: Value in scope | FromAny per OP_LOAD (4.5% memory + CPU) | ~10% faster variable access |
| Task 3: Pre-compile consts | fromConst per OP_PUSH_CONST | ~5% faster const access |
| Task 4: Batch ctx check | select per opcode | ~5-10% faster short templates |
| Task 5: Pre-size buffer | Builder regrowth (51% memory) | Fewer large allocations |
| Task 6: Scope reuse | scope.New per loop iter | ~15% reduction in loop allocs |

Combined, these should bring Grove's loop and complex benchmarks below Pongo2's numbers while maintaining all existing functionality.
