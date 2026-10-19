package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

// bannedLibraries contains normalized Python package names that must not be used.
// Go slices cannot be constants; keep this list unexported and modify it only here.
var bannedLibraries = []string{}

func main() {
	reqFile := os.Getenv("INPUT_REQUIREMENTS_FILE")

	if len(bannedLibraries) == 0 {
		fmt.Println("No blacklist provided. Skipping scan.")
		return
	}

	file, err := os.Open(reqFile)
	if err != nil {
		fmt.Printf("::error::Could not open file %s: %v\n", reqFile, err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Extract just the package name (e.g., "flask==2.0" -> "flask")
		pkg := strings.Split(line, "=")[0]
		pkg = strings.Split(pkg, ">")[0]
		pkg = strings.Split(pkg, "<")[0]
		pkg = strings.TrimSpace(strings.ToLower(pkg))

		if slices.Contains(bannedLibraries, pkg) {
			// The ::error:: syntax tells GitHub Actions to highlight this as a red error in the UI
			fmt.Printf("::error::Blacklisted package found: %s\n", pkg)
			os.Exit(1)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("::error::Error reading file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Scan passed. No blacklisted packages found.")
}
