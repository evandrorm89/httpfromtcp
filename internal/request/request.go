package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	State       RequestState
}

func newRequest() *Request {
	return &Request{
		State: StateInitialized,
	}
}

func (r *Request) done() bool {
	return r.State == StateDone
}

var (
	ErrBadRequestLine = errors.New("bad request line")
	ErrInvalidMethod  = errors.New("invalid method")
	ErrInvalidVersion = errors.New("invalid http version")
	ErrDoneState      = errors.New("trying to read data in a done state")
	ErrUnknownState   = errors.New("unkown State")
	SEPARATOR         = []byte("\r\n")
)

type RequestState int

const (
	StateInitialized RequestState = iota
	StateDone
)

const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	buff := make([]byte, bufferSize, bufferSize)
	readToIndex := 0
	r := newRequest()

	for !r.done() {
		if readToIndex >= len(buff) {
			newBuff := make([]byte, len(buff)*2)
			copy(newBuff, buff)
			buff = newBuff
		}

		n, err := reader.Read(buff[readToIndex:])
		if err != nil {
			if err == io.EOF {
				r.State = StateDone
				break
			}
			return nil, err
		}

		readToIndex += n
		readN, err := r.parse(buff[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(buff, buff[readN:readToIndex])
		readToIndex -= readN
	}

	return r, nil
}

func parseRequestLine(b []byte) (*RequestLine, int, error) {
	idx := bytes.Index(b, SEPARATOR)
	if idx == -1 {
		return nil, 0, nil
	}

	startLine := b[:idx]
	read := idx + len(SEPARATOR)

	parts := bytes.Split(startLine, []byte(" "))
	if len(parts) != 3 {
		return nil, 0, ErrBadRequestLine
	}

	httpParts := bytes.Split(parts[2], []byte("/"))
	if len(httpParts) != 2 || string(httpParts[0]) != "HTTP" || string(httpParts[1]) != "1.1" {
		return nil, 0, ErrInvalidVersion
	}

	return &RequestLine{
		HttpVersion:   string(httpParts[1]),
		RequestTarget: string(parts[1]),
		Method:        string(parts[0]),
	}, read, nil
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for {
		switch r.State {
		case StateInitialized:
			rl, n, err := parseRequestLine(data[read:])
			if err != nil {
				return 0, err
			}
			if n == 0 {
				break outer
			}

			r.RequestLine = *rl
			read += n

			fmt.Println("Parsed Request Line:", r.RequestLine)
			fmt.Println("read:", r)

			r.State = StateDone
		case StateDone:
			break outer
		}
	}
	return read, nil
}
