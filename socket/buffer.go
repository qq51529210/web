package socket

import (
	"encoding/binary"
	"io"
	"sync"
)

var (
	readBufferPool   sync.Pool
	encodeBufferPool sync.Pool
)

func init() {
	encodeBufferPool.New = func() interface{} {
		return &encodeBuffer{}
	}
	readBufferPool.New = func() interface{} {
		return &readBuffer{}
	}
}

// Use for read a segmented message.
type readBuffer struct {
	buf []byte
	len int
}

// Make sure b can stores n bytes data.
func (b *readBuffer) grow(n int) {
	m := b.len + n
	if m > len(b.buf) {
		nb := make([]byte, m)
		copy(nb, b.buf)
		b.buf = nb
	}
}

func (b *readBuffer) Reset() {
	b.len = 0
}

func (b *readBuffer) ReadN(r io.Reader, n int) error {
	b.grow(n)
	n, err := io.ReadFull(r, b.buf[b.len:b.len+n])
	if err != nil {
		return err
	}
	b.len += n
	return nil
}

// Use for encode message.
type encodeBuffer struct {
	buf []byte
	len int
}

func (b *encodeBuffer) Reset() {
	b.len = 0
}

// Make sure b can stores n bytes data.
func (b *encodeBuffer) grow(n int) {
	m := b.len + n
	if m > len(b.buf) {
		nb := make([]byte, m)
		copy(nb, b.buf)
		b.buf = nb
	}
}

func (b *encodeBuffer) Put8(n byte) {
	b.grow(1)
	b.buf[b.len] = n
	b.len++
}

// Append BigEndian n
func (b *encodeBuffer) Put16(n uint16) {
	b.grow(2)
	binary.BigEndian.PutUint16(b.buf[b.len:], n)
	b.len += 2
}

// Append BigEndian n
func (b *encodeBuffer) Put32(n uint32) {
	b.grow(4)
	binary.BigEndian.PutUint32(b.buf[b.len:], n)
	b.len += 4
}

// Append BigEndian n
func (b *encodeBuffer) Put64(n uint64) {
	b.grow(8)
	binary.BigEndian.PutUint64(b.buf[b.len:], n)
	b.len += 8
}

// Append n random bytes.
func (b *encodeBuffer) PutRandom(n int) {
	b.grow(n)
	random.Read(b.buf[b.len:])
	b.len += n
}

// Append bytes.
func (b *encodeBuffer) PutBytes(d []byte) {
	b.grow(len(d))
	b.len += copy(b.buf[b.len:], d)
}

// Append string.
func (b *encodeBuffer) PutString(s string) {
	b.grow(len(s))
	b.len += copy(b.buf[b.len:], s)
}
