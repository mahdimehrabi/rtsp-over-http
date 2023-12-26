package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

func handleRequest(conn net.Conn) {
	defer conn.Close()

	// Create a buffered reader to read from the connection
	reader := bufio.NewReader(conn)

	// Read the first line of the request (request line)
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading request line:", err)
		return
	}

	// Split the request line into method, path, and protocol
	parts := strings.Fields(requestLine)
	if len(parts) != 3 {
		fmt.Println("Invalid request line:", requestLine)
		return
	}

	method := parts[0]
	path := parts[1]
	protocol := parts[2]

	// Print the request information
	fmt.Printf("Method: %s\nPath: %s\nProtocol: %s\n", method, path, protocol)

	// Respond with a simple message
	response := "HTTP/1.0 200 OK\r\n" +
		"Content-Type: text/plain\r\n" +
		"\r\n" +
		"Hello, this is a manual HTTP 1.0 server!\n"

	if method == http.MethodGet {
		for {
			// Write the response to the connection
			_, err = conn.Write([]byte(response))
			if err != nil {
				fmt.Println("Error writing response:", err)
				return
			}
			time.Sleep(1 * time.Second)
		}
	} else {
		for {
			bt := make([]byte, 1024)
			// Read the first line of the request (request line)
			_, err := conn.Read(bt)
			if err != nil {
				fmt.Println("Error reading request line:", err)
				return
			}
			fmt.Println("readed from post", string(bt))
		}
	}
}

func setup() {
	// Listen on port 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on port 8080...")

	// Accept incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Handle each incoming connection in a new goroutine
		go handleRequest(conn)
	}
}
