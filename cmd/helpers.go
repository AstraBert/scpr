package cmd

import (
	"log/slog"
	"regexp"
	"strings"
	"sync"
)

func stringLevelToSlogLevel(level string) slog.Leveler {
	level = strings.ToLower(level)
	switch level {
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "debug":
		return slog.LevelDebug
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func validateAllowedDomains(domains []string) []string {
	polishedDomains := make([]string, 0, len(domains))
	pattern := regexp.MustCompile(`^([A-Za-z0-9]([A-Za-z0-9-]{0,61}[A-Za-z0-9])?\.)+[A-Za-z]{2,6}$`)
	for _, domain := range domains {
		domain = strings.TrimPrefix(domain, "https://")
		if pattern.MatchString(domain) {
			polishedDomains = append(polishedDomains, domain)
		}
	}
	return polishedDomains
}

type outputFiles struct {
	mu    sync.RWMutex
	paths []string
}

func (o *outputFiles) append(path string) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.paths = append(o.paths, path)
}

func (o *outputFiles) getFiles() []string {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.paths
}

func newOutputFiles() *outputFiles {
	return &outputFiles{
		paths: []string{},
	}
}
