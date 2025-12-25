package response

import (
	"errors"
	"io"

	"http.ppichler94.io/internal/headers"
)

type Writer struct {
	w     io.Writer
	state writerState
}

type writerState int

const (
	StatusLine writerState = iota
	Headers
	Body
)

func NewWriter(w io.Writer) Writer {
	return Writer{w, StatusLine}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.state != StatusLine {
		return errors.New("The status line must be written first")
	}

	err := WriteStatusLine(w.w, statusCode)
	if err != nil {
		return err
	}

	w.state = Headers
	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.state != Headers {
		return errors.New("The headers must be written after the status line and before the body")
	}

	err := WriteHeaders(w.w, headers)
	if err != nil {
		return err
	}

	w.state = Body
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.state != Body {
		return 0, errors.New("The body must be written after the headers")
	}

	return w.w.Write(p)
}
