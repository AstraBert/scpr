package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/cobra"
)

var url string
var showHelp bool

var rootCmd = &cobra.Command{
	Use:   "scpr",
	Short: "scpr is a simple web scraper based on Colly.",
	Long:  "scpr is web scraper based on Colly that returns HTML web pages as markdown text.",
	Run: func(cmd *cobra.Command, args []string) {
		if showHelp {
			_ = cmd.Help()
		} else if url == "" {
			fmt.Println("Missing required option: `--url/-u`")
			_ = cmd.Help()
		} else {
			Scrape(url)
		}
	},
}

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Start a stdio MCP server powered by scpr",
	Long:  "Start an MCP server (over stdio transport) that allows you to perform the web scraping operations, powered by scpr.",
	Run: func(cmd *cobra.Command, args []string) {
		server := GetMcpServer()
		log.Println("Starting scpr MCP server...")
		if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Oops. An error while executing scpr '%s'\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&url, "url", "u", "", "URL of the web page to scrape. Required.")
	rootCmd.Flags().BoolVarP(&showHelp, "help", "h", false, "Show the help message and exit.")

	rootCmd.AddCommand(mcpCmd)

	_ = rootCmd.MarkFlagRequired("pattern")
}
