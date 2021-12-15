package jwt

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"hash"
	"sync"
)

var (
	ps256GenPool sync.Pool
	ps384GenPool sync.Pool
	ps512GenPool sync.Pool
)

func init() {
	ps256GenPool.New = func() interface{} {
		return NewPSGenerator(HS256, nil, nil)
	}
	ps384GenPool.New = func() interface{} {
		return NewPSGenerator(HS384, nil, nil)
	}
	ps512GenPool.New = func() interface{} {
		return NewPSGenerator(HS512, nil, nil)
	}
}

func NewPSGenerator(alg Alg, key *rsa.PrivateKey, opt *rsa.PSSOptions) *PSGenerator {
	s := new(PSGenerator)
	s.Init(alg, key, opt)
	return s
}

func NewPSVerifier(alg Alg, key *rsa.PublicKey, opt *rsa.PSSOptions) *PSVerifier {
	s := new(PSVerifier)
	s.Init(alg, key, opt)
	return s
}

func GeneratePS256(header, payload map[string]interface{}, key *rsa.PrivateKey, opt *rsa.PSSOptions) (string, error) {
	return generatePS(header, payload, key, opt, &ps256GenPool)
}

func GeneratePS384(header, payload map[string]interface{}, key *rsa.PrivateKey, opt *rsa.PSSOptions) (string, error) {
	return generatePS(header, payload, key, opt, &ps384GenPool)
}

func GeneratePS512(header, payload map[string]interface{}, key *rsa.PrivateKey, opt *rsa.PSSOptions) (string, error) {
	return generatePS(header, payload, key, opt, &ps512GenPool)
}

func generatePS(header, payload map[string]interface{}, key *rsa.PrivateKey, opt *rsa.PSSOptions, pool *sync.Pool) (string, error) {
	g := pool.Get().(*PSGenerator)
	g.key = key
	g.opt = opt
	token, err := g.Generate(header, payload)
	pool.Put(g)
	return token, err
}

type PSGenerator struct {
	alg    Alg
	enc    encoder
	key    *rsa.PrivateKey
	opt    *rsa.PSSOptions
	hash   hash.Hash
	crypto crypto.Hash
	sign   []byte
}

func (g *PSGenerator) Init(alg Alg, key *rsa.PrivateKey, opt *rsa.PSSOptions) {
	g.alg = alg
	g.key = key
	g.opt = opt
	g.crypto = alg.CryptoHash()
	g.hash = g.crypto.New()
	g.sign = make([]byte, g.hash.Size())
	g.enc.Init()
}

func (g *PSGenerator) Generate(header, payload map[string]interface{}) (string, error) {
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
	data, err := rsa.SignPSS(rand.Reader, g.key, g.crypto, g.sign, g.opt)
	if err != nil {
		return "", err
	}
	// '.' between payload and signature.
	g.enc.token = append(g.enc.token, '.')
	// Base64 hash signature.
	g.enc.Base64(data)
	return string(g.enc.token), nil
}

type PSVerifier struct {
	alg    Alg
	key    *rsa.PublicKey
	opt    *rsa.PSSOptions
	hash   hash.Hash
	crypto crypto.Hash
	sign   []byte
}

func (v *PSVerifier) Init(alg Alg, key *rsa.PublicKey, opt *rsa.PSSOptions) {
	v.alg = alg
	v.key = key
	v.opt = opt
	v.crypto = alg.CryptoHash()
	v.hash = v.crypto.New()
	v.sign = make([]byte, v.hash.Size())
}

func (v *PSVerifier) Verify(data string, signature []byte) error {
	b := string2bytes(data)
	// Hash data.
	v.hash.Reset()
	v.hash.Write(b)
	v.hash.Sum(v.sign[:0])
	// RSA PSS verify.
	return rsa.VerifyPSS(v.key, v.crypto, v.sign, signature, v.opt)
}
