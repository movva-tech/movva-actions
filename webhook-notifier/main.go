package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const (
	URL = "localhost:8080/api/release/notify"
)

func main() {
	secret := os.Getenv("INPUT_WEBHOOK_SECRET")
	repo := os.Getenv("INPUT_REPOSITORY")
	version := os.Getenv("INPUT_VERSION")

	if secret == "" {
		fmt.Println("::error::Missing URL or Secret inputs")
		os.Exit(1)
	}

	payload := map[string]string{
		"repository": repo,
		"version":    version,
		"status":     "success",
		"message":    "New release bundled and published successfully.",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("::error::Failed to marshal JSON payload: %v\n", err)
		os.Exit(1)
	}

	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("::error::Failed to create HTTP request: %v\n", err)
		os.Exit(1)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+secret)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("::error::Failed to execute HTTP request: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		fmt.Printf("::error::Server responded with error status: %d\n", resp.StatusCode)
		os.Exit(1)
	}

	fmt.Printf("Successfully sent notification for %s release %s\n", repo, version)
}
