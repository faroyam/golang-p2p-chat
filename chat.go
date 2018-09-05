package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	certSender, err := tls.LoadX509KeyPair("certs/client.pem", "certs/client.key")
	if err != nil {
		log.Fatalf("client: loadkeys: %s", err)
	}
	configSender := tls.Config{Certificates: []tls.Certificate{certSender}, InsecureSkipVerify: true}

	certListener, err := tls.LoadX509KeyPair("certs/server.pem", "certs/server.key")
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
	}
	configListener := tls.Config{Certificates: []tls.Certificate{certListener}}

	if len(os.Args) < 3 {
		fmt.Println("first argument is your username\nsecond argument is remote ip addr")
		return
	}
	username := os.Args[1]
	remoteAddr := os.Args[2]
	var wg sync.WaitGroup
	m := make(chan message)
	c := make(chan bool)

	fmt.Printf("Welcome, %s\nREMOTE: %s\n", username, remoteAddr)

	go listener(&configListener, &wg, c)
	go sender(remoteAddr, &configSender, m, c)
	go read(m, username)
	wg.Add(3)

	wg.Wait()
}

type message struct {
	Username string `json:"username"`
	Text     string `json:"text"`
}

func listener(config *tls.Config, wg *sync.WaitGroup, c chan bool) {
	defer wg.Done()
	ln, err := tls.Listen("tcp", ":49228", config)
	if err != nil {
		log.Println("Starting listener error!", err)
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

func sender(remoteAddr string, config *tls.Config, m chan message, c chan bool) {
CONNECTION:
	for {
		conn, err := tls.Dial("tcp", remoteAddr+":49228", config)
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
		if text, _ := reader.ReadString('\n'); text != "\n" {

			msg.Text = strings.TrimSpace(text)
			m <- msg
		}
	}
}
