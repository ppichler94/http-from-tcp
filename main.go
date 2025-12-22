package main

import (
	"fmt"
	"io"
	"os"
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
	f, err := os.Open("messages.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	lines := getLinesChannel(f)
	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}
}
