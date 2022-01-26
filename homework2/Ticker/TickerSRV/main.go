package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

type client chan<- string

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string)
)

func main() {
	l, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	go msgFunc()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func broadcaster() {
	clients := make(map[client]bool)
	for {
		select {
		case msg := <-messages:
			for cli := range clients {
				cli <- msg
			}

		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}


func msgFunc() {
	userChan := make(chan string)
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("Enter text: ")
			text, _ := reader.ReadString('\n')
			text = text[:len(text)-1]
			userChan <- text
		}
	}()

	ticker := time.NewTicker(1 * time.Second)
	var message string
	for {
		<-ticker.C
		select {
		case userMsg := <-userChan:
			message = time.Now().Format("15:04:05") + " " + userMsg
		default: 
			message = time.Now().Format("15:04:05")
		}
		messages <- message
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string)

	entering <- ch

	for msg := range ch {
		_, err := fmt.Fprintln(conn, msg)
		if err != nil {
			break
		}
	}

	leaving <- ch
	conn.Close()
}
