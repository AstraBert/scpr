package cmd

import (
	"fmt"
	"log/slog"
	"slices"
	"sync"
	"testing"
)

func TestStringLevelToSlogLevel(t *testing.T) {
	if level := stringLevelToSlogLevel("info"); level != slog.LevelInfo {
		t.Fatalf("Expected info to be %d, got %d", slog.LevelInfo, level)
	}
	if level := stringLevelToSlogLevel("error"); level != slog.LevelError {
		t.Fatalf("Expected info to be %d, got %d", slog.LevelError, level)
	}
	if level := stringLevelToSlogLevel("debug"); level != slog.LevelDebug {
		t.Fatalf("Expected info to be %d, got %d", slog.LevelDebug, level)
	}
	if level := stringLevelToSlogLevel("warn"); level != slog.LevelWarn {
		t.Fatalf("Expected info to be %d, got %d", slog.LevelWarn, level)
	}
	if level := stringLevelToSlogLevel("default"); level != slog.LevelInfo {
		t.Fatalf("Expected info to be %d, got %d", slog.LevelInfo, level)
	}
}

func TestValidateAllowedDomains(t *testing.T) {
	testCases := []struct {
		allowedDomains  []string
		expectedDomains []string
	}{
		{allowedDomains: []string{"hello.com", "my123.bye", "123.456"}, expectedDomains: []string{"hello.com", "my123.bye"}},
		{allowedDomains: []string{"www.hello.com", "xyz.my123.bye", "123.not.12"}, expectedDomains: []string{"www.hello.com", "xyz.my123.bye"}},
		{allowedDomains: []string{"https://www.hello.com", "https://123.564"}, expectedDomains: []string{"www.hello.com"}},
	}
	for _, tc := range testCases {
		validated := validateAllowedDomains(tc.allowedDomains)
		if !slices.Equal(validated, tc.expectedDomains) {
			t.Errorf("Expected %v, got %v", tc.expectedDomains, validated)
		}
	}
}

func TestOutputFilesRaceConditions(t *testing.T) {
	outputFiles := newOutputFiles()
	var wg sync.WaitGroup

	for i := range 10 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := range 100 {
				key := fmt.Sprintf("key-%d-%d", id, j)
				outputFiles.append(key)
			}
		}(i)
	}

	for i := range 10 {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for range 100 {
				_ = outputFiles.getFiles()
			}
		}(i)
	}

	wg.Wait()
}
