package main

import "net/http"

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("OK"))
}
