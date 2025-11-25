package main

import (
    "net/http"
)

func main() {
    // Serve current directory on /
    http.Handle("/", http.FileServer(http.Dir(".")))

    // Start server
    http.ListenAndServe(":8080", nil)
}
