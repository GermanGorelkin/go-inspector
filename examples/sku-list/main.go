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
	offset := flag.Int("offset", 0, "Pagination offset")
	limit := flag.Int("limit", 10, "Items per page")
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

	// Get SKU page
	ctx := context.Background()
	pag, err := client.Sku.GetSKU(ctx, *offset, *limit)
	if err != nil {
		log.Fatalf("Failed to get SKU list: %v", err)
	}

	// Parse SKU results
	skus, err := client.Sku.ToSku(pag.Results)
	if err != nil {
		log.Fatalf("Failed to parse SKU results: %v", err)
	}

	// Build response with pagination info
	result := map[string]interface{}{
		"count":    pag.Count,
		"next":     pag.Next,
		"previous": pag.Previous,
		"offset":   *offset,
		"limit":    *limit,
		"skus":     skus,
	}

	// Output JSON result
	output, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal result: %v", err)
	}
	fmt.Println(string(output))
}
