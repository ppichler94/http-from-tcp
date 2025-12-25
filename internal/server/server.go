package server

import (
	"fmt"
	"net"
	"sync/atomic"

	"http.ppichler94.io/internal/request"
	"http.ppichler94.io/internal/response"
)

type Server struct {
	listener net.Listener
	handler  Handler
	closed   atomic.Bool
}

func Serve(port int, handler Handler) (*Server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	s := &Server{listener: l, handler: handler}

	go s.listen()

	return s, nil
}

func (s *Server) Close() error {
	err := s.listener.Close()
	s.closed.Store(true)
	return err
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	req, err := request.RequestFromReader(conn)
	if err != nil {
		hErr := &HandlerError{Status: response.BadRequest, Message: err.Error()}
		hErr.write(conn)
		return
	}

	writer := response.NewWriter(conn)

	s.handler(&writer, req)
}
