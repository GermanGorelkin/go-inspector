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
	types := flag.String("types", "FACING_COUNT", "Comma-separated report types")
	webhook := flag.String("webhook", "", "Webhook URL for async notifications (optional)")
	visit := flag.Int("visit", 0, "Visit ID to associate with recognition (optional)")
	flag.Parse()

	// Validate required flags
	if *images == "" {
		fmt.Fprintf(os.Stderr, "Error: -images flag is required\n\n")
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

	// Parse report types
	reportTypes := strings.Split(*types, ",")
	for i, rt := range reportTypes {
		reportTypes[i] = strings.TrimSpace(rt)
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

	// Build recognition request
	req := inspector.RecognizeRequest{
		Images:      imageIDs,
		ReportTypes: reportTypes,
		Webhook:     *webhook,
	}
	if *visit > 0 {
		req.Visit = *visit
	}

	// Trigger recognition
	ctx := context.Background()
	resp, err := client.Recognize.Recognize(ctx, req)
	if err != nil {
		log.Fatalf("Failed to trigger recognition: %v", err)
	}

	// Output JSON result
	output, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal result: %v", err)
	}
	fmt.Println(string(output))
}
