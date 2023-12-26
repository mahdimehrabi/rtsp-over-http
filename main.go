package main

import (
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

func main() {
	//dataConnectionEndpoint := "http://root:admin123456@172.24.90.42:80/axis-media/media.amp"
	//commandConnectionEndpoint := "http://root:admin123456@172.24.90.42:80/axis-media/media.amp"
	//rtspEndpoint := "rtsp://root:admin123456@172.24.90.42:80/axis-media/media.amp"

	go setup()
	dataConnectionEndpoint := "http://root:admin123456@localhost:8080/dump"
	commandConnectionEndpoint := "http://root:admin123456@localhost:8080/dump"
	rtspEndpoint := "rtsp://root:admin123456@localhost:8080/dump"

	sessionID := uuid.New().String()
	establishDataConnection(dataConnectionEndpoint, sessionID)
	time.Sleep(1 * time.Second) //wait for establishing data connection

	//commandConnectionEndpoint = replaceCommandConnectionIP(respData, commandConnectionEndpoint)
	conn := establishCommandConnection(commandConnectionEndpoint, sessionID)
	time.Sleep(1 * time.Second)
	fmt.Println("command connection established successfully")

	go describeCommand(conn, commandConnectionEndpoint, rtspEndpoint, 1)
	fmt.Println("describe command has been sent")

	time.Sleep(100 * time.Second)
}

func establishDataConnection(serverURL string, sessionID string) {
	// Parse the server URL
	parsedURL, err := url.Parse(serverURL)
	if err != nil {
		fmt.Println("Error parsing server URL:", err)
		return
	}

	// Create a TCP connection to the server
	conn, err := net.Dial("tcp", parsedURL.Host)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}

	request := fmt.Sprintf("GET %s HTTP/1.0\r\n"+
		"Host: %s\r\n"+
		"x-sessioncookie: %s\r\n"+
		"Cache-Control: no-store\r\n"+
		"Pragma: no-cache\r\n"+
		"Accept: application/x-rtsp-tunnelled\r\n"+
		"Content-Length: 32767\r\n"+
		"\r\n", parsedURL.Path,
		parsedURL.Hostname(), sessionID)

	_, err = conn.Write([]byte(request))
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	go receiveData(conn)
}

func receiveData(resp net.Conn) {
	for {
		bt := make([]byte, 1024)
		n, err := resp.Read(bt)
		if err != nil {
			if err == io.EOF {
				fmt.Println("data connection closed (EOF)")
				return
			}
			log.Println("error reading data:", err)
			return
		}
		fmt.Printf("received data from data connection: %s\n", bt[:n])
		time.Sleep(100 * time.Millisecond)
	}
}

func establishCommandConnection(serverURL string, sessionID string) net.Conn {
	// Parse the server URL
	parsedURL, err := url.Parse(serverURL)
	if err != nil {
		fmt.Println("Error parsing server URL:", err)
		return nil
	}

	// Create a TCP connection to the server
	conn, err := net.Dial("tcp", parsedURL.Host)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return nil
	}

	request := fmt.Sprintf("POST %s HTTP/1.0\r\n"+
		"Host: %s\r\n"+
		"x-sessioncookie: %s\r\n"+
		"Cache-Control: no-store\r\n"+
		"Pragma: no-cache\r\n"+
		"Accept: application/x-rtsp-tunnelled\r\n"+
		"Content-Length: 32767\r\n"+
		"\r\n", parsedURL.Path,
		parsedURL.Hostname(), sessionID)

	_, err = conn.Write([]byte(request))
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil
	}

	return conn
}

func replaceCommandConnectionIP(resp *http.Response, commandURL string) string {
	connectionIP := resp.Header.Get("x-server-ip-address")
	if connectionIP == "" {
		return commandURL
	}
	u, err := url.Parse(commandURL)
	if err != nil {
		fmt.Println(err)
		return commandURL
	}
	u.Host = connectionIP
	return u.String()
}

func describeCommand(req net.Conn, commandEndpoint string, rtspEndpoint string, cseqValue int) {
	rtspCommand := fmt.Sprintf("DESCRIBE %s RTSP/1.0\r\nCSeq: %d\r\nUser-Agent: Axis AMC\r\nAccept: application/sdp\r\n\r\n", rtspEndpoint, cseqValue)
	encoder := base64.NewEncoder(base64.StdEncoding, req)
	_, err := encoder.Write([]byte(rtspCommand))
	if err != nil {
		log.Fatal(err)
	}
}
