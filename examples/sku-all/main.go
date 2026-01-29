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
	pageSize := flag.Int("page-size", 100, "Items per page")
	max := flag.Int("max", 0, "Maximum total items to fetch (0 = all)")
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

	// Fetch SKUs using iterator
	ctx := context.Background()
	iterator := client.Sku.IterateSKU(ctx, *pageSize)

	var allSKUs []inspector.Sku
	pageNum := 0

	for {
		pageNum++
		fmt.Fprintf(os.Stderr, "Fetching page %d...\n", pageNum)

		page, err := iterator.Next()
		if err != nil {
			log.Fatalf("Failed to fetch SKU page: %v", err)
		}
		if page == nil {
			break // No more pages
		}

		allSKUs = append(allSKUs, page...)
		fmt.Fprintf(os.Stderr, "Retrieved %d SKUs (total so far: %d)\n", len(page), len(allSKUs))

		// Check if we've reached the maximum
		if *max > 0 && len(allSKUs) >= *max {
			allSKUs = allSKUs[:*max]
			fmt.Fprintf(os.Stderr, "Reached maximum limit of %d SKUs\n", *max)
			break
		}
	}

	fmt.Fprintf(os.Stderr, "Total SKUs retrieved: %d\n", len(allSKUs))

	// Output JSON result
	output, err := json.MarshalIndent(allSKUs, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal result: %v", err)
	}
	fmt.Println(string(output))
}
