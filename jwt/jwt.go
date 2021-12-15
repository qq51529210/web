package jwt

import (
	"bytes"
	"crypto"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"unsafe"
)

var (
	errInvalidJWT  = errors.New("invalid jwt")
	errAlgType     = errors.New("alg must be string type")
	errAlgNotfound = errors.New("algorithm not found")
	decoderPool    = sync.Pool{}
)

func init() {
	decoderPool.New = func() interface{} {
		d := new(decoder)
		d.Decoder = json.NewDecoder(&d.Buffer)
		return d
	}
}

// Standard JWT Claims.
const (
	TYP = "typ"
	ALG = "alg"
	ISS = "iss"
	SUB = "sub"
	AUD = "aud"
	EXP = "exp"
	NBF = "nbf"
	IAT = "iat"
	JTI = "jti"
)

// Algorithm.
type Alg string

func (alg Alg) CryptoHash() crypto.Hash {
	switch alg {
	case HS256:
		return crypto.SHA256
	case HS384:
		return crypto.SHA384
	case HS512:
		return crypto.SHA512
	case PS256:
		return crypto.SHA256
	case PS384:
		return crypto.SHA384
	case PS512:
		return crypto.SHA512
	case ES256:
		return crypto.SHA256
	case ES384:
		return crypto.SHA384
	case ES512:
		return crypto.SHA512
	case RS256:
		return crypto.SHA256
	case RS384:
		return crypto.SHA384
	case RS512:
		return crypto.SHA512
	default:
		return 0
	}
}

// All algorithms.
const (
	HS256 Alg = "HS256"
	HS384 Alg = "HS384"
	HS512 Alg = "HS512"
	PS256 Alg = "PS256"
	PS384 Alg = "PS384"
	PS512 Alg = "PS512"
	ES256 Alg = "ES256"
	ES384 Alg = "ES384"
	ES512 Alg = "ES512"
	RS256 Alg = "RS256"
	RS384 Alg = "RS384"
	RS512 Alg = "RS512"
)

// Zero copy reference
func string2bytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

type Generator interface {
	// Generate JWT and return.
	Generate(header, payload map[string]interface{}) (string, error)
}

// Use for encode jwt header and payload.
type encoder struct {
	*json.Encoder
	// json.Encoder output writer.
	bytes.Buffer
	// Buffer for jwt.
	token []byte
}

// Initialize json encoder.
func (e *encoder) Init() {
	e.Encoder = json.NewEncoder(&e.Buffer)
}

// Encode jwt header and payload into buffer.
func (e *encoder) Enc(header, payload map[string]interface{}) error {
	header[TYP] = "JWT"
	// Reset token buffer.
	e.token = e.token[:0]
	// Encode header to JSON.
	e.Buffer.Reset()
	err := e.Encoder.Encode(header)
	if err != nil {
		return err
	}
	// Base64 header JSON
	e.Base64(e.Buffer.Bytes())
	// '.' between header and payload.
	e.token = append(e.token, '.')
	// Encode payload to JSON.
	e.Buffer.Reset()
	err = e.Encoder.Encode(payload)
	if err != nil {
		return err
	}
	// Base64 payload JSON
	e.Base64(e.Buffer.Bytes())
	return nil
}

// Base64 encode data d then append to e.token.
func (e *encoder) Base64(d []byte) {
	i := len(e.token)
	m := i + base64.RawURLEncoding.EncodedLen(len(d))
	if m > cap(e.token) {
		b := make([]byte, m)
		copy(b, e.token)
		e.token = b
	} else {
		e.token = e.token[:m]
	}
	base64.RawURLEncoding.Encode(e.token[i:], d)
}

type Verifier interface {
	// Verify signature by data.
	// Param "data" is JWT data part(header.payload, base64).
	// Param "signature" is JWT signature part(base64 decoded).
	Verify(data string, signature []byte) error
}

// Verifier table, key is alg.
type Verifiers map[Alg]Verifier

// Use for encode jwt header and payload.
type decoder struct {
	*json.Decoder
	bytes.Buffer
	// Buffer for base64 decode
	buff []byte
}

// Base64 decode s to d.sign.
func (d *decoder) Base64(s string) error {
	b := string2bytes(s)
	// Base64 decode
	n := base64.RawURLEncoding.DecodedLen(len(b))
	if cap(d.buff) < n {
		d.buff = make([]byte, n)
	} else {
		d.buff = d.buff[:n]
	}
	_, err := base64.RawURLEncoding.Decode(d.buff, b)
	return err
}

// Base64 decode then json decode
func (d *decoder) Dec(s string) (map[string]interface{}, error) {
	// Base64 decode
	err := d.Base64(s)
	if err != nil {
		return nil, err
	}
	// Json decode
	data := make(map[string]interface{})
	d.Buffer.Reset()
	d.Buffer.Write(d.buff)
	err = d.Decoder.Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// Verify token signature, return header and payload.
// Function verifier return a Verifier to verify token(You may have a Verifier pool),
// return nil means does not support token's algorihm, and function Verify will return error.
func Verify(token string, verifier func(Alg) Verifier) (map[string]interface{}, map[string]interface{}, error) {
	i1, i2, err := SplitJWT(token)
	if err != nil {
		return nil, nil, err
	}
	decoder := decoderPool.Get().(*decoder)
	// Decode header
	header, err := decoder.Dec(token[:i1])
	if err != nil {
		decoderPool.Put(decoder)
		return nil, nil, err
	}
	// Alg
	val, ok := header[ALG]
	if !ok {
		decoderPool.Put(decoder)
		return nil, nil, errAlgNotfound
	}
	alg, ok := val.(string)
	if !ok {
		decoderPool.Put(decoder)
		return nil, nil, errAlgType
	}
	// Verifiier to verify.
	ver := verifier(Alg(alg))
	if ver == nil {
		decoderPool.Put(decoder)
		return nil, nil, fmt.Errorf("unsupported algorithm %s", alg)
	}
	// Base64 decode signature.
	err = decoder.Base64(token[i2+1:])
	if err != nil {
		decoderPool.Put(decoder)
		return nil, nil, err
	}
	// Verify.
	err = ver.Verify(token[:i2], decoder.buff)
	if err != nil {
		decoderPool.Put(decoder)
		return nil, nil, err
	}
	// Decode payload
	payload, err := decoder.Dec(token[i1+1 : i2])
	if err != nil {
		decoderPool.Put(decoder)
		return nil, nil, err
	}
	decoderPool.Put(decoder)
	return header, payload, nil
}

// Split JWT, return index of '.'.
func SplitJWT(token string) (int, int, error) {
	// '.' follow header.
	i1 := strings.IndexByte(token, '.')
	if i1 < 0 {
		return 0, 0, errInvalidJWT
	}
	// '.' follow payload.
	i2 := strings.IndexByte(token[i1+1:], '.') + 1
	if i2 == 0 || i2 == len(token) {
		return 0, 0, errInvalidJWT
	}
	return i1, i1 + i2, nil
}
