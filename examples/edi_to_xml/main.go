package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/pbagtoltol/omniwriter/pkg/omniwriter"
)

func main() {
	schema, err := os.ReadFile("schema.json")
	if err != nil {
		fmt.Printf("Error reading schema: %v\n", err)
		os.Exit(1)
	}

	input, err := os.Open("input.edi")
	if err != nil {
		fmt.Printf("Error opening input: %v\n", err)
		os.Exit(1)
	}
	defer input.Close()

	req := omniwriter.TransformRequest{
		SourceFormat: omniwriter.FormatEDI,
		TargetFormat: omniwriter.FormatXML,
		Mapping:      schema,
		Input:        input,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Println("EDI → XML Transformation")
	fmt.Println("=========================")

	result, err := omniwriter.Transform(ctx, req)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile("output.xml", result.Output, 0644); err != nil {
		fmt.Printf("Error writing output: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n✅ Success!\n")
	fmt.Printf("Records: %d\n", result.Stats.Records)
	fmt.Printf("Output: output.xml\n\n")
	fmt.Println("Generated XML:")
	fmt.Println(string(result.Output))
}
