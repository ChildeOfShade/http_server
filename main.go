package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	mux := http.NewServeMux()

	// file server ONLY under /app/
	fsHandler := apiCfg.middlewareMetricsInc(
	    http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))),
	)
	mux.Handle("/app/", fsHandler)
	mux.Handle("/app", fsHandler)


	// API routes
	mux.HandleFunc("/healthz", handlerReadiness)
	mux.HandleFunc("/reset", apiCfg.handlerReset)
	mux.HandleFunc("/metrics", apiCfg.handlerMetrics)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
