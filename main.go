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
	dataConnectionEndpoint := "http://root:admin123456@172.24.90.99/axis-media/media.amp?videocodec=h264&audio=0 "
	commandConnectionEndpoint := "http://root:admin123456@172.24.90.99/axis-media/media.amp?videocodec=h264&audio=0 "
	rtspEndpoint := "rtsp://root:admin123456@172.24.90.99/axis-media/media.amp?videocodec=h264&audio=0 "
	sessionID := uuid.New().String()
	respData, _ := establishDataConnection(dataConnectionEndpoint, sessionID)
	if respData.StatusCode != http.StatusOK {
		log.Fatalf("failed to establish data connection status code %d", respData.StatusCode)
	}
	commandConnectionEndpoint = replaceCommandConnectionIP(respData, commandConnectionEndpoint)
	fmt.Println("data connection established status=%d", respData.StatusCode)
	respCommand, reqBody, _ := establishCommandConnection(commandConnectionEndpoint, sessionID)
	if respCommand.StatusCode != http.StatusOK {
		b := bytes.NewBuffer(make([]byte, 0))
		respCommand.Body.Close()
		io.Copy(b, respCommand.Body)
		log.Fatalf("failed to establish data connection status code %d,data:%s", respCommand.StatusCode, b.String())
	}
	fmt.Println("data connection established status=%d", respCommand.StatusCode)

	describeCommand(reqBody, rtspEndpoint)
	fmt.Println("describe command has been sent")

	go receiveData(respData.Body)
	receiveData(respCommand.Body)
}

func establishDataConnection(url string, sessionID string) (*http.Response, *http.Client) {
	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		panic(err)
	}
	req.Header.Add("x-sessioncookie", sessionID)
	req.Header.Add("Cache-Control", "no-store")
	req.Header.Add("Pragma", "no-cache")
	req.Header.Add("Accept", " application/x-rtsp-tunnelled")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	return resp, &client
}

func receiveData(reader io.Reader) {
	for {
		data := bytes.NewBuffer(make([]byte, 0))
		fmt.Println("received data from data connection", data)
		time.Sleep(100 * time.Millisecond)
	}
}

func establishCommandConnection(url string, sessionID string) (*http.Response, *bytes.Buffer, *http.Client) {
	client := http.Client{}
	reqBody := bytes.NewBuffer(make([]byte, 0))
	req, err := http.NewRequest(http.MethodGet, url, reqBody)
	if err != nil {
		panic(err)
	}
	req.Header.Add("x-sessioncookie", sessionID)
	req.Header.Add("Content-Length", "32767")
	req.Header.Add("Pragma", "no-cache")
	req.Header.Add("Accept", " application/x-rtsp-tunnelled")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	return resp, reqBody, &client
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

func describeCommand(reqBody *bytes.Buffer, rtspEndpoint string) {
	encoder := base64.NewEncoder(base64.StdEncoding, reqBody)
	_, err := encoder.Write([]byte(fmt.Sprintf("DESCRIBE %s RTSP/1.0\nCSeq: 1\nUser-Agent: Axis AMC\nAccept: application/sdp", rtspEndpoint)))
	if err != nil {
		log.Fatal(err)
	}
}
