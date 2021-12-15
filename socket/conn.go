package socket

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Code byte

func (c Code) String() string {
	switch c {
	case CodeText:
		return "text"
	case CodeBinary:
		return "binary"
	case CodeClose:
		return "close"
	case CodePing:
		return "ping"
	case CodePong:
		return "pong"
	default:
		return "continuation"
	}
}

const (
	codeContinuation Code = 0
	CodeText         Code = 1
	CodeBinary       Code = 2
	CodeClose        Code = 8
	CodePing         Code = 9
	CodePong         Code = 10
)

var (
	_Fin  = [2]byte{0x00, 0x80}
	_Mask = [2]byte{0x00, 0x80}
)

// Mask data with key
func maskData(data, key []byte) {
	for i := 0; i < len(data); i++ {
		data[i] ^= key[i%4]
	}
}

type Conn struct {
	conn io.ReadWriteCloser
	mask byte
}

// Write code type data.
// If data length bigger than payload, it will be split into multiple frames.
func (c *Conn) Write(code Code, data []byte, payload int) error {
	if len(data) <= payload {
		return c.writeFrame(_Fin[1], code, data)
	}
	switch code {
	case CodeBinary, CodeText:
	default:
		return fmt.Errorf(`could not split "%s" data into frames`, code.String())
	}
	// If length of data bigger than c.maxSize[1] split into frames.
	p := data
	// First frame, fin=0, code!=0.
	err := c.writeFrame(_Fin[0], code, p[:payload])
	if err != nil {
		return err
	}
	p = p[payload:]
	// Continuation frame, fin=0, code=0.
	for len(p) > payload {
		err := c.writeFrame(_Fin[0], codeContinuation, p[:payload])
		if err != nil {
			return err
		}
		p = p[payload:]
	}
	// The last frame, fin=1, opcode=0.
	return c.writeFrame(_Fin[1], codeContinuation, p)
}

// Write a frame.
func (c *Conn) writeFrame(fin byte, code Code, data []byte) error {
	b := encodeBufferPool.Get().(*encodeBuffer)
	b.Reset()
	// Encode fin and code.
	b.Put8(fin | byte(code))
	// Encode mask and payload length.
	payload := len(data)
	if payload < 126 {
		b.Put8(c.mask | byte(payload))
	} else if payload <= 0xffff {
		b.Put8(c.mask | 126)
		b.Put16(uint16(payload))
	} else {
		b.Put8(c.mask | 127)
		b.Put64(uint64(payload))
	}
	// Encode mask key.
	if c.mask != 0 {
		// Append random key.
		keyIdx1 := b.len
		b.PutRandom(4)
		keyIdx2 := b.len
		// Append mask data.
		dataIdx := b.len
		b.PutBytes(data)
		maskData(b.buf[dataIdx:], b.buf[keyIdx1:keyIdx2])
		// Write buffer.
		_, err := c.conn.Write(b.buf[:b.len])
		encodeBufferPool.Put(b)
		return err
	}
	// Write header.
	_, err := c.conn.Write(b.buf[:b.len])
	if err != nil {
		encodeBufferPool.Put(b)
		return err
	}
	encodeBufferPool.Put(b)
	// Write data.
	_, err = c.conn.Write(data)
	return err
}

// Read a complete message then call handle.
// If message length bigger than maxLen, it return error.
func (c *Conn) ReadLoop(maxLen int, handle func(Code, []byte) error) error {
	var (
		fin     byte
		code    Code
		mask    byte
		length  int
		err     error
		header  [8]byte
		key     [4]byte
		payload [16]readBuffer
		b       *readBuffer
	)
	for {
		// Read header.
		_, err = io.ReadFull(c.conn, header[:2])
		if err != nil {
			return err
		}
		// Decode fin
		fin = header[0] & _Fin[1]
		// Decode code
		code = Code(header[0] & 0x0f)
		// Decode mask
		mask = header[1] & _Mask[1]
		// Decode length
		length = int(header[1] & 0x7f)
		switch length {
		case 126:
			_, err = io.ReadFull(c.conn, header[:2])
			if err != nil {
				return err
			}
			length = int(binary.BigEndian.Uint16(header[:]))
		case 127:
			_, err = io.ReadFull(c.conn, header[:8])
			if err != nil {
				return err
			}
			length = int(binary.BigEndian.Uint64(header[:]))
		}
		// Decode key
		if mask != 0 {
			_, err = io.ReadFull(c.conn, key[:])
			if err != nil {
				return err
			}
		}
		// Read payload.
		b = &payload[code]
		idx := b.len
		err = b.ReadN(c.conn, length)
		if err != nil {
			return err
		}
		// Mask payload data.
		if mask != 0 {
			maskData(b.buf[idx:b.len], key[:])
		}
		// Last frame of a message
		if fin == _Fin[1] {
			// Call back handle.
			handle(code, b.buf[:b.len])
			// Reset buffer.
			b.Reset()
		}
	}
}

// Write a CodeClose frame, then call Closer.
func (c *Conn) Close() error {
	c.Write(CodeClose, nil, 128)
	return c.conn.Close()
}
