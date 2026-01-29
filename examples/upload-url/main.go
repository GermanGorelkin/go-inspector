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
	// Define flags
	url := flag.String("url", "", "Image URL to upload (required)")
	flag.Parse()

	// Validate required flags
	if *url == "" {
		fmt.Fprintf(os.Stderr, "Error: -url flag is required\n\n")
		flag.Usage()
		os.Exit(1)
	}

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

	// Upload image by URL
	ctx := context.Background()
	image, err := client.Image.UploadByURL(ctx, *url)
	if err != nil {
		log.Fatalf("Failed to upload image: %v", err)
	}

	// Output JSON result
	output, err := json.MarshalIndent(image, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal result: %v", err)
	}
	fmt.Println(string(output))
}
