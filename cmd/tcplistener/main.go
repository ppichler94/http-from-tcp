package main

import (
	"fmt"
	"net"

	"http.ppichler94.io/internal/request"
)

func main() {
	l, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Println("Error listening on port 42069:", err)
		return
	}

	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		req, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Println("Error reading request:", err)
			return
		}

		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n-Target: %v\n-Version: %v\n", req.RequestLine.Method, req.RequestLine.RequestTarget, req.RequestLine.HttpVersion)

		fmt.Println("Headers:")
		for key, value := range req.Headers {
			fmt.Printf("- %s: %s\n", key, value)
		}

		fmt.Println("Body:")
		fmt.Println(string(req.Body))
	}

}
