package jwt

import "testing"

func test_Generate(gen Generator) (string, error) {
	header := make(map[string]interface{})
	payload := make(map[string]interface{})
	header["1"] = 1
	payload["2"] = "2"
	return gen.Generate(header, payload)
}

func test_Verify(t *testing.T, token string, v Verifier) {
	header, payload, err := Verify(token, func(a Alg) Verifier { return v })
	if err != nil {
		t.Fatal(err)
	}
	val, ok := header["1"]
	if !ok {
		t.FailNow()
	}
	if n, ok := val.(float64); !ok || n != 1 {
		t.FailNow()
	}
	val, ok = payload["2"]
	if !ok {
		t.FailNow()
	}
	if n, ok := val.(string); !ok || n != "2" {
		t.FailNow()
	}
}

func benchmark_Generate(b *testing.B, g Generator) {
	b.ReportAllocs()
	b.ResetTimer()
	header := make(map[string]interface{})
	payload := make(map[string]interface{})
	header["1"] = 1
	payload["2"] = "2"
	for i := 0; i < b.N; i++ {
		g.Generate(header, payload)
	}
}

func benchmark_Verify(b *testing.B, t string, v Verifier) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Verify(t, func(a Alg) Verifier {
			return v
		})
	}
}
