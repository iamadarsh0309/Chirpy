package main

import (
	"fmt"
	"log"
	"net/http"
	// "net/http"
)

func main() {
	fmt.Println("Hello, World!")
	const port = "8080"

	// fileServer := http.FileServer(http.Dir("./static"))
	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	mux.Handle("/", http.FileServer(http.Dir("./static")))

	log.Println("Serving on port : %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
