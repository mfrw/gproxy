// Author: Muhammad Falak <falakreyaz@gmail.com>
// Date:   09-March-2021

// This package acts as a simple TCP Reverse proxy.
// The main usage of this is to relay SSH connections.
package main

import (
	"flag"
	"io"
	"log"
	"net"
)

var (
	localAddr  = flag.String("l", ":8080", "host:port to listen on")
	remoteAddr = flag.String("r", "172.28.188.5:22", "host:port to forward to")
	prefix     = flag.String("p", "proxy: ", "String to prefix log output")
)

func forward(conn net.Conn) {
	cl, err := net.Dial("tcp", *remoteAddr)
	if err != nil {
		defer conn.Close()
		log.Printf("Dial failed: %v", err)
		return
	}
	defer cl.Close()
	log.Printf("Forwarding from %v to %v\n", conn.LocalAddr(), cl.RemoteAddr())
	go io.Copy(cl, conn)
	io.Copy(conn, cl)
}

func main() {
	flag.Parse()
	log.SetPrefix(*prefix + ": ")

	listener, err := net.Listen("tcp", *localAddr)
	if err != nil {
		log.Fatalf("Failed to setup listener: %v", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("ERROR: failed to accept listener: %v", err)
		}
		log.Printf("Accepted connection from %v\n", conn.RemoteAddr().String())
		go forward(conn)
	}
}
