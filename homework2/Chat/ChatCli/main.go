package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	fmt.Print("Enter your nickname:")
	reader := bufio.NewReader(os.Stdin)
	nickname, _, err := reader.ReadLine()
	if err != nil {
		fmt.Println("cannot read data, program exit")
		return
	}

	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	_, err = conn.Write(nickname)
	if err != nil {
		fmt.Print(err)
		return
	}
	go func() {
		io.Copy(os.Stdout, conn)
	}()

	io.Copy(conn, os.Stdin)
	fmt.Printf("%s: exit", conn.LocalAddr())
}