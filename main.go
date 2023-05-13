package main

import (
	"fmt"
	"os"

	"net-cat/server"
	"net-cat/tools"
)

func main() {
	port := tools.CheckArgs(os.Args)
	if len(port) == 0 {
		fmt.Println("Try another port")
		return
	}

	server.StartServer(port)
}
