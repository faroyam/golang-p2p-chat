package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("first argument ia username\nsecond argument is remote ip addr")
		return
	}
	username := os.Args[1]
	remoteAddr := os.Args[2]
	var wg sync.WaitGroup
	m := make(chan message)
	c := make(chan bool)

	fmt.Printf("Welcome, %s\nREMOTE: %s\n", username, remoteAddr)

	go server(&wg, c)
	go send(remoteAddr, m, c)
	go read(m, username)
	wg.Add(1)

	wg.Wait()
}

type message struct {
	Username string `json:"username"`
	Text     string `json:"text"`
}

func server(wg *sync.WaitGroup, c chan bool) {
	defer wg.Done()
	ln, err := net.Listen("tcp", ":49228")
	if err != nil {
		log.Println("Starting listener error!")
		return
	}
	defer ln.Close()
	for {
		fmt.Print("\n")
		log.Println("Listening...")
		conn, err := ln.Accept()
		if err != nil {
			continue
		}
		defer conn.Close()
		var msg = message{}
		for {
			err := json.NewDecoder(conn).Decode(&msg)

			if err != nil {
				c <- false
				log.Println("Disconnected!")
				break
			}
			if msg.Text != "" {
				log.Printf("%s: %s", msg.Username, msg.Text)
				fmt.Print(">>")
			}
		}
	}
}

func send(remoteAddr string, m chan message, c chan bool) {
CONNECTION:
	for {
		conn, err := net.Dial("tcp", remoteAddr+":49228")

		if err != nil {
			log.Println("Connecting...")
			time.Sleep(3 * time.Second)
			continue
		}
		log.Println("Connected!")
		fmt.Print(">>")

		for {
			select {
			case msg := <-m:
				err := json.NewEncoder(conn).Encode(msg)
				if err != nil {
					continue CONNECTION
				}

			case alive := <-c:
				if !alive {
					continue CONNECTION
				}
			}
		}
	}
}

func read(m chan message, username string) {
	var reader = bufio.NewReader(os.Stdin)
	var msg = message{Username: username}
	for {
		fmt.Print(">>")
		if text, _ := reader.ReadString('\n'); text != "" {
			if text != "\n" {
				msg.Text = strings.TrimSpace(text)
				m <- msg
			}
		}
	}
}
