package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("📡 Serving on http://localhost:8080")
	http.Handle("/", http.FileServer(http.Dir("dist")))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
