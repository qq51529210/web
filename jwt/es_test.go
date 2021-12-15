package jwt

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"
)

func Test_ES(t *testing.T) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	// 256
	token, err := test_Generate(NewESGenerator(ES256, key))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(token)
	test_Verify(t, token, NewESVerifier(ES256, &key.PublicKey))
	// 384
	token, err = test_Generate(NewESGenerator(ES384, key))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(token)
	test_Verify(t, token, NewESVerifier(ES384, &key.PublicKey))
	// 512
	token, err = test_Generate(NewESGenerator(ES512, key))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(token)
	test_Verify(t, token, NewESVerifier(ES512, &key.PublicKey))
}

func Benchmark_ES256_Generate(b *testing.B) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	benchmark_Generate(b, NewESGenerator(ES256, key))
}

func Benchmark_ES384_Generate(b *testing.B) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	benchmark_Generate(b, NewESGenerator(ES384, key))
}

func Benchmark_ES512_Generate(b *testing.B) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	benchmark_Generate(b, NewESGenerator(ES512, key))
}

func Benchmark_ES256_Verify(b *testing.B) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	token, err := test_Generate(NewESGenerator(ES256, key))
	if err != nil {
		b.Fatal(err)
	}
	benchmark_Verify(b, token, NewESVerifier(ES256, &key.PublicKey))
}

func Benchmark_ES384_Verify(b *testing.B) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	token, err := test_Generate(NewESGenerator(ES384, key))
	if err != nil {
		b.Fatal(err)
	}
	benchmark_Verify(b, token, NewESVerifier(ES384, &key.PublicKey))
}

func Benchmark_ES512_Verify(b *testing.B) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	token, err := test_Generate(NewESGenerator(ES512, key))
	if err != nil {
		b.Fatal(err)
	}
	benchmark_Verify(b, token, NewESVerifier(ES512, &key.PublicKey))
}
