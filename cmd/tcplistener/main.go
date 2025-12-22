package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string)

	go func() {
		defer close(out)
		defer f.Close()

		buf := make([]byte, 8)
		line := ""
		for {
			n, err := f.Read(buf)
			if err != nil {
				break
			}

			parts := strings.Split(string(buf[:n]), "\n")
			line += parts[0]

			for i := 1; i < len(parts); i++ {
				out <- line
				line = parts[i]
			}
		}

		if len(line) > 0 {
			out <- line
		}

	}()

	return out
}

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
		lines := getLinesChannel(conn)
		for line := range lines {
			fmt.Printf("read: %s\n", line)
		}
	}
}
