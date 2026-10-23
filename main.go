package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})

}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not found", 405)
	}
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hello" {
		http.Error(w, "404 not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	fmt.Fprintf(w, "Hello")
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	address := r.FormValue("Address")

	fmt.Fprintf(w, "Name: %s\nAddress: %s", name, address)
}

func main() {
	const port = "8080"

	mux := http.NewServeMux()

	apiCfg := apiConfig{fileserverHits: atomic.Int32{}}

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("./static")))))

	// mux.Handle("/app/assets/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not found", 405)
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})

	mux.HandleFunc("GET /metrics/", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /reset", apiCfg.handlerReset)

	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/form", formHandler)

	log.Printf("Serving on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
