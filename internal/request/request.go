package request

import (
	"bytes"
	"errors"
	"io"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

var (
	ErrBadRequestLine = errors.New("bad request line")
	ErrInvalidMethod  = errors.New("invalid method")
	ErrInvalidVersion = errors.New("invalid http version")
)

func RequestFromReader(reader io.Reader) (*Request, error) {
	b, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	rl, err := parseRequestLine(b)
	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: rl,
	}, nil
}

func parseRequestLine(b []byte) (RequestLine, error) {
	lines := bytes.Split(b, []byte("\r\n"))

	parts := bytes.Split(lines[0], []byte(" "))
	if len(parts) != 3 {
		return RequestLine{}, ErrBadRequestLine
	}

	method := string(parts[0])
	if method != "GET" && method != "POST" {
		return RequestLine{}, ErrInvalidMethod
	}
	target := string(parts[1])
	versionParts := bytes.Split(parts[2], []byte("/"))
	version := string(versionParts[1])
	if version != "1.1" {
		return RequestLine{}, ErrInvalidVersion
	}

	return RequestLine{
		HttpVersion:   version,
		RequestTarget: target,
		Method:        method,
	}, nil
}
