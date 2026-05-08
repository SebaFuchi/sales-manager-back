package main

import (
	"log"
	"os"
	"template/internal/server"
)

// @title Template
// @version 1.0.0
// @BasePath /api

func main() {
	// When using local development uncomment this line of code with your own port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	serv, err := server.New(port)
	if err != nil {
		log.Fatal(err)
	}

	// Ensure all default logger output goes to stdout with timestamp and short file info
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
	log.SetPrefix("Template Api: ")

	serv.Start()
}
