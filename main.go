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

	apiCfg := apiConfig{}

	mux := http.NewServeMux()

	// --- FILE SERVER ---
	fsHandler := apiCfg.middlewareMetricsInc(
		http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))),
	)
	mux.Handle("/app/", fsHandler)
	mux.Handle("/app", fsHandler)

	// --- API ROUTES ---
	mux.HandleFunc("/api/healthz", exactPathAndMethod("/api/healthz", http.MethodGet, handlerReadiness))

	// --- ADMIN ROUTES ---
	mux.HandleFunc("/admin/metrics", apiCfg.handlerAdminMetrics)
	mux.HandleFunc("/admin/reset", apiCfg.handlerAdminReset)

	// --- JSON parsers ---
	mux.HandleFunc("/api/validate_chirp", apiCfg.handlerValidateChirp)

	// --- START SERVER ---
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}



	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}

// --------------------
// UTILITIES
// --------------------

// Enforce exact path + correct method
func exactPathAndMethod(path, method string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != path {
			http.NotFound(w, r)
			return
		}
		if r.Method != method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Method Not Allowed"))
			return
		}
		handler(w, r)
	}
}
