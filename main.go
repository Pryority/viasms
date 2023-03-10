package main

import (
	"fmt"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}

	defer ln.Close()

	fmt.Println("Server listening on port 8080")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting:", err.Error())
			continue
		}

		go handleConnection(conn)
	}
}

type Client struct {
	conn net.Conn
}

var clients []Client

func handleConnection(conn net.Conn) {
	defer conn.Close()

	client := Client{conn}
	clients = append(clients, client)

	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			break
		}

		msg := string(buf[:n])

		// Broadcast message to all clients
		for _, c := range clients {
			if c.conn != conn {
				_, err := c.conn.Write([]byte(msg))
				if err != nil {
					fmt.Println("Error writing:", err.Error())
					continue
				}
			}
		}
	}

	// Remove client from list of clients
	for i, c := range clients {
		if c.conn == conn {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}
}
