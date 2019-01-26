package main

import (
	"log"
	"net/http"
)

// HomeTaskFunc handles "/" page
func HomeTaskFunc(w http.ResponseWriter, r *http.Request) {
	log.Println("Incoming request")
	w.Write([]byte("Hello"))
}

func main() {
	PORT := ":8080"
	log.Println("Running server on " + PORT)
	http.HandleFunc("/", HomeTaskFunc)
	log.Fatal(http.ListenAndServe(PORT, nil))
}
