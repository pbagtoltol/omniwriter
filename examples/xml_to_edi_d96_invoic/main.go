package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pbagtoltol/omniwriter/pkg/omniwriter"
)

func main() {
	schema, err := os.ReadFile("schema.json")
	if err != nil {
		fmt.Printf("Error reading schema: %v\n", err)
		os.Exit(1)
	}

	input, err := os.Open("input.xml")
	if err != nil {
		fmt.Printf("Error opening input: %v\n", err)
		os.Exit(1)
	}
	defer input.Close()

	req := omniwriter.TransformRequest{
		SourceFormat: omniwriter.FormatXML,
		TargetFormat: omniwriter.FormatEDI,
		Mapping:      schema,
		Input:        input,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println("XML to EDIFACT D96A INVOIC Transformation")
	fmt.Println("==========================================")
	fmt.Println()

	result, err := omniwriter.Transform(ctx, req)
	if err != nil {
		fmt.Printf("Transformation error: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile("output.edi", result.Output, 0644); err != nil {
		fmt.Printf("Error writing output: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Transformation complete")
	fmt.Printf("Records processed: %d\n", result.Stats.Records)
	fmt.Printf("Output size: %d bytes\n", len(result.Output))
	fmt.Printf("Warnings: %d\n\n", len(result.Warnings))

	// Display the EDIFACT output
	fmt.Println("Generated EDIFACT INVOIC D96A:")
	fmt.Println("==============================")
	output := string(result.Output)
	segments := strings.Split(output, "'")
	for i, seg := range segments {
		if strings.TrimSpace(seg) != "" {
			fmt.Printf("%2d: %s'\n", i+1, seg)
		}
	}
	fmt.Println("==============================")
	fmt.Println()

	// Show message structure
	fmt.Println("Message Structure:")
	fmt.Println("  UNH - Message Header (INVOIC D96A)")
	fmt.Println("  BGM - Beginning of Message")
	fmt.Println("  DTM - Document Date")
	fmt.Println("  NAD - Buyer Party")
	fmt.Println("  NAD - Supplier Party")
	fmt.Println("  LIN/QTY/PRI/MOA - Line Item Details (2 items)")
	fmt.Println("  UNS - Section Control")
	fmt.Println("  MOA - Total Amount")
	fmt.Println("  UNT - Message Trailer")
	fmt.Println()
	fmt.Printf("Output saved to: output.edi\n")
}
