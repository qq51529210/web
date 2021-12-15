package socket

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

var (
	serverRequiredHeader = map[string]string{
		"Upgrade":               "websocket",
		"Connection":            "Upgrade",
		"Sec-WebSocket-Version": "13",
	}
	clientRequiredHeader = map[string]string{
		"Upgrade":    "websocket",
		"Connection": "Upgrade",
	}
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// Implement this interface to handle message.
type Handler interface {
	HandleText(*Conn, []byte)
	HandleBinary(*Conn, []byte)
	HandleClose(*Conn, []byte)
	HandlePing(*Conn, []byte)
	HandlePong(*Conn, []byte)
}

// Server side connection.
func Accept(res http.ResponseWriter, req *http.Request) (*Conn, error) {
	// Check required headers
	err := checkRequiredHeader(req.Header, serverRequiredHeader)
	if err != nil {
		return nil, err
	}
	// Check header
	key := req.Header.Get("Sec-Websocket-Key")
	if key == "" {
		return nil, fmt.Errorf(`header "Sec-Websocket-Key" is required`)
	}
	// Set response headers.
	res.Header().Set("Upgrade", "websocket")
	res.Header().Set("Connection", "Upgrade")
	res.Header().Set("Sec-WebSocket-Accept", GenSecWebSocketAccept(key))
	// Set status code.
	res.WriteHeader(http.StatusSwitchingProtocols)
	// Hijack net.Conn
	h, o := res.(http.Hijacker)
	if !o {
		return nil, errors.New("hijacker response conn fail")
	}
	conn, buf, err := h.Hijack()
	if nil != err {
		return nil, err
	}
	// 响应给客户端
	err = buf.Flush()
	if nil != err {
		return nil, err
	}
	// Conn
	return &Conn{conn: conn, mask: _Mask[0]}, nil
}

// Client side connection.
func Dial(req *http.Request, conn io.ReadWriteCloser) (*Conn, error) {
	// Set required header
	secWebSocketKey := GenSecWebSocketKey()
	req.Header.Set("Connection", "Upgrade")
	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("Sec-WebSocket-Version", "13")
	req.Header.Set("Sec-WebSocket-Key", secWebSocketKey)
	// Write http request.
	err := req.Write(conn)
	if err != nil {
		return nil, err
	}
	// Read http response.
	rd := bufio.NewReader(conn)
	res, err := http.ReadResponse(rd, req)
	if err != nil {
		return nil, err
	}
	// Check status code.
	if http.StatusSwitchingProtocols != res.StatusCode {
		return nil, fmt.Errorf("status code %d", res.StatusCode)
	}
	// Check required header.
	err = checkRequiredHeader(res.Header, clientRequiredHeader)
	if err != nil {
		return nil, err
	}
	// Check "Sec-WebSocket-Accept" value.
	key := req.Header.Get("Sec-Websocket-Accept")
	if key != GenSecWebSocketAccept(secWebSocketKey) {
		return nil, fmt.Errorf(`invalid "Sec-Websocket-Accept" value %s`, key)
	}
	// Conn
	return &Conn{conn: conn, mask: _Mask[1]}, nil
}

func GenSecWebSocketAccept(webSocketKey string) string {
	var buf bytes.Buffer
	buf.WriteString(webSocketKey)
	buf.WriteString("258EAFA5-E914-47DA-95CA-C5AB0DC85B11")
	hash := sha1.New()
	hash.Write(buf.Bytes())
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func GenSecWebSocketKey() string {
	var b [16]byte
	random.Read(b[:])
	return base64.StdEncoding.EncodeToString(b[:])
}

func checkRequiredHeader(request http.Header, required map[string]string) error {
	for k, v := range required {
		vs := request.Values(k)
		for i := 0; i < len(vs); i++ {
			if vs[i] == v {
				return nil
			}
		}
		return fmt.Errorf(`header "%s: %s" is required`, k, v)
	}
	return nil
}
