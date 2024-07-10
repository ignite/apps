package main

import (
	"log"
	"os"

	"github.com/ignite/apps/spaceship/pkg/ssh"
)

func main() {
	if len(os.Args) < 6 {
		log.Fatalf("Usage: %s <user> <password> <host> <port> <app_path>\n", os.Args[0])
	}

	user := os.Args[1]
	password := os.Args[2]
	host := os.Args[3]
	port := os.Args[4]
	appPath := os.Args[5]

	ssh.New(user, password, host, port, appPath)
}
