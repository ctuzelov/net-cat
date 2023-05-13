package tools

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
)

func CheckArgs(args []string) string {
	Port := ""
	if len(args) == 2 {
		Port = os.Args[1]
	} else if len(args) == 1 {
		Port = "8989"
	} else {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return ""
	}
	return Port
}

func Welcome(conn net.Conn) {
	content, _ := ioutil.ReadFile("tools/welcome.txt")
	fmt.Fprint(conn, string(content))
}
