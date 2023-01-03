package main

import (
	"fmt"
	"net/http"

	"github.com/bradleyshawkins/go-clean-architecture/refactored-carsonoid/handlers"
)

func main() {
	mux := handlers.NewMux()

	fmt.Println("Starting router...")
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Printf("Error received: %v\n", err)
	}
}
