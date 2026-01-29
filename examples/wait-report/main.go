package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/germangorelkin/go-inspector/inspector"
)

func main() {
	// Define flags
	reportID := flag.Int("id", 0, "Report ID (required)")
	interval := flag.Duration("interval", 2*time.Second, "Polling interval")
	timeout := flag.Duration("timeout", 60*time.Second, "Overall timeout")
	reportType := flag.String("type", "", "Report type for parsing: FACING_COUNT, PRICE_TAGS, REALOGRAM (optional)")
	flag.Parse()

	// Validate required flags
	if *reportID == 0 {
		fmt.Fprintf(os.Stderr, "Error: -id flag is required\n\n")
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

	// Wait for report with progress callback
	ctx := context.Background()
	opts := &inspector.ReportWaitOptions{
		Interval: *interval,
		Timeout:  *timeout,
		OnProgress: func(r *inspector.Report) {
			fmt.Fprintf(os.Stderr, "Report %d status: %s\n", r.ID, r.Status)
		},
	}

	report, err := client.Report.WaitForReport(ctx, *reportID, opts)
	if err != nil {
		log.Fatalf("Failed to wait for report: %v", err)
	}

	fmt.Fprintf(os.Stderr, "Report ready!\n")

	// Parse report based on type if specified
	var result interface{}
	if *reportType != "" {
		switch strings.ToUpper(*reportType) {
		case inspector.ReportTypeFACING_COUNT:
			parsed, err := client.Report.ToFacingCount(report.Json)
			if err != nil {
				log.Fatalf("Failed to parse FACING_COUNT report: %v", err)
			}
			result = parsed
		case inspector.ReportTypePRICE_TAGS:
			parsed, err := client.Report.ToPriceTags(report.Json)
			if err != nil {
				log.Fatalf("Failed to parse PRICE_TAGS report: %v", err)
			}
			result = parsed
		case inspector.ReportTypeREALOGRAM:
			parsed, err := client.Report.ToRealogram(report.Json)
			if err != nil {
				log.Fatalf("Failed to parse REALOGRAM report: %v", err)
			}
			result = parsed
		default:
			log.Fatalf("Unknown report type: %s (valid: FACING_COUNT, PRICE_TAGS, REALOGRAM)", *reportType)
		}
	} else {
		result = report
	}

	// Output JSON result
	output, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal result: %v", err)
	}
	fmt.Println(string(output))
}
