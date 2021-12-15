package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
)

func Test_RS(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	// 256
	token, err := test_Generate(NewRSGenerator(RS256, key))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(token)
	test_Verify(t, token, NewRSVerifier(RS256, &key.PublicKey))
	// 384
	token, err = test_Generate(NewRSGenerator(RS384, key))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(token)
	test_Verify(t, token, NewRSVerifier(RS384, &key.PublicKey))
	// 512
	token, err = test_Generate(NewRSGenerator(RS512, key))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(token)
	test_Verify(t, token, NewRSVerifier(RS512, &key.PublicKey))
}

func Benchmark_RS256_Generate(b *testing.B) {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	benchmark_Generate(b, NewRSGenerator(RS256, key))
}

func Benchmark_RS384_Generate(b *testing.B) {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	benchmark_Generate(b, NewRSGenerator(RS384, key))
}

func Benchmark_RS512_Generate(b *testing.B) {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	benchmark_Generate(b, NewRSGenerator(RS512, key))
}

func Benchmark_RS256_Verify(b *testing.B) {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	token, err := test_Generate(NewRSGenerator(RS256, key))
	if err != nil {
		b.Fatal(err)
	}
	benchmark_Verify(b, token, NewRSVerifier(RS256, &key.PublicKey))
}

func Benchmark_RS384_Verify(b *testing.B) {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	token, err := test_Generate(NewRSGenerator(RS384, key))
	if err != nil {
		b.Fatal(err)
	}
	benchmark_Verify(b, token, NewRSVerifier(RS384, &key.PublicKey))
}

func Benchmark_RS512_Verify(b *testing.B) {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	token, err := test_Generate(NewRSGenerator(RS512, key))
	if err != nil {
		b.Fatal(err)
	}
	benchmark_Verify(b, token, NewRSVerifier(RS512, &key.PublicKey))
}
