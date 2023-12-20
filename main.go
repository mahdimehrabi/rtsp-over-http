package main

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"time"
)

func main() {
	sessionID := uuid.New().String()
	dataReader, _ := establishDataConnection("http://172.24.90.99/axis-media/media.amp", sessionID)
	fmt.Println("data connection established status=%d", dataReader.Status)
	receiveData(dataReader.Body)
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
