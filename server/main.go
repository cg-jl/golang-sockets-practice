package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	connHost = "localhost"
	connPort = "9999"
	connType = "tcp"
	myPass   = "0(\\6`\tBfFk"
)

func main() {
	fmt.Println("Starting " + connType + " server on " + connHost + ":" + connPort)
	l, err := net.Listen(connType, connHost+":"+connPort)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()

	for {
		// accept next client
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error connecting:", err.Error())
			return
		}
		fmt.Println("Client connected.")

		fmt.Println("Client " + conn.RemoteAddr().String() + " connected.")

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	bw := bufio.NewWriter(conn)
	br := bufio.NewReader(conn)
	for {
		bw.WriteString("Password: ")
		bw.Flush()
		buffer, err := br.ReadBytes('\n')
		if err != nil {
			fmt.Println("Cient left.")
			conn.Close()
			return
		}

		buffer = buffer[:len(buffer)-1]
		passwd := string(buffer)

		if len(passwd) < 10 {
			passwd += strings.Repeat("\x00", 10-len(passwd))
		}

		if passwd == myPass {
			bw.WriteString("Yes!\n")
			bw.Flush()
			conn.Close()
			return
		}
		for i, chr := range passwd {
			if i >= len(myPass) {
				break
			}
			if chr == rune(myPass[i]) {
				bw.WriteRune('1')
			} else {
				bw.WriteRune('0')
			}
		}
		bw.WriteRune('\n')
		bw.Flush()

	}
}
