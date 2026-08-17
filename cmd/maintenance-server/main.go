package main

import (
	"log"
	"net/http"
	"os"

	"github.com/1260124186-cc/maintenance-work-order-service/internal/store"
	"github.com/1260124186-cc/maintenance-work-order-service/internal/transport"
)

func main() {
	address := os.Getenv("MAINTENANCE_ADDR")
	if address == "" {
		address = ":8080"
	}

	handler := transport.NewServer(store.NewMemoryRepository())
	log.Printf("maintenance work-order service listening on %s", address)
	log.Fatal(http.ListenAndServe(address, handler))
}
