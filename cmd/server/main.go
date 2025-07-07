package main

import (
	"log"
	"net/http"
)

func main() {
	const port = ":8080"
	const dir = "dist"

	log.Printf("📡 Serving static site at http://localhost%s\n", port)
	fs := http.FileServer(http.Dir(dir))
	http.Handle("/", fs)

	log.Fatal(http.ListenAndServe(port, nil))
}
