package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/germangorelkin/go-inspector/inspector"
)

func main() {
	// Define flags (none required - server assigns defaults)
	flag.Parse()

	// Get credentials from environment
	apiKey := os.Getenv("API_KEY")
	instance := os.Getenv("INSTANCE")
	if apiKey == "" || instance == "" {
		log.Fatal("API_KEY and INSTANCE environment variables must be set")
	}

	// Create client
	client, err := inspector.NewClient(inspector.ClientConf{
		APIKey:   apiKey,
		Instance: instance,
		Timeout:  30 * time.Second,
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Create visit (server assigns defaults)
	ctx := context.Background()
	visit, err := client.Visit.AddVisit(ctx)
	if err != nil {
		log.Fatalf("Failed to create visit: %v", err)
	}

	// Output JSON result
	output, err := json.MarshalIndent(visit, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal result: %v", err)
	}
	fmt.Println(string(output))
}
