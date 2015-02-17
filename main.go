package main

import (
	"log"
	"net"
	"os"
	"sshtunnel/comm"

	"code.google.com/p/go.crypto/ssh"
	"code.google.com/p/go.crypto/ssh/terminal"
)

func getPass() string {
	state, err := terminal.MakeRaw(0)
	if err != nil {
		log.Fatal(err)
	}
	defer terminal.Restore(0, state)
	term := terminal.NewTerminal(os.Stdout, ">")
	password, err := term.ReadPassword("Password: ")
	if err != nil {
		log.Fatal(err)
	}
	return password
}

func main() {
	port := os.Args[2]
	addr := os.Args[1] + ":22"
	pass := getPass()

	config := &ssh.ClientConfig{
		User: os.Getenv("LOGNAME"),
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
	}

	remoteAddr := "127.0.0.1:" + port

	local, err := net.Listen("tcp", remoteAddr)
	if err != nil {
		log.Fatalf("Listen failed on local address (%s): %s", remoteAddr, err)
	}
	defer local.Close()

	for {
		localConn, err := local.Accept()
		defer localConn.Close()
		if err != nil {
			log.Fatalf("listen.Accept failed: %v", err)
		}
		go comm.Forward(localConn, config, addr, remoteAddr)
	}
}
