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
var recursive bool
var allowedDomains []string
var maxDepth int
var parallel int
var output string
var logLevel string
var showHelp bool

var rootCmd = &cobra.Command{
	Use:   "scpr",
	Short: "scpr is a simple web scraper based on Colly.",
	Long:  "scpr is web scraper based on Colly that returns HTML web pages as markdown text.",
	Run: func(cmd *cobra.Command, args []string) {
		if showHelp {
			_ = cmd.Help()
			return
		} else if url == "" {
			fmt.Println("Missing required option: `--url/-u`")
			_ = cmd.Help()
			return
		} else if output == "" {
			fmt.Println("Missing required option: `--output/-o`")
			_ = cmd.Help()
			return
		} else {
			if recursive && len(allowedDomains) == 0 {
				fmt.Println("For security reasons, you need to provide a list of allowed domains to visit if you wish to scrape pages recursively.")
				_ = cmd.Help()
				return
			}
			files, err := ScraperImpl(url, output, logLevel, parallel, maxDepth, recursive, allowedDomains)
			if err != nil {
				fmt.Printf("An error occurred during scraping: %s\n", err.Error())
				return
			}
			if len(files) > 0 {
				fmt.Println("Ouput files:")
				for _, p := range files {
					fmt.Println(p)
				}
			} else {
				fmt.Println("No results produced.")
			}
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
	rootCmd.Flags().StringVarP(&output, "output", "o", "", "Output folder to save the markdown files to. Required.")
	rootCmd.Flags().StringVarP(&logLevel, "log", "l", "info", "Logging level. Defaults to info. Allowed values: 'info', 'debug', 'warn' and 'error'.")
	rootCmd.Flags().IntVarP(&parallel, "parallel", "p", 1, "Maximum number of threads for paralle scraping. Defaults to 1.")
	rootCmd.Flags().StringSliceVarP(&allowedDomains, "allowed", "a", []string{}, "Allowed domains (required for recursive scraping).")
	rootCmd.Flags().IntVarP(&maxDepth, "max", "m", 1, "Maximum depth of linked pages to scrape. Defaults to 1 (scrapes only the original page)")
	rootCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "Whether or not to scrape the pages recursively (scrape a page and all the pages linked to it). Only works if a list of allowed domains is passed")
	rootCmd.Flags().BoolVarP(&showHelp, "help", "h", false, "Show the help message and exit.")

	rootCmd.AddCommand(mcpCmd)

	_ = rootCmd.MarkFlagRequired("pattern")
}
