package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/germangorelkin/go-inspector/inspector"
)

func main() {
	// Define flags
	images := flag.String("images", "", "Comma-separated image IDs (required)")
	skuID := flag.Int("sku", 0, "SKU ID that was mistakenly recognized (required)")
	scene := flag.String("scene", "", "Scene UUID (required)")
	message := flag.String("message", "", "Error description (optional)")
	flag.Parse()

	// Validate required flags
	if *images == "" {
		fmt.Fprintf(os.Stderr, "Error: -images flag is required\n\n")
		flag.Usage()
		os.Exit(1)
	}
	if *skuID == 0 {
		fmt.Fprintf(os.Stderr, "Error: -sku flag is required\n\n")
		flag.Usage()
		os.Exit(1)
	}
	if *scene == "" {
		fmt.Fprintf(os.Stderr, "Error: -scene flag is required\n\n")
		flag.Usage()
		os.Exit(1)
	}

	// Parse image IDs
	imageIDStrs := strings.Split(*images, ",")
	var imageIDs []int
	for _, idStr := range imageIDStrs {
		idStr = strings.TrimSpace(idStr)
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Fatalf("Invalid image ID '%s': %v", idStr, err)
		}
		imageIDs = append(imageIDs, id)
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

	// Build recognition error request
	req := &inspector.RecognitionErrorRequest{
		Images:  imageIDs,
		SkuId:   *skuID,
		Scene:   *scene,
		Message: *message,
	}

	// Submit recognition error
	ctx := context.Background()
	resp, err := client.Recognize.RecognitionError(ctx, req)
	if err != nil {
		log.Fatalf("Failed to submit recognition error: %v", err)
	}

	// Output JSON result
	output, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal result: %v", err)
	}
	fmt.Println(string(output))
}
