package server

import (
	"fmt"
	"io"

	"http.ppichler94.io/internal/request"
	"http.ppichler94.io/internal/response"
)

type HandlerError struct {
	Message string
	Status  response.StatusCode
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func (error HandlerError) write(w io.Writer) {
	headers := response.GetDefaultHeaders(len(error.Message))
	err := response.WriteStatusLine(w, error.Status)
	if err != nil {
		fmt.Println("Error writing status line:", err)
		return
	}
	err = response.WriteHeaders(w, headers)
	if err != nil {
		fmt.Println("Error writing headers:", err)
		return
	}
	_, err = w.Write([]byte(error.Message))
	if err != nil {
		fmt.Println("Error writing response body:", err)
		return
	}
}
