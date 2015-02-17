package comm

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"code.google.com/p/go.crypto/ssh"
)

func Forward(localConn net.Conn, config *ssh.ClientConfig, addr string, remoteAddr string) {
	defer localConn.Close()

	// Setup sshClientConn (type *ssh.ClientConn)
	clientConn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		log.Fatalf("TCP Dial to %s failed: %s", addr, err)
	}

	// Setup sshConn (type net.Conn)
	conn, err := clientConn.Dial("tcp", remoteAddr)
	if err != nil {
		log.Fatalf("Tunneling failed for %s: %s", remoteAddr, err)
	}

	// Copy localConn.Reader to sshConn.Writer
	go func() {
		fmt.Print("Request: ")
		rdr := io.TeeReader(localConn, os.Stdout)
		_, err = io.Copy(conn, rdr)
		if err != nil {
			log.Fatalf("io.Copy failed: %v", err)
		}
		fmt.Println()
	}()

	// Copy sshConn.Reader to localConn.Writer
	go func() {
		rdr := io.TeeReader(conn, os.Stdout)
		_, err = io.Copy(localConn, rdr)
		if err != nil {
			log.Fatalf("Copy failed: %s", err)
		}
		io.Copy(conn, os.Stdout)
	}()
}
