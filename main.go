package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Printf("Usage: %s <giturl>", os.Args[0])
		os.Exit(1)
	}

	repoUrl := args[0]

	if !isValidUrl(repoUrl) {
		fmt.Printf("Invalid URL: %s\n", repoUrl)
		os.Exit(1)
	}

	repoName, err := getRepoName(repoUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting repository name: %v\n", err)
		os.Exit(1)
	}

	// Construct the GitHub API URL
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s", repoName)

	resp, err := http.Get(apiURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error making HTTP request: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(os.Stderr, "HTTP request returned status code %d\n", resp.StatusCode)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(resp.Body)

	if !scanner.Scan() {
		fmt.Fprintf(os.Stderr, "Error scanning response body: %v\n", scanner.Err())
		os.Exit(1)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(scanner.Bytes(), &data); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing JSON: %v\n", err)
		os.Exit(1)
	}

	sizeVar, ok := data["size"].(float64)
	if !ok {
		fmt.Fprintf(os.Stderr, "Error parsing size from response\n")
		os.Exit(1)
	}

	sizeStr := convertSize(sizeVar)

	fmt.Printf("Repo size: %v\n", sizeStr)
}

func getRepoName(url string) (string, error) {
	if !strings.HasPrefix(url, "https://github.com/") {
		return "", errors.New("invalid repository URL")
	}

	// Strip the protocol and any trailing ".git" extension from the URL
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimSuffix(url, ".git")

	// Split the URL into its component parts
	parts := strings.Split(url, "/")
	if len(parts) < 3 {
		return "", errors.New("invalid repository URL")
	}

	username := parts[len(parts)-2]
	reponame := parts[len(parts)-1]

	return fmt.Sprintf("%s/%s", username, reponame), nil
}

func convertSize(size float64) string {
    if size < 1024 {
        return fmt.Sprintf("%.2f KB", size)
    } else if size < 1048576 {
        return fmt.Sprintf("%.2f MB", float64(size)/1024)
    } else if size < 1073741824 {
        return fmt.Sprintf("%.2f GB", float64(size)/1048576)
    } else if size < 1099511627776 {
        return fmt.Sprintf("%.2f TB", float64(size)/1073741824)
    } else {
        return fmt.Sprintf("%.2f PB", float64(size)/1099511627776)
    }
}

func isValidUrl(url string) bool {
	// Regular expression for validating a GitHub repository URL
	// Assumes that the URL has the format: https://github.com/<username>/<repository>
	r, err := regexp.Compile(`^https:\/\/github\.com\/[a-zA-Z0-9]+\/[a-zA-Z0-9_-]+$`)

	if err != nil {
		panic(err)
	}

	return r.MatchString(url)
}