package router

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"testing"
)

func Benchmark_Hash_MD5(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		MD5String("hash md5")
	}
}

func Benchmark_STD_MD5(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	d := []byte("hash md5")
	hd := make([]byte, md5.Size*2)
	for i := 0; i < b.N; i++ {
		h := md5.New()
		h.Write(d)
		h.Sum(nil)
		hex.Encode(hd, d)
	}
}

func Benchmark_Hash_SHA1(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		MD5String("hash sha1")
	}
}

func Benchmark_STD_SHA1(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	d := []byte("hash sha1")
	hd := make([]byte, sha1.Size*2)
	for i := 0; i < b.N; i++ {
		h := sha1.New()
		h.Write(d)
		h.Sum(nil)
		hex.Encode(hd, d)
	}
}

func Benchmark_Hash_SHA256(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		MD5String("hash sha256")
	}
}

func Benchmark_STD_SHA256(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	d := []byte("hash sha256")
	hd := make([]byte, sha256.Size*2)
	for i := 0; i < b.N; i++ {
		h := sha256.New()
		h.Write(d)
		h.Sum(nil)
		hex.Encode(hd, d)
	}
}

func Benchmark_SHA512(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		MD5String("hash sha512")
	}
}

func Benchmark_STD_SHA512(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	d := []byte("hash sha512")
	hd := make([]byte, sha512.Size*2)
	for i := 0; i < b.N; i++ {
		h := sha512.New()
		h.Write(d)
		h.Sum(nil)
		hex.Encode(hd, d)
	}
}
