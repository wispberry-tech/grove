package benchmarks

import (
	"testing"
)

// ---------- Parse-only benchmarks ----------

func BenchmarkParse_Simple(b *testing.B) {
	benchParse(b, SimpleTemplates)
}

func BenchmarkParse_Loop(b *testing.B) {
	benchParse(b, LoopTemplates)
}

func BenchmarkParse_Conditional(b *testing.B) {
	benchParse(b, ConditionalTemplates)
}

func BenchmarkParse_Complex(b *testing.B) {
	benchParse(b, ComplexTemplates)
}

// ---------- Render-only benchmarks (pre-parsed) ----------

func BenchmarkRender_Simple(b *testing.B) {
	benchRender(b, SimpleTemplates, WrapSimple())
}

func BenchmarkRender_Loop(b *testing.B) {
	benchRender(b, LoopTemplates, WrapLoop())
}

func BenchmarkRender_Conditional(b *testing.B) {
	benchRender(b, ConditionalTemplates, WrapConditional())
}

func BenchmarkRender_Complex(b *testing.B) {
	benchRender(b, ComplexTemplates, WrapComplex())
}

// ---------- Full (parse + render) benchmarks ----------

func BenchmarkFull_Simple(b *testing.B) {
	benchFull(b, SimpleTemplates, WrapSimple())
}

func BenchmarkFull_Loop(b *testing.B) {
	benchFull(b, LoopTemplates, WrapLoop())
}

func BenchmarkFull_Conditional(b *testing.B) {
	benchFull(b, ConditionalTemplates, WrapConditional())
}

func BenchmarkFull_Complex(b *testing.B) {
	benchFull(b, ComplexTemplates, WrapComplex())
}

// ---------- Helpers ----------

func benchParse(b *testing.B, templates map[string]string) {
	for _, eng := range AllEngines() {
		src := templates[eng.Name()]
		b.Run(eng.Name(), func(b *testing.B) {
			b.ReportAllocs()
			// For Grove, use ForceParse to bypass the LRU cache and measure
			// actual lex+parse+compile cost each iteration.
			if ge, ok := eng.(*groveEngine); ok {
				for b.Loop() {
					if err := ge.ForceParse(src); err != nil {
						b.Fatal(err)
					}
				}
			} else {
				for b.Loop() {
					if err := eng.Parse("bench", src); err != nil {
						b.Fatal(err)
					}
				}
			}
		})
	}
}

func benchRender(b *testing.B, templates map[string]string, data map[string]any) {
	for _, eng := range AllEngines() {
		src := templates[eng.Name()]
		if err := eng.Parse("bench", src); err != nil {
			b.Fatalf("%s parse: %v", eng.Name(), err)
		}
		d := EngineData(eng, data)
		b.Run(eng.Name(), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				if _, err := eng.Render("bench", d); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func benchFull(b *testing.B, templates map[string]string, data map[string]any) {
	for _, eng := range AllEngines() {
		src := templates[eng.Name()]
		d := EngineData(eng, data)
		b.Run(eng.Name(), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				if _, err := eng.ParseAndRender("bench", src, d); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
