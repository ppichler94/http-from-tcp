package request

import (
	"errors"
	"io"
	"slices"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\r\n")
	if len(lines) < 1 {
		return nil, errors.New("empty request")
	}

	requestLine, err := parseRequestLine(lines[0])
	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: requestLine,
	}, nil
}

func parseRequestLine(line string) (RequestLine, error) {
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return RequestLine{}, errors.New("invalid request line")
	}

	if !strings.HasPrefix(parts[2], "HTTP/") {
		return RequestLine{}, errors.New("invalid HTTP version format")
	}

	version := strings.TrimPrefix(parts[2], "HTTP/")
	if version != "1.1" {
		return RequestLine{}, errors.New("unsupported HTTP version")
	}

	methods := []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS", "PATCH"}
	if !slices.Contains(methods, parts[0]) {
		return RequestLine{}, errors.New("unsupported HTTP method")
	}

	return RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   version,
	}, nil
}
