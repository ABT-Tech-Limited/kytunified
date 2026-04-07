package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ABT-Tech-Limited/beosin-go"
	"github.com/ABT-Tech-Limited/kytunified/kyt"
	beosinprovider "github.com/ABT-Tech-Limited/kytunified/provider/beosin"
	"github.com/ABT-Tech-Limited/kytunified/registry"
)

func main() {
	// Get credentials from environment
	appID := os.Getenv("BEOSIN_APP_ID")
	appSecret := os.Getenv("BEOSIN_APP_SECRET")

	if appID == "" || appSecret == "" {
		log.Fatal("Please set BEOSIN_APP_ID and BEOSIN_APP_SECRET environment variables")
	}

	// Example 1: Direct provider creation
	fmt.Println("=== Example 1: Direct Provider Creation ===")
	directExample(appID, appSecret)

	// Example 2: Using registry
	fmt.Println("\n=== Example 2: Using Registry ===")
	registryExample(appID, appSecret)
}

func directExample(appID, appSecret string) {
	// Create beosin client
	client := beosin.NewClient(
		appID,
		appSecret,
		beosin.WithTimeout(30*time.Second),
		beosin.WithDebug(true),
	)

	// Create provider with client
	provider := beosinprovider.New(client, beosinprovider.WithV4())
	defer provider.Close()

	// Test provider configuration
	testProvider(provider)

	// Run assessments
	runAssessments(provider)
}

func registryExample(appID, appSecret string) {
	// Create beosin client
	client := beosin.NewClient(
		appID,
		appSecret,
		beosin.WithTimeout(30*time.Second),
		beosin.WithDebug(true),
	)

	// Register provider with the registry
	registry.MustRegisterBeosin(client, beosinprovider.WithV4())

	// Get provider from registry
	provider, err := registry.GetBeosin()
	if err != nil {
		log.Printf("Failed to get provider: %v", err)
		return
	}
	defer provider.Close()

	// List available providers
	fmt.Printf("Available providers: %v\n", registry.List())

	// Test provider configuration
	testProvider(provider)

	// Run assessments
	runAssessments(provider)
}

func testProvider(provider kyt.Provider) {
	ctx := context.Background()

	fmt.Printf("\n--- Provider Test ---\n")
	result := provider.Test(ctx)
	if result.Err != nil {
		fmt.Printf("Test inconclusive (non-business error): %v\n", result.Err)
		return
	}
	if !result.Valid {
		fmt.Printf("Configuration invalid: %s\n", result.Reason)
		return
	}
	fmt.Println("Configuration OK")
}

func runAssessments(provider kyt.Provider) {
	ctx := context.Background()

	fmt.Printf("Provider: %s\n", provider.Name())

	// Example address (Vitalik's address)
	testAddress := "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"

	// Assess address risk
	fmt.Printf("\n--- Address Risk Assessment ---\n")
	result, err := provider.AddressRisk(ctx, &kyt.AddressRiskRequest{
		ChainID: kyt.ChainIDETH,
		Address: testAddress,
	})
	if err != nil {
		if kyt.IsRetryable(err) {
			fmt.Println("Assessment pending, retry later")
		} else if kyt.IsValidation(err) {
			fmt.Printf("Validation error: %v\n", err)
		} else {
			fmt.Printf("Error: %v\n", err)
		}
	} else {
		printResult(result)
	}
}

func printResult(result *kyt.RiskResult) {
	fmt.Printf("Risk Level: %s\n", result.Level)
	fmt.Printf("Score: %.2f\n", result.Score)

	fmt.Printf("Provider: %s (API: %s)\n",
		result.Metadata.Provider, result.Metadata.APIVersion)

	if result.Detail != nil {
		fmt.Printf("Detail: %+v\n", result.Detail)
	}

	// Example of threshold checking
	if result.IsHighRisk() {
		fmt.Println("\nWARNING: High risk detected!")
	}
	if result.IsCritical() {
		fmt.Println("\nCRITICAL: Immediate action required!")
	}
}
