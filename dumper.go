package main

import (
	"fmt"
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

	// Respond with a simple message
	w.WriteHeader(http.StatusOK)
	for {
		time.Sleep(1 * time.Second)
		fmt.Fprint(w, "Request data dumped successfully!")
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
