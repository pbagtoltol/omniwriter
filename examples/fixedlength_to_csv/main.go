package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pbagtoltol/omniwriter"
)

func main() {
	schema, err := os.ReadFile("transform.json")
	if err != nil {
		log.Fatalf("Failed to read schema: %v", err)
	}

	input, err := os.Open("transactions.txt")
	if err != nil {
		log.Fatalf("Failed to open input: %v", err)
	}
	defer input.Close()

	result, err := omniwriter.Transform(context.Background(),
		omniwriter.TransformRequest{
			SourceFormat: omniwriter.FormatCustom, // fixedlength2
			TargetFormat: omniwriter.FormatCSV,
			Mapping:      schema,
			Input:        input,
		})

	if err != nil {
		log.Fatalf("Transform failed: %v", err)
	}

	if err := os.WriteFile("output.csv", result.Output, 0644); err != nil {
		log.Fatalf("Failed to write output: %v", err)
	}

	fmt.Println("Transformation complete")
	fmt.Printf("Processed %d records\n", result.Stats.Records)
	fmt.Println("Output written to output.csv")
}
