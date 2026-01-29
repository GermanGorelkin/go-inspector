package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/germangorelkin/go-inspector/inspector"
)

func main() {
	// Define flags
	url := flag.String("url", "", "Image URL to upload")
	file := flag.String("file", "", "Path to local image file")
	types := flag.String("types", "FACING_COUNT", "Comma-separated report types")
	wait := flag.Bool("wait", true, "Wait for reports to complete")
	retailChain := flag.String("retail-chain", "", "Retail chain identifier (optional)")
	flag.Parse()

	// Validate required flags - need either URL or file
	if *url == "" && *file == "" {
		fmt.Fprintf(os.Stderr, "Error: either -url or -file flag is required\n\n")
		flag.Usage()
		os.Exit(1)
	}
	if *url != "" && *file != "" {
		fmt.Fprintf(os.Stderr, "Error: cannot specify both -url and -file flags\n\n")
		flag.Usage()
		os.Exit(1)
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

	ctx := context.Background()

	// Step 1: Upload image
	fmt.Fprintf(os.Stderr, "Step 1: Uploading image...\n")
	var image inspector.Image
	if *url != "" {
		fmt.Fprintf(os.Stderr, "  Uploading from URL: %s\n", *url)
		image, err = client.Image.UploadByURL(ctx, *url)
	} else {
		fmt.Fprintf(os.Stderr, "  Uploading from file: %s\n", *file)
		f, err := os.Open(*file)
		if err != nil {
			log.Fatalf("Failed to open file: %v", err)
		}
		defer f.Close()
		filename := filepath.Base(*file)
		image, err = client.Image.Upload(ctx, f, filename)
	}
	if err != nil {
		log.Fatalf("Failed to upload image: %v", err)
	}
	fmt.Fprintf(os.Stderr, "  ✓ Image uploaded: ID=%d, Size=%dx%d\n", image.ID, image.Width, image.Height)

	// Step 2: Create visit
	fmt.Fprintf(os.Stderr, "\nStep 2: Creating visit...\n")
	visit, err := client.Visit.AddVisit(ctx)
	if err != nil {
		log.Fatalf("Failed to create visit: %v", err)
	}
	fmt.Fprintf(os.Stderr, "  ✓ Visit created: ID=%d\n", visit.ID)

	// Step 3: Trigger recognition
	fmt.Fprintf(os.Stderr, "\nStep 3: Triggering recognition...\n")
	fmt.Fprintf(os.Stderr, "  Report types: %s\n", strings.Join(reportTypes, ", "))
	recReq := inspector.RecognizeRequest{
		Images:      []int{image.ID},
		ReportTypes: reportTypes,
		Visit:       visit.ID,
		RetailChain: *retailChain,
	}
	recResp, err := client.Recognize.Recognize(ctx, recReq)
	if err != nil {
		log.Fatalf("Failed to trigger recognition: %v", err)
	}
	fmt.Fprintf(os.Stderr, "  ✓ Recognition started: ID=%d\n", recResp.ID)
	fmt.Fprintf(os.Stderr, "  Report IDs: %v\n", recResp.Reports)

	// Step 4: Wait for reports (if enabled)
	result := map[string]interface{}{
		"image":       image,
		"visit":       visit,
		"recognition": recResp,
		"reports":     make(map[string]interface{}),
	}

	if *wait {
		fmt.Fprintf(os.Stderr, "\nStep 4: Waiting for reports to complete...\n")
		reports := make(map[string]interface{})

		for reportType, reportID := range recResp.Reports {
			fmt.Fprintf(os.Stderr, "  Waiting for %s (ID=%d)...\n", reportType, reportID)

			opts := &inspector.ReportWaitOptions{
				Interval: 2 * time.Second,
				Timeout:  60 * time.Second,
				OnProgress: func(r *inspector.Report) {
					fmt.Fprintf(os.Stderr, "    Status: %s\n", r.Status)
				},
			}

			report, err := client.Report.WaitForReport(ctx, reportID, opts)
			if err != nil {
				log.Fatalf("Failed to wait for report %s: %v", reportType, err)
			}

			// Parse report based on type
			var parsed interface{}
			switch reportType {
			case inspector.ReportTypeFACING_COUNT:
				parsed, err = client.Report.ToFacingCount(report.Json)
			case inspector.ReportTypePRICE_TAGS:
				parsed, err = client.Report.ToPriceTags(report.Json)
			case inspector.ReportTypeREALOGRAM:
				parsed, err = client.Report.ToRealogram(report.Json)
			default:
				parsed = report.Json
			}
			if err != nil {
				log.Fatalf("Failed to parse report %s: %v", reportType, err)
			}

			reports[reportType] = parsed
			fmt.Fprintf(os.Stderr, "  ✓ %s ready\n", reportType)
		}

		result["reports"] = reports
		fmt.Fprintf(os.Stderr, "\n✓ All reports completed successfully!\n\n")
	} else {
		fmt.Fprintf(os.Stderr, "\nStep 4: Skipped (use -wait=true to wait for reports)\n\n")
	}

	// Output JSON result
	output, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal result: %v", err)
	}
	fmt.Println(string(output))
}
