package headers

import (
	"errors"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func isToken(s string) bool {
	for _, c := range s {
		found := false
		if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') {
			found = true
		}

		switch c {
		case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~':
			found = true
		}

		if !found {
			return false
		}
	}
	return true
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	bytesRead := strings.Index(string(data), "\r\n")
	if bytesRead == -1 {
		return 0, false, nil // Incomplete line
	}

	if strings.HasPrefix(string(data), "\r\n") {
		return bytesRead + 2, true, nil // End of headers
	}

	headerLine := string(data[:bytesRead])
	parts := strings.SplitN(headerLine, ":", 2)
	if len(parts) != 2 {
		return 0, false, errors.New("invalid header format")
	}

	key := strings.TrimLeft(parts[0], " ")
	key = strings.ToLower(key)
	if !isToken(key) {
		return 0, false, errors.New("invalid header key")
	}
	value := strings.TrimSpace(parts[1])

	if strings.HasSuffix(key, " ") {
		return 0, false, errors.New("invalid header spacing")
	}

	if v, contains := h[key]; contains {
		h[key] = v + ", " + value
	} else {
		h[key] = value
	}

	return bytesRead + 2, false, nil

}
