package main

import (
	"fmt"
	"net/http"
)

func server1() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Server 1")
	})

	http.ListenAndServe(":8080", mux)
}

func server2() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Server 2")
	})

	http.ListenAndServe(":9090", mux)
}

func main() {
	go server1()   // runs in parallel
	server2()      // main thread
}