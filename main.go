package main

import (
    "net/http"
)

func main() {
    mux := http.NewServeMux()

    // Step 1: readiness endpoint
    mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/plain; charset=utf-8")
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("OK"))
    })

    // Step 2: fileserver at /app/
    fileServer := http.FileServer(http.Dir("."))
    mux.Handle("/app/", http.StripPrefix("/app/", fileServer))

    http.ListenAndServe(":8080", mux)
}
