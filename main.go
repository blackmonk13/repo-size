package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {

	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Printf("Usage: %s <giturl>", os.Args[0])
		os.Exit(1)
	}

	repoUrl := os.Args[1]

	repoName := getRepoName(repoUrl)

	// Construct the GitHub API URL
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s", repoName)

	// fmt.Println(repoName)
	// fmt.Println(apiURL)
	// fmt.Println(repoUrl)

	resp, err := http.Get(apiURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// fmt.Println("Response status:", resp.Status)

	scanner := bufio.NewScanner(resp.Body)

	scan := scanner.Scan()

	if !scan {
		os.Exit(1)
	}

	// fmt.Println(scanner.Text())

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	var data map[string]interface{}
	err = json.Unmarshal(scanner.Bytes(), &data)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing JSON: %v\n", err)
		os.Exit(1)
	}
	
	// fmt.Println(data["size"])
	// fmt.Printf("%T\n", data["size"])
	sizeVar, ok := data["size"].(float64)

	if !ok {
		os.Exit(1)
	}

	sizeStr := convertSize(sizeVar)

	// Print the size of the repository
	fmt.Printf("Repo size: %v\n", sizeStr)
}

func getRepoName(url string) string {
	// Strip the protocol and any trailing ".git" extension from the URL
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimSuffix(url, ".git")

	// Split the URL into its component parts
	parts := strings.Split(url, "/")

	username := parts[len(parts)-2]
	reponame := parts[len(parts)-1]
	// Return the repository name in the format "username/reponame"
	repository := []string{username, reponame}
	return strings.Join(repository, "/")
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
