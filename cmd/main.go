package main

import (
	"aTES/infrastructure"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Loading the configuration.
	config, err := infrastructure.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Connecting to the database.
	db, err := infrastructure.InitDB(config)
	if err != nil {
		log.Fatalf("Error inititalising the database: %v", err)
	}
	defer db.Close()

	// Initialising HTTP handlers.
	httpHandlers := infrastructure.NewHTTPHandlers(db)

	// Setting up routs.
	http.HandleFunc("/tasks", httpHandlers.TaskHandler)
	http.HandleFunc("/accounting", httpHandlers.AccountingHandler)

	// Starting the HTTP server.
	addr := fmt.Sprintf(":%d", config.Port)
	fmt.Printf("Starting server on %s...\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Error starting the server: %v", err)
	}
}
