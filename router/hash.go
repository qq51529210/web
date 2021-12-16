package router

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
	"sync"
)

var (
	md5Pool    sync.Pool
	sha1Pool   sync.Pool
	sha256Pool sync.Pool
	sha512Pool sync.Pool
)

func init() {

	md5Pool.New = func() interface{} {
		return &hashBuffer{
			hash: md5.New(),
			buf:  make([]byte, 0, md5.Size*2),
			sum:  make([]byte, md5.Size),
		}
	}
	sha1Pool.New = func() interface{} {
		return &hashBuffer{
			hash: sha1.New(),
			buf:  make([]byte, 0, sha1.Size*2),
			sum:  make([]byte, sha1.Size),
		}
	}
	sha256Pool.New = func() interface{} {
		return &hashBuffer{
			hash: sha256.New(),
			buf:  make([]byte, 0, sha256.Size*2),
			sum:  make([]byte, sha256.Size),
		}
	}
	sha512Pool.New = func() interface{} {
		return &hashBuffer{
			hash: sha512.New(),
			buf:  make([]byte, 0, sha512.Size*2),
			sum:  make([]byte, sha512.Size),
		}
	}
}

type hashBuffer struct {
	hash hash.Hash
	buf  []byte
	sum  []byte
}

func (h *hashBuffer) HashString(s string) string {
	h.buf = h.buf[:0]
	h.buf = append(h.buf, s...)
	return h.hashAndHex()
}

func (h *hashBuffer) Hash(b []byte) string {
	h.buf = h.buf[:0]
	h.buf = append(h.buf, b...)
	return h.hashAndHex()
}

func (h *hashBuffer) hashAndHex() string {
	h.hash.Reset()
	h.hash.Write(h.buf)
	h.hash.Sum(h.sum[:0])
	h.buf = h.buf[:h.hash.Size()*2]
	hex.Encode(h.buf, h.sum)
	return string(h.buf)
}

func hashStringWithPool(p *sync.Pool, s string) string {
	h := p.Get().(*hashBuffer)
	s = h.HashString(s)
	p.Put(h)
	return s
}

func hashWithPool(p *sync.Pool, b []byte) string {
	h := p.Get().(*hashBuffer)
	s := h.Hash(b)
	p.Put(h)
	return s
}

// Return hex string of s MD5.
func MD5(b []byte) string {
	return hashWithPool(&md5Pool, b)
}

// Return hex string of s SHA1.
func SHA1(b []byte) string {
	return hashWithPool(&sha1Pool, b)
}

// Return hex string of s SHA256.
func SHA256(b []byte) string {
	return hashWithPool(&sha256Pool, b)
}

// Return hex string of s SHA512.
func SHA512(b []byte) string {
	return hashWithPool(&sha512Pool, b)
}

// Return hex string of s MD5.
func MD5String(s string) string {
	return hashStringWithPool(&md5Pool, s)
}

// Return hex string of s SHA1.
func SHA1String(s string) string {
	return hashStringWithPool(&sha1Pool, s)
}

// Return hex string of s SHA256.
func SHA256String(s string) string {
	return hashStringWithPool(&sha256Pool, s)
}

// Return hex string of s SHA512.
func SHA512String(s string) string {
	return hashStringWithPool(&sha512Pool, s)
}
