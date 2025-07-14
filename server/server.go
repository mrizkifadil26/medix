package server

import (
	"log"
	"net/http"
)

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s", r.Method, r.URL.Path)
		handler.ServeHTTP(w, r)
	})
}

func Serve(distPath string, port string) error {
	log.Printf("Serving %s on http://localhost:%s", distPath, port)

	fs := http.FileServer(http.Dir(distPath))
	return http.ListenAndServe(":"+port, logRequest(fs))
}
