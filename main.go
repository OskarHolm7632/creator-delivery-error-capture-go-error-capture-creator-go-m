package main

import (
	"log"
	"net/http"
)

func main() {
	client, err := NewInfraiClient()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/delivery", deliveryHandler(client))
	log.Println("creator delivery service listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
