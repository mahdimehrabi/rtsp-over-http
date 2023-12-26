package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

func main() {
	//dataConnectionEndpoint := "http://root:admin123456@172.24.90.42/axis-media/media.amp"
	//commandConnectionEndpoint := "http://root:admin123456@172.24.90.42/axis-media/media.amp"
	//rtspEndpoint := "rtsp://root:admin123456@172.24.90.42/axis-media/media.amp"

	go setup()
	dataConnectionEndpoint := "http://root:admin123456@localhost:8080/dump"
	commandConnectionEndpoint := "http://root:admin123456@localhost:8080/dump"
	rtspEndpoint := "rtsp://root:admin123456@localhost:8080/dump"

	sessionID := uuid.New().String()
	go establishDataConnection(dataConnectionEndpoint, sessionID)
	time.Sleep(1 * time.Second) //wait for establishing data connection

	//commandConnectionEndpoint = replaceCommandConnectionIP(respData, commandConnectionEndpoint)
	commandReq, _ := establishCommandConnection(commandConnectionEndpoint, sessionID)
	time.Sleep(1 * time.Second)
	fmt.Println("command connection established successfully")

	go describeCommand(commandReq, commandConnectionEndpoint, rtspEndpoint, 1, sessionID)
	fmt.Println("describe command has been sent")

	time.Sleep(100 * time.Second)
}

func establishDataConnection(url string, sessionID string) (*http.Response, *http.Client) {
	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("x-sessioncookie", sessionID)
	req.Header.Add("Cache-Control", "no-store")
	req.Header.Add("Pragma", "no-cache")
	req.Header.Add("Accept", " application/x-rtsp-tunnelled")
	req.Header.Add("Content-Length", "32767")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	receiveData(resp)

	return resp, &client
}

func receiveData(resp *http.Response) {
	for {
		bf := bytes.NewBuffer(make([]byte, 0))
		_, err := io.Copy(bf, resp.Body)
		if err != nil {
			if err == io.EOF {
				fmt.Println("data connection closed (EOF)")
				return
			}
			log.Println("error reading data:", err)
			return
		}
		fmt.Printf("received data from data connection: %s\n", bf.String())
		time.Sleep(100 * time.Millisecond)
	}
}

func establishCommandConnection(url string, sessionID string) (*http.Request, *http.Client) {
	client := http.Client{}
	reqBody := bytes.NewBuffer(make([]byte, 0))
	req, err := http.NewRequest(http.MethodPost, url, reqBody)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("x-sessioncookie", sessionID)
	req.Header.Add("Content-Length", "32767")
	req.Header.Add("Pragma", "no-cache")
	req.Header.Add("Accept", "application/x-rtsp-tunnelled")
	go client.Do(req)
	time.Sleep(3 * time.Second)
	return req, &client
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

func describeCommand(req *http.Request, commandEndpoint string, rtspEndpoint string, cseqValue int, sessionID string) {
	reqBody := bytes.NewBuffer(make([]byte, 0))
	rtspCommand := fmt.Sprintf("DESCRIBE %s RTSP/1.0\r\nCSeq: %d\r\nUser-Agent: Axis AMC\r\nAccept: application/sdp\r\n\r\n", rtspEndpoint, cseqValue)
	encoder := base64.NewEncoder(base64.StdEncoding, reqBody)
	_, err := encoder.Write([]byte(rtspCommand))
	if err != nil {
		log.Fatal(err)
	}
	err = req.Write(reqBody)
	if err != nil {
		log.Fatal(err)
	}
}
func setupCommand(reqBody *bytes.Buffer, rtspEndpoint string, cseqValue int) {
	encoder := base64.NewEncoder(base64.StdEncoding, reqBody)

	setupCommand := fmt.Sprintf("SETUP %s RTSP/1.0\r\nCSeq: %d\r\nTransport: <transport-specifier>\r\n\r\n", rtspEndpoint, cseqValue)
	_, err := encoder.Write([]byte(setupCommand))
	if err != nil {
		log.Fatal(err)
	}
}

func playCommand(reqBody *bytes.Buffer, rtspEndpoint string, cseqValue int) {
	encoder := base64.NewEncoder(base64.StdEncoding, reqBody)

	playCommand := fmt.Sprintf("PLAY %s RTSP/1.0\r\nCSeq: %d\r\n\r\n", rtspEndpoint, cseqValue)
	_, err := encoder.Write([]byte(playCommand))
	if err != nil {
		log.Fatal(err)
	}
}

func teardownCommand(reqBody *bytes.Buffer, rtspEndpoint string, cseqValue int) {
	encoder := base64.NewEncoder(base64.StdEncoding, reqBody)

	teardownCommand := fmt.Sprintf("TEARDOWN %s RTSP/1.0\r\nCSeq: %d\r\n\r\n", rtspEndpoint, cseqValue)
	_, err := encoder.Write([]byte(teardownCommand))
	if err != nil {
		log.Fatal(err)
	}
}
