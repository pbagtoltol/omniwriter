package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/pbagtoltol/omniwriter/pkg/omniwriter"
)

func main() {
	// Read the schema
	schema, err := os.ReadFile("schema.json")
	if err != nil {
		fmt.Printf("Error reading schema: %v\n", err)
		os.Exit(1)
	}

	// Read the EDI input
	input, err := os.Open("input.edi")
	if err != nil {
		fmt.Printf("Error opening input file: %v\n", err)
		os.Exit(1)
	}
	defer input.Close()

	// Create transformation request
	req := omniwriter.TransformRequest{
		SourceFormat: omniwriter.FormatEDI,
		TargetFormat: omniwriter.FormatCSV,
		Mapping:      schema,
		Input:        input,
	}

	// Execute transformation with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println("Starting EDI to CSV transformation...")
	fmt.Println("=====================================")

	result, err := omniwriter.Transform(ctx, req)
	if err != nil {
		fmt.Printf("Transformation error: %v\n", err)
		os.Exit(1)
	}

	// Write output
	err = os.WriteFile("output.csv", result.Output, 0644)
	if err != nil {
		fmt.Printf("Error writing output: %v\n", err)
		os.Exit(1)
	}

	// Print results
	fmt.Println("\n✅ Transformation complete!")
	fmt.Printf("Records processed: %d\n", result.Stats.Records)
	fmt.Printf("Output size: %d bytes\n", len(result.Output))
	fmt.Printf("Output written to: output.csv\n")

	if len(result.Warnings) > 0 {
		fmt.Println("\nWarnings:")
		for _, w := range result.Warnings {
			fmt.Printf("  - %s\n", w.Message)
		}
	}

	// Display the output
	fmt.Println("\nGenerated CSV Output:")
	fmt.Println("=====================")
	fmt.Println(string(result.Output))
}
