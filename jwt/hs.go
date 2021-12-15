package jwt

import (
	"crypto/hmac"
	"hash"
)

func NewHSGenerator(alg Alg, key []byte) *HSGenerator {
	s := new(HSGenerator)
	s.Init(alg, key)
	return s
}

func NewHSVerifier(alg Alg, key []byte) *HsVerifier {
	s := new(HsVerifier)
	s.Init(alg, key)
	return s
}

func GenerateHS256(header, payload map[string]interface{}, key []byte) (string, error) {
	return NewHSGenerator(HS256, key).Generate(header, payload)
}

func GenerateHS384(header, payload map[string]interface{}, key []byte) (string, error) {
	return NewHSGenerator(HS384, key).Generate(header, payload)
}

func GenerateHS512(header, payload map[string]interface{}, key []byte) (string, error) {
	return NewHSGenerator(HS512, key).Generate(header, payload)
}

type HSGenerator struct {
	enc  encoder
	alg  Alg
	hash hash.Hash
	sign []byte
}

func (g *HSGenerator) Init(alg Alg, key []byte) {
	g.alg = alg
	g.hash = hmac.New(alg.CryptoHash().New, key)
	g.sign = make([]byte, g.hash.Size())
	g.enc.Init()
}

func (g *HSGenerator) Generate(header, payload map[string]interface{}) (string, error) {
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
	// '.' between payload and signature.
	g.enc.token = append(g.enc.token, '.')
	// Base64 hash signature.
	g.enc.Base64(g.sign)
	return string(g.enc.token), nil
}

type HsVerifier struct {
	alg  Alg
	hash hash.Hash
	sign []byte
}

func (v *HsVerifier) Init(alg Alg, key []byte) {
	v.alg = alg
	v.hash = hmac.New(alg.CryptoHash().New, key)
	v.sign = make([]byte, v.hash.Size())
}

func (v *HsVerifier) Verify(data string, signature []byte) error {
	b := string2bytes(data)
	// Hash data.
	v.hash.Reset()
	v.hash.Write(b)
	v.hash.Sum(v.sign[:0])
	// Compare data with signature.
	if len(v.sign) != len(signature) {
		return errInvalidJWT
	}
	for i := 0; i < len(v.sign); i++ {
		if v.sign[i] != signature[i] {
			return errInvalidJWT
		}
	}
	return nil
}
