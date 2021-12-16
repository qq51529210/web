package router

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type HandleFunc func(ctx *Context)

var (
	errSeekOffset = errors.New("seek: invalid offset")
	// Create compressor functions.
	compressFunc = []func(io.Writer) io.WriteCloser{
		func(w io.Writer) io.WriteCloser {
			return gzip.NewWriter(w)
		},
		func(w io.Writer) io.WriteCloser {
			return zlib.NewWriter(w)
		},
		func(w io.Writer) io.WriteCloser {
			wc, _ := flate.NewWriter(w, flate.DefaultCompression)
			return wc
		},
	}
	compressName = []string{
		"gzip",
		"zlib",
		"deflate",
	}
)

const (
	gzipCompress = iota
	zlibCompress
	deflateCompress
)

type FileHandler struct {
	file string
}

func (h *FileHandler) Handle(ctx *Context) {
	http.ServeFile(ctx.ResponseWriter, ctx.Request, h.file)
	ctx.Next()
}

func NewFileHandler(file string) (*FileHandler, error) {
	fileInfo, err := os.Stat(file)
	if err != nil {
		return nil, err
	}
	if fileInfo.IsDir() {
		return nil, fmt.Errorf("%s is a directory", file)
	}
	return &FileHandler{file: file}, nil
}

type dataSeeker struct {
	b []byte
	i int64
}

func (s *dataSeeker) Seek(o int64, w int) (int64, error) {
	switch w {
	case io.SeekStart:
	case io.SeekCurrent:
		o += s.i
	case io.SeekEnd:
		o += int64(len(s.b))
	}
	if o < 0 {
		return 0, errSeekOffset
	}
	if o > int64(len(s.b)) {
		o = int64(len(s.b))
	}
	s.i = o
	return o, nil
}

func (s *dataSeeker) Read(p []byte) (int, error) {
	if s.i >= int64(len(s.b)) {
		return 0, io.EOF
	}
	n := copy(p, s.b[s.i:])
	s.i += int64(n)
	return n, nil
}

func (s *dataSeeker) Write(b []byte) (int, error) {
	s.b = append(s.b, b...)
	return len(b), nil
}

type CacheHandler struct {
	contentType    string
	modTime        time.Time
	data           []byte
	compressedData [3][]byte
}

func (h *CacheHandler) Handle(ctx *Context) {
	if h.contentType != "" {
		ctx.ResponseWriter.Header().Add("Content-Type", h.contentType)
	}
	for _, s := range strings.Split(ctx.Request.Header.Get("Accept-Encoding"), ",") {
		switch s {
		case "*", "gzip":
			h.serveContent(ctx, gzipCompress)
			ctx.Next()
			return
		case "zlib":
			h.serveContent(ctx, zlibCompress)
			ctx.Next()
			return
		case "deflate":
			h.serveContent(ctx, deflateCompress)
			ctx.Next()
			return
		default:
			continue
		}
	}
	// Handler does not has client compressions.
	http.ServeContent(ctx.ResponseWriter, ctx.Request, "", h.modTime, &dataSeeker{b: h.data})
	ctx.Next()
}

func (h *CacheHandler) serveContent(ctx *Context, n int) {
	// Compress data if is empty.
	if len(h.compressedData[n]) < 1 {
		var buf bytes.Buffer
		w := compressFunc[n](&buf)
		w.Write(h.data)
		w.Close()
		h.compressedData[n] = append(h.compressedData[n], buf.Bytes()...)
	}
	// Response compressed data.
	if len(h.compressedData[n]) < len(h.data) {
		ctx.ResponseWriter.Header().Add("Content-Encoding", compressName[n])
		http.ServeContent(ctx.ResponseWriter, ctx.Request, "", h.modTime, &dataSeeker{b: h.compressedData[n]})
		return
	}
	// Response origin data.
	http.ServeContent(ctx.ResponseWriter, ctx.Request, "", h.modTime, &dataSeeker{b: h.data})
}

func NewCacheHandlerFromFile(file string) (*CacheHandler, error) {
	fileInfo, err := os.Stat(file)
	if err != nil {
		return nil, err
	}
	if fileInfo.IsDir() {
		return nil, fmt.Errorf("%s is a directory", file)
	}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return &CacheHandler{
		contentType: mime.TypeByExtension(filepath.Ext(file)),
		modTime:     fileInfo.ModTime(),
		data:        data,
	}, nil
}
