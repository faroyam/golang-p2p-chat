package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

const (
	prompt = ">> "
)

func connListenerHandler(ln net.Listener) {
	defer ln.Close()
SERVER:
	for {
		servConn, err := ln.Accept()
		if err != nil {
			continue
		}
		defer servConn.Close()
		for {
			buf := make([]byte, 1024)
			data, err := servConn.Read(buf)
			message := string(buf[:data])

			switch {
			case err != nil:
				fmt.Println("Connection failed!")
				continue SERVER
			case message != "":
				fmt.Println(time.Now().Format("15:04:05 "), message)
				fmt.Print(prompt)
			}
		}
	}
}

func clientHandler(remoteAddr string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		clientConn, err := net.Dial("tcp", remoteAddr+":8081")
		if err != nil {
			fmt.Println("Connecting...")
			time.Sleep(3 * time.Second)
			continue
		}
		fmt.Println("Connected!")
		for {
			fmt.Print(prompt)
			text, _ := reader.ReadString('\n')
			message := []byte(strings.TrimSpace(text))
			if len(message) != 0 {
				_, err := clientConn.Write(message)
				if err != nil {
					break
				}
			}
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("first arguement is remote ip addr")
		return
	}
	remoteAddr := os.Args[1]
	fmt.Println("REMOTE: ", remoteAddr)
	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal(err)
	}
	go clientHandler(remoteAddr)
	connListenerHandler(ln)
}
