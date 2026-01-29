package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/germangorelkin/go-inspector/inspector"
)

func main() {
	// Define flags
	filePath := flag.String("file", "", "Path to local image file (required)")
	flag.Parse()

	// Validate required flags
	if *filePath == "" {
		fmt.Fprintf(os.Stderr, "Error: -file flag is required\n\n")
		flag.Usage()
		os.Exit(1)
	}

	// Get credentials from environment
	apiKey := os.Getenv("API_KEY")
	instance := os.Getenv("INSTANCE")
	if apiKey == "" || instance == "" {
		log.Fatal("API_KEY and INSTANCE environment variables must be set")
	}

	// Open file
	file, err := os.Open(*filePath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// Create client
	client, err := inspector.NewClient(inspector.ClientConf{
		APIKey:   apiKey,
		Instance: instance,
		Timeout:  30 * time.Second,
	})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Upload image from file
	ctx := context.Background()
	filename := filepath.Base(*filePath)
	image, err := client.Image.Upload(ctx, file, filename)
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
