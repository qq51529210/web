package jwt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"hash"
	"sync"
)

var (
	rs256GenPool sync.Pool
	rs384GenPool sync.Pool
	rs512GenPool sync.Pool
)

func init() {
	rs256GenPool.New = func() interface{} {
		return NewRSGenerator(HS256, nil)
	}
	rs384GenPool.New = func() interface{} {
		return NewRSGenerator(HS384, nil)
	}
	rs512GenPool.New = func() interface{} {
		return NewRSGenerator(HS512, nil)
	}
}

func NewRSGenerator(alg Alg, key *rsa.PrivateKey) *RSGenerator {
	s := new(RSGenerator)
	s.Init(alg, key)
	return s
}

func NewRSVerifier(alg Alg, key *rsa.PublicKey) *RSVerifier {
	s := new(RSVerifier)
	s.Init(alg, key)
	return s
}

func GenerateRS256(header, payload map[string]interface{}, key *rsa.PrivateKey) (string, error) {
	return generateRS(header, payload, key, &rs256GenPool)
}

func GenerateRS384(header, payload map[string]interface{}, key *rsa.PrivateKey) (string, error) {
	return generateRS(header, payload, key, &rs384GenPool)
}

func GenerateRS512(header, payload map[string]interface{}, key *rsa.PrivateKey) (string, error) {
	return generateRS(header, payload, key, &rs512GenPool)
}

func generateRS(header, payload map[string]interface{}, key *rsa.PrivateKey, pool *sync.Pool) (string, error) {
	g := pool.Get().(*RSGenerator)
	g.key = key
	token, err := g.Generate(header, payload)
	pool.Put(g)
	return token, err
}

type RSGenerator struct {
	alg    Alg
	enc    encoder
	key    *rsa.PrivateKey
	hash   hash.Hash
	crypto crypto.Hash
	sign   []byte
}

func (g *RSGenerator) Init(alg Alg, key *rsa.PrivateKey) {
	g.alg = alg
	g.key = key
	g.crypto = alg.CryptoHash()
	g.hash = g.crypto.New()
	g.sign = make([]byte, g.hash.Size())
	g.enc.Init()
}

func (g *RSGenerator) Generate(header, payload map[string]interface{}) (string, error) {
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
	data, err := rsa.SignPKCS1v15(rand.Reader, g.key, g.crypto, g.sign)
	if err != nil {
		return "", err
	}
	// '.' between payload and signature.
	g.enc.token = append(g.enc.token, '.')
	// Base64 hash signature.
	g.enc.Base64(data)
	return string(g.enc.token), nil
}

type RSVerifier struct {
	alg    Alg
	key    *rsa.PublicKey
	hash   hash.Hash
	crypto crypto.Hash
	sign   []byte
}

func (v *RSVerifier) Init(alg Alg, key *rsa.PublicKey) {
	v.alg = alg
	v.key = key
	v.crypto = alg.CryptoHash()
	v.hash = v.crypto.New()
	v.sign = make([]byte, v.hash.Size())
}

func (v *RSVerifier) Verify(data string, signature []byte) error {
	b := string2bytes(data)
	// Hash data.
	v.hash.Reset()
	v.hash.Write(b)
	v.hash.Sum(v.sign[:0])
	// RSA verify.
	return rsa.VerifyPKCS1v15(v.key, v.crypto, v.sign, signature)
}
