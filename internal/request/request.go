package request

import (
	"errors"
	"io"
	"slices"
	"strconv"
	"strings"

	"http.ppichler94.io/internal/headers"
)

const bufferSize = 8

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        []byte
	state       requestState
}

type requestState int

const (
	Initialized    requestState = 0
	ParsingHeaders requestState = 1
	ParsingBody    requestState = 2
	Done           requestState = 3
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize)
	readToIndex := 0
	request := &Request{
		state:   Initialized,
		Headers: headers.NewHeaders(),
	}

	for request.state != Done {
		n, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if err == io.EOF {
				if request.state != Done {
					return nil, errors.New("incomplete request")
				}
				request.state = Done
				return request, nil
			}
		}
		readToIndex += n
		if readToIndex == len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		bytesRead, err := request.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}
		if bytesRead > 0 {
			// Shift remaining bytes to the beginning of the buffer
			copy(buf, buf[bytesRead:])
			readToIndex -= bytesRead
		}
	}

	return request, nil
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.state != Done {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return totalBytesParsed, nil
		}
		totalBytesParsed += n
	}
	return 0, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.state {
	case Initialized:
		line, bytesRead, err := parseRequestLine(string(data))
		if err != nil {
			return 0, err
		}
		if bytesRead == 0 {
			return 0, nil
		}
		r.RequestLine = line
		r.state = ParsingHeaders
		return bytesRead, nil
	case ParsingHeaders:
		bytesRead, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			r.state = ParsingBody
		}
		return bytesRead, nil
	case ParsingBody:
		if r.Headers.Get("Content-Length") == "" {
			r.state = Done
			return 0, nil
		}
		bodyLength, err := strconv.Atoi(r.Headers.Get("Content-Length"))
		if err != nil {
			return 0, err
		}
		if len(data) == bodyLength {
			r.Body = data[:bodyLength]
			r.state = Done
			return bodyLength, nil
		}
		if len(data) > bodyLength {
			return 0, errors.New("body too long")
		}
	}

	return 0, nil
}

func parseRequestLine(line string) (RequestLine, int, error) {

	bytesRead := strings.Index(line, "\r\n")
	if bytesRead <= 0 {
		return RequestLine{}, 0, nil // Incomplete line
	}

	parts := strings.Split(line[:bytesRead], " ")
	if len(parts) != 3 {
		return RequestLine{}, 0, errors.New("invalid request line")
	}

	if !strings.HasPrefix(parts[2], "HTTP/") {
		return RequestLine{}, 0, errors.New("invalid HTTP version format")
	}

	version := strings.TrimPrefix(parts[2], "HTTP/")
	if version != "1.1" {
		return RequestLine{}, 0, errors.New("unsupported HTTP version")
	}

	methods := []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS", "PATCH"}
	if !slices.Contains(methods, parts[0]) {
		return RequestLine{}, 0, errors.New("unsupported HTTP method")
	}

	return RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   version,
	}, bytesRead + 2, nil
}
