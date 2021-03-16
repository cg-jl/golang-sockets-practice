package main

import (
	"bufio"
	"fmt"
	"context"
	"math/rand"
	"net"
	"syscall"
	"os"
	"strings"
	"time"
)

const (
	connHost = "localhost"
	connPort = "9999"
	connType = "tcp"
	printable  = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~ \t"
)

func reusePort(network, address string, conn syscall.RawConn) error {
	return conn.Control(func(descriptor uintptr) {
		syscall.SetsockoptInt(int(descriptor), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	})
}

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

	fmt.Printf("Generated password: \"%s\"\n", password)


	config := &net.ListenConfig{Control: reusePort}



	l, err := config.Listen(context.Background(), "tcp", "localhost:9999")
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

		if len(passwd) < len(myPass) {
			passwd += strings.Repeat("\x00", len(myPass)-len(passwd))
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
