package cmd

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ScrapeParams struct {
	Url string `json:"url" jsonschema:"URL of the page you want to scrape"`
}

func ScprMcp(ctx context.Context, req *mcp.CallToolRequest, args ScrapeParams) (*mcp.CallToolResult, any, error) {
	res, err := ScraperImpl(args.Url)
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
	server := mcp.NewServer(&mcp.Implementation{Name: "scpr-mcp", Version: "1.0.0"}, &mcp.ServerOptions{Instructions: "Use this MCP server in order to perform scrape the content of a web page based on its URL.", HasTools: true, HasPrompts: false, HasResources: false})
	mcp.AddTool(server, &mcp.Tool{Name: "scpr", Description: "`scpr` is a tool that allows scraping of web pages based on their URLs. It takes, as arguments, the URL of the page to scrape."}, ScprMcp)
	return server
}
