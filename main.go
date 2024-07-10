package main

import (
	"fmt"
	"net/http"
)

const port = "8080"

func main() {
	// Creating a new serve mux
	mux := http.NewServeMux()

	// Creating a file server
	mux.Handle("/", http.FileServer(http.Dir("./")))

	// mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	// Creating a server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Println("Serving on port " + port)
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
