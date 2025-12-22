package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return
	}

	con, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("Error dialing UDP:", err)
		return
	}
	defer con.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from stdin:", err)
			continue
		}
		_, err = con.Write([]byte(line))
		if err != nil {
			fmt.Println("Error sending UDP message:", err)
		}
	}
}
