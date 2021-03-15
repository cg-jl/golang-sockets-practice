package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
)

const (
	connHost = "localhost"
	connPort = "9999"
	connType = "tcp"
	printable  = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~ \t\n\r\x0b\x0c"
)

func generateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = printable[rand.Intn(len(printable))]
	}

	return string(b)
}


func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("Starting " + connType + " server on " + connHost + ":" + connPort)
	password := generateRandomString(6 + rand.Intn(34))
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

		go handleConnection(conn, password)
	}
}

func handleConnection(conn net.Conn, myPass string) {
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
