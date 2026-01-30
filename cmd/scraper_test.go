package cmd

import (
	"fmt"
	"os"
	"testing"
)

func TestScraperSimple(t *testing.T) {
	url := "https://example-files.online-convert.com/document/txt/example.txt"
	output := "test"
	result, err := ScraperImpl(url, output, "info", 1, 1, false, []string{})
	if err != nil {
		t.Fatalf("An error occurred: %s", err.Error())
	}
	if len(result) != 1 {
		t.Fatalf("Expected to retrieve one result, got %d", len(result))
	}
	for _, r := range result {
		if _, err := os.Stat(r); os.IsNotExist(err) {
			t.Fatalf("%s does not exist", r)
		}
	}
}

func TestScraperComplex(t *testing.T) {
	url := "https://clelia.dev/2025-09-20-gen-z-in-open-source"
	output := "test"
	allowedDomains := []string{"github.blog", "clelia.dev"}
	maxDepth := 2
	parallel := 2
	recursive := true
	results, err := ScraperImpl(url, output, "info", parallel, maxDepth, recursive, allowedDomains)
	if err != nil {
		t.Fatalf("An error occurred: %s", err.Error())
	}
	if len(results) != 16 {
		t.Fatalf("Expected to retrieve 16 results, got %d", len(results))
	}
	fmt.Println(results)
	for _, r := range results {
		if _, err := os.Stat(r); os.IsNotExist(err) {
			t.Fatalf("%s does not exist", r)
		}
	}
}

func TestScraperNotFound(t *testing.T) {
	url := "https://example-files.online-convert.com/document/txt/example.t"
	output := "test"
	result, err := ScraperImpl(url, output, "info", 1, 1, false, []string{})
	if err != nil {
		t.Fatalf("An error occurred: %s", err.Error())
	}
	if len(result) != 0 {
		t.Fatalf("Expected to retrieve 0 results, got %d", len(result))
	}
}

func TestScraperDomainNotAllowed(t *testing.T) {
	url := "https://clelia.dev/2025-09-20-gen-z-in-open-source"
	output := "test"
	allowedDomains := []string{"github.blog"}
	maxDepth := 2
	parallel := 2
	recursive := true
	_, err := ScraperImpl(url, output, "info", parallel, maxDepth, recursive, allowedDomains)
	if err == nil {
		t.Fatal("An error occurred was expected, none occurred")
	} else {
		if err.Error() != "Forbidden domain" {
			t.Fatalf("Expecting 'Forbidden domain' as error message, got %s", err.Error())
		}
	}
}
