package main

import (
	"fmt"
	"net/http"
)

const port = "8080"

func main() {
	// Creating a new serve mux
	// ServeMux is an HTTP request multiplexer that matches the URL of each incoming request against
	// a list of registered patterns and calls the handler for the pattern that most closely matches
	// the URL.
	mux := http.NewServeMux()

	// Creating a static file server
	// The http.FileServer creates a handler that serves HTTP requests with the contents of the file
	// system.
	// This will take index.html file from the root directory of the project.
	mux.Handle("/", http.FileServer(http.Dir(".")))

	// Creating a health check
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte("OK"))
	})

	// This is wrong way of implementing a static file server.
	// mux.HandleFunc("/app/assets", func(w http.ResponseWriter, r *http.Request) {
	// 	w.WriteHeader(http.StatusOK)
	// 	w.Write([]byte(`<pre><a href="logo.png">logo.png</a></pre>`))
	// })

	// mux.HandleFunc("/app", func(w http.ResponseWriter, r *http.Request) {
	// 	w.WriteHeader(http.StatusOK)
	// 	w.Write([]byte(`<html>Welcome to Chirpy</html>`))
	// })

	// Creating a static file server for /app
	// The http.StripPrefix strips the prefix from the URL and serves the request using the file
	// server.
	mux.Handle("/app/*", http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	// Creating a server
	// The Server type is an HTTP server. The pointer server points to a new Server value with the
	// specified network address and handler. 
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	// Starting the server using ListenAndServe method of the Server type
	err := server.ListenAndServe()
	fmt.Println("Serving on port " + port)
	if err != nil {
		panic(err)
	}
}
