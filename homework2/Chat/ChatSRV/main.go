package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type client struct {
	ch       chan<- string
	nickname string
}

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string)
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()
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
				cli.ch <- msg
			}

		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli.ch)
		}
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string)

	go clientWriter(conn, ch)
	buffer := make([]byte, 100)
	_, err := conn.Read(buffer)
	if err != nil {
		fmt.Print(err)
		return
	}
	who := string(buffer)

	ch <- "your nickname is set to " + who
	messages <- who + " has arrived"
	entering <- client{ch, who}

	inputMsg := bufio.NewScanner(conn)
	for inputMsg.Scan() {
		messages <- who + ": " + inputMsg.Text()
	}
	leaving <- client{ch, who}
	messages <- who + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}
