package cmd

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ScrapeParams struct {
	Url            string   `json:"url" jsonschema:"required,description=URL of the web page to scrape"`
	Recursive      bool     `json:"recursive" jsonschema:"description=Whether or not to scrape the pages recursively (scrape a page and all the pages linked to it). Only works if a list of allowed domains is passed"`
	AllowedDomains []string `json:"allowed" jsonschema:"description=Allowed domains (required for recursive scraping)"`
	MaxDepth       int      `json:"max" jsonschema:"description=Maximum depth of linked pages to scrape. Should default to 1 (scrapes only the original page)"`
	Parallel       int      `json:"parallel" jsonschema:"description=Maximum number of threads for parallel scraping. Should default to 1"`
	Output         string   `json:"output" jsonschema:"required,description=Output folder to save the markdown files to"`
}

func ScprMcp(ctx context.Context, req *mcp.CallToolRequest, args ScrapeParams) (*mcp.CallToolResult, any, error) {
	res, err := ScraperImpl(args.Url, args.Output, "debug", args.Parallel, args.MaxDepth, args.Recursive, args.AllowedDomains)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("An error occurred: %s\n", err.Error())},
			},
		}, nil, nil
	}
	contents := []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Page content for %s:\n\n```md\n%s\n```", args.Url, res)}}
	return &mcp.CallToolResult{
		Content: contents,
	}, nil, nil
}

func GetMcpServer() *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{Name: "scpr-mcp", Version: "1.0.0"}, &mcp.ServerOptions{Instructions: "Use this MCP server in order to perform scrape the content of web pages based on their URLs.", HasTools: true, HasPrompts: false, HasResources: false})
	mcp.AddTool(server, &mcp.Tool{Name: "scpr", Description: "`scpr` is a tool that allows scraping of web pages based on their URLs. It takes, as arguments, the URL of the page to scrape and an output folder to save the scraped content. Optionally, it supports recursive scraping of linked pages with configurable depth, parallel processing, and domain restrictions."}, ScprMcp)
	return server
}
