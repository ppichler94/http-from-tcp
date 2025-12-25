package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"http.ppichler94.io/internal/request"
	"http.ppichler94.io/internal/response"
	"http.ppichler94.io/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w io.Writer, req *request.Request) *server.HandlerError {
	if req.RequestLine.RequestTarget == "/yourproblem" {
		return &server.HandlerError{
			Message: "Your problem is not my problem\n",
			Status:  response.BadRequest,
		}
	}

	if req.RequestLine.RequestTarget == "/myproblem" {
		return &server.HandlerError{
			Message: "Woopsie, my bad\n",
			Status:  response.ServerError,
		}
	}

	_, err := w.Write([]byte("All good, frfr\n"))
	if err != nil {
		return &server.HandlerError{
			Message: err.Error(),
			Status:  response.ServerError,
		}
	}

	return nil
}
