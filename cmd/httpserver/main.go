package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"http.ppichler94.io/internal/request"
	"http.ppichler94.io/internal/response"
	"http.ppichler94.io/internal/server"
)

const port = 8080

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

func handler(w *response.Writer, req *request.Request) {
	var message string
	var status response.StatusCode
	if req.RequestLine.RequestTarget == "/yourproblem" {
		message = "<html>\n  <head>\n    <title>400 Bad Request</title>\n  </head>\n  <body>\n    <h1>Bad Request</h1>\n    <p>Your request honestly kinda sucked.</p>\n  </body>\n</html>"
		status = response.BadRequest
	} else if req.RequestLine.RequestTarget == "/myproblem" {
		message = "<html>\n  <head>\n    <title>500 Internal Server Error</title>\n  </head>\n  <body>\n    <h1>Internal Server Error</h1>\n    <p>Okay, you know what? This one is on me.</p>\n  </body>\n</html>"
		status = response.ServerError
	} else {
		message = "<html>\n  <head>\n    <title>200 OK</title>\n  </head>\n  <body>\n    <h1>Success!</h1>\n    <p>Your request was an absolute banger.</p>\n  </body>\n</html>"
		status = response.OK
	}

	w.WriteStatusLine(status)
	h := response.GetDefaultHeaders(len(message))
	h.Set("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody([]byte(message))
}
