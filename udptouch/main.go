package main

import (
	"fmt"
	"os"

	"github.com/chaosannals/trial-go-udptouch/util"
)

func main() {
	argc := len(os.Args)
	if argc > 3 {
		id := os.Args[1]
		target := os.Args[2]
		server := os.Args[3]
		fmt.Printf("id: %s touch: %s by %s\n", id, target, server)
		client, err := util.NewTouchClient(server, id, target)
		if err != nil {
			fmt.Printf("client new err: %v", err)
		}
		if err := client.Start(); err != nil {
			fmt.Printf("client err: %v\n", err)
		}
	} else if argc > 1 {
		fmt.Println("args: [id] [target] [server]")
	} else {
		port := 38383
		fmt.Printf("server start, port: %d", port)
		server := util.NewTouchServer(port)
		if err := server.Serve(); err != nil {
			fmt.Printf("server err: %v\n", err)
		}
	}
}
