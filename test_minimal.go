package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// Test route
	mux.HandleFunc("/api/auth/login", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Route /api/auth/login called!")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"test": "ok"}`)
	})

	mux.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Route /api/login called!")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"test": "ok old"}`)
	})

	fmt.Println("Starting server on :8089")
	http.ListenAndServe(":8089", mux)
}
