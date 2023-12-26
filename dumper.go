package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func dumpHandler(w http.ResponseWriter, r *http.Request) {
	// Dump request method and URL
	fmt.Printf("Request: %s %s\n", r.Method, r.URL.Path)

	// Dump request headers
	fmt.Println("Headers:")
	for key, values := range r.Header {
		for _, value := range values {
			fmt.Printf("  %s: %s\n", key, value)
		}
	}

	// Dump request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v\n", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Println("Body:")
	fmt.Printf("%s\n", body)
	w.Header().Add("Content-Length", r.Header.Get("Content-Length"))
	// Respond with a simple message if its GET
	if r.Method == http.MethodGet {
		w.WriteHeader(http.StatusOK)

		for i := 0; i < 5; i++ {
			time.Sleep(1 * time.Second)
			_, err := w.Write([]byte(" a data connection simple response "))
			if err != nil {
				log.Fatal(err)
			}
		}
	} else {

		//request is post
		for {
			buf := bytes.NewBuffer(make([]byte, 0))
			_, err := io.Copy(buf, r.Body)
			if err != nil {
				if err == io.EOF {
					fmt.Println("eof request post")
					return
				}
				log.Println("error reading request post:", err)
				return
			}
			fmt.Printf("received data from data request post: %s\n", buf.String())
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func setup() {
	// Register the dumpHandler function to handle requests to the "/dump" path
	http.HandleFunc("/dump", dumpHandler)

	// Start the HTTP server on port 8080
	port := 8080
	fmt.Printf("Server listening on :%d...\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
