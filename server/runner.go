package server

import (
	"fmt"
	"log"
	"net"
	"sync"
)

func StartServer(port string) {
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", port))
	if err != nil {
		fmt.Println(err)
		log.Fatal("ERROR on starting server")
	}

	defer listener.Close()

	log.Printf("Listening on the port: %s", port)
	chat := Chat{[]Usr{}, make(chan Message), []Message{}, sync.Mutex{}}

	go chat.Massenger()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection from client: %s", err)
		} else {
			go chat.ClientServe(conn)
		}
	}
}
