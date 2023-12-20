package main

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"net/url"
	"time"
)

func main() {
	dataConnectionEndpoint := "http://172.24.90.99/axis-media/media.amp"
	commandConnectionEndpoint := "http://172.24.90.99/axis-media/media.amp"
	sessionID := uuid.New().String()
	resp, _ := establishDataConnection(dataConnectionEndpoint, sessionID)

	replaceCommandConnectionIP(resp, commandConnectionEndpoint)
	fmt.Println("data connection established status=%d", resp.Status)
	receiveData(resp.Body)
}

func establishDataConnection(url string, sessionID string) (*http.Response, *http.Client) {
	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		panic(err)
	}
	req.Header.Add("x-sessioncookie", sessionID)
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

func establishCommandConnection(url string, sessionID string) (*http.Response, *http.Client) {
	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		panic(err)
	}
	req.Header.Add("x-sessioncookie", sessionID)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	return resp, &client
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
