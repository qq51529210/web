package jwt

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"hash"
	"math/big"
	"math/bits"
	"sync"
)

// Use for ES algorithm.
const _s = bits.UintSize / 8

// Use for ES algorithm.
type bigint struct {
	b    []byte
	r, s big.Int
}

// Copy from:
// go/src/math/big/nat.go func (z nat) bytes(buf []byte) (i int)
// go/src/math/big/int.go func (x *Int) Bytes() []byte
func (b *bigint) Encode(r, s *big.Int) []byte {
	sz, rz := s.Bits(), r.Bits()
	n := (len(sz) + len(rz)) * _s
	if cap(b.b) < n {
		b.b = make([]byte, n)
	} else {
		b.b = b.b[:n]
	}
	i := len(b.b)
	for _, d := range sz {
		for j := 0; j < _s; j++ {
			i--
			b.b[i] = byte(d)
			d >>= 8
		}
	}
	for i < len(b.b) && b.b[i] == 0 {
		i++
	}
	for _, d := range rz {
		for j := 0; j < _s; j++ {
			i--
			b.b[i] = byte(d)
			d >>= 8
		}
	}
	for i < len(b.b) && b.b[i] == 0 {
		i++
	}
	return b.b[i:]
}

func (b *bigint) Decode(buf []byte) {
	n := len(buf) / 2
	b.r.SetBytes(buf[:n])
	b.s.SetBytes(buf[n:])
}

var (
	es256GenPool sync.Pool
	es384GenPool sync.Pool
	es512GenPool sync.Pool
)

func init() {
	es256GenPool.New = func() interface{} {
		return NewESGenerator(HS256, nil)
	}
	es384GenPool.New = func() interface{} {
		return NewESGenerator(HS384, nil)
	}
	es512GenPool.New = func() interface{} {
		return NewESGenerator(HS512, nil)
	}
}

func NewESGenerator(alg Alg, key *ecdsa.PrivateKey) *ESGenerator {
	s := new(ESGenerator)
	s.Init(alg, key)
	return s
}

func NewESVerifier(alg Alg, key *ecdsa.PublicKey) *ESVerifier {
	s := new(ESVerifier)
	s.Init(alg, key)
	return s
}

func GenerateES256(header, payload map[string]interface{}, key *ecdsa.PrivateKey) (string, error) {
	return generateES(header, payload, key, &es256GenPool)
}

func GenerateES384(header, payload map[string]interface{}, key *ecdsa.PrivateKey) (string, error) {
	return generateES(header, payload, key, &es384GenPool)
}

func GenerateES512(header, payload map[string]interface{}, key *ecdsa.PrivateKey) (string, error) {
	return generateES(header, payload, key, &es512GenPool)
}

func generateES(header, payload map[string]interface{}, key *ecdsa.PrivateKey, pool *sync.Pool) (string, error) {
	g := pool.Get().(*ESGenerator)
	g.key = key
	token, err := g.Generate(header, payload)
	pool.Put(g)
	return token, err
}

type ESGenerator struct {
	alg    Alg
	enc    encoder
	key    *ecdsa.PrivateKey
	bi     bigint
	hash   hash.Hash
	crypto crypto.Hash
	sign   []byte
}

func (g *ESGenerator) Init(alg Alg, key *ecdsa.PrivateKey) {
	g.alg = alg
	g.key = key
	g.crypto = alg.CryptoHash()
	g.hash = g.crypto.New()
	g.sign = make([]byte, g.hash.Size())
	g.enc.Init()
}

func (g *ESGenerator) Generate(header, payload map[string]interface{}) (string, error) {
	header[ALG] = g.alg
	// Encode
	err := g.enc.Enc(header, payload)
	if err != nil {
		return "", err
	}
	// Generate hash signature.
	g.hash.Reset()
	g.hash.Write(g.enc.token)
	g.hash.Sum(g.sign[:0])
	r, s, err := ecdsa.Sign(rand.Reader, g.key, g.sign)
	if err != nil {
		return "", err
	}
	// '.' between payload and signature.
	g.enc.token = append(g.enc.token, '.')
	// Base64 hash signature.
	g.enc.Base64(g.bi.Encode(r, s))
	return string(g.enc.token), nil
}

type ESVerifier struct {
	alg    Alg
	key    *ecdsa.PublicKey
	bi     bigint
	hash   hash.Hash
	crypto crypto.Hash
	sign   []byte
}

func (v *ESVerifier) Init(alg Alg, key *ecdsa.PublicKey) {
	v.alg = alg
	v.key = key
	v.crypto = alg.CryptoHash()
	v.hash = v.crypto.New()
	v.sign = make([]byte, v.hash.Size())
}

func (v *ESVerifier) Verify(data string, signature []byte) error {
	b := string2bytes(data)
	// Hash data.
	v.hash.Reset()
	v.hash.Write(b)
	v.hash.Sum(v.sign[:0])
	v.bi.Decode(signature)
	// ECDSA verify.
	if !ecdsa.Verify(v.key, v.sign, &v.bi.r, &v.bi.s) {
		return errInvalidJWT
	}
	return nil
}
