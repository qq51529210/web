package jwt

import (
	"testing"
)

func Test_HS(t *testing.T) {
	// 256
	key := []byte("HS256")
	token, err := test_Generate(NewHSGenerator(HS256, key))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(token)
	test_Verify(t, token, NewHSVerifier(HS256, key))
	// 384
	key = []byte("HS384")
	token, err = test_Generate(NewHSGenerator(HS384, key))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(token)
	test_Verify(t, token, NewHSVerifier(HS384, key))
	// 512
	key = []byte("HS512")
	token, err = test_Generate(NewHSGenerator(HS512, key))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(token)
	test_Verify(t, token, NewHSVerifier(HS512, key))
}

func Benchmark_HS256_Generate(b *testing.B) {
	benchmark_Generate(b, NewHSGenerator(HS256, []byte("HS256")))
}

func Benchmark_HS384_Generate(b *testing.B) {
	benchmark_Generate(b, NewHSGenerator(HS384, []byte("HS384")))
}

func Benchmark_HS512_Generate(b *testing.B) {
	benchmark_Generate(b, NewHSGenerator(HS512, []byte("HS512")))
}

func Benchmark_HS256_Verify(b *testing.B) {
	key := []byte("HS256")
	token, err := test_Generate(NewHSGenerator(HS256, key))
	if err != nil {
		b.Fatal(err)
	}
	benchmark_Verify(b, token, NewHSVerifier(HS256, key))
}

func Benchmark_HS384_Verify(b *testing.B) {
	key := []byte("HS384")
	token, err := test_Generate(NewHSGenerator(HS384, key))
	if err != nil {
		b.Fatal(err)
	}
	benchmark_Verify(b, token, NewHSVerifier(HS384, key))
}

func Benchmark_HS512_Verify(b *testing.B) {
	key := []byte("HS512")
	token, err := test_Generate(NewHSGenerator(HS512, key))
	if err != nil {
		b.Fatal(err)
	}
	benchmark_Verify(b, token, NewHSVerifier(HS512, key))
}
