package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/pbagtoltol/omniwriter"
)

func main() {
	ctx := context.Background()

	schema, err := os.ReadFile("schema.json")
	if err != nil {
		fmt.Printf("Error reading schema: %v\n", err)
		os.Exit(1)
	}

	input, err := os.Open("input.txt")
	if err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		os.Exit(1)
	}
	defer input.Close()

	req := omniwriter.TransformRequest{
		SourceFormat: omniwriter.FormatCSV,
		TargetFormat: omniwriter.FormatEDI,
		Mapping:      schema,
		Input:        input,
	}

	fmt.Println("Flat File to EDI X12 850 Purchase Order Transformation")
	fmt.Println("======================================================")
	fmt.Println()

	result, err := omniwriter.Transform(ctx, req)
	if err != nil {
		fmt.Printf("Transformation error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Transformation complete")
	fmt.Printf("Records processed: %d\n", result.Stats.Records)
	fmt.Println()
	fmt.Println("EDI Output:")
	fmt.Println(strings.TrimSpace(string(result.Output)))
}
