package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"testing"
)

func Test_PS(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	var opt rsa.PSSOptions
	// 256
	token, err := test_Generate(NewPSGenerator(PS256, key, &opt))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(token)
	test_Verify(t, token, NewPSVerifier(PS256, &key.PublicKey, &opt))
	// 384
	token, err = test_Generate(NewPSGenerator(PS384, key, &opt))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(token)
	test_Verify(t, token, NewPSVerifier(PS384, &key.PublicKey, &opt))
	// 512
	token, err = test_Generate(NewPSGenerator(PS512, key, &opt))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(token)
	test_Verify(t, token, NewPSVerifier(PS512, &key.PublicKey, &opt))
}

func Benchmark_PS256_Generate(b *testing.B) {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	benchmark_Generate(b, NewPSGenerator(PS256, key, &rsa.PSSOptions{}))
}

func Benchmark_PS384_Generate(b *testing.B) {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	benchmark_Generate(b, NewPSGenerator(PS384, key, &rsa.PSSOptions{}))
}

func Benchmark_PS512_Generate(b *testing.B) {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	benchmark_Generate(b, NewPSGenerator(PS512, key, &rsa.PSSOptions{}))
}

func Benchmark_PS256_Verify(b *testing.B) {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	var opt rsa.PSSOptions
	token, err := test_Generate(NewPSGenerator(PS256, key, &opt))
	if err != nil {
		b.Fatal(err)
	}
	benchmark_Verify(b, token, NewPSVerifier(PS256, &key.PublicKey, &opt))
}

func Benchmark_PS384_Verify(b *testing.B) {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	var opt rsa.PSSOptions
	token, err := test_Generate(NewPSGenerator(PS384, key, &opt))
	if err != nil {
		b.Fatal(err)
	}
	benchmark_Verify(b, token, NewPSVerifier(PS384, &key.PublicKey, &opt))
}

func Benchmark_PS512_Verify(b *testing.B) {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	var opt rsa.PSSOptions
	token, err := test_Generate(NewPSGenerator(PS512, key, &opt))
	if err != nil {
		b.Fatal(err)
	}
	benchmark_Verify(b, token, NewPSVerifier(PS512, &key.PublicKey, &opt))
}
