# scpr

`scpr` is a simple and straightforward webscraping CLI tool made to scrape page as markdown content, and mainly developed to be used by coding agents either as an [MCP server](./.mcp.json) or as a [skill](.claude/skills/web-scraping/SKILL.md).

`scpr` is written in Go and based on [colly](https://github.com/gocolly/colly) for web scraping and [`html-to-markdown`](https://github.com/JohannesKaufmann/html-to-markdown) for converting HTML pages to markdown.

## Installation

Install with Go (v1.24+ required):

```bash
go install github.com/AstraBert/scpr
```

Here's the updated usage section:

## Usage

### As a CLI tool

Basic usage (scrape a single page):

```bash
scpr --url https://example.com --output ./scraped
```

This will scrape the page and save it as a markdown file in the `./scraped` folder.

**Recursive scraping**

To scrape a page and all linked pages within the same domain:

```bash
scpr --url https://example.com --output ./scraped --recursive --allowed example.com --max 3
```

**Parallel scraping**

Speed up recursive scraping with multiple threads:

```bash
scpr --url https://example.com --output ./scraped --recursive --allowed example.com --max 2 --parallel 5
```

**Additional options**

- `--log` - Set logging level (info, debug, warn, error)
- `--max` - Maximum depth of pages to follow (default: 1)
- `--parallel` - Number of concurrent threads (default: 1)
- `--allowed` - Allowed domains for recursive scraping (can be specified multiple times)

For more details, run:

```bash
scpr --help
```

### As a stdio MCP server

Start the MCP server with:

```bash
scpr mcp
```

And configure it in agents using:

```json
{
  "mcpServers": {
    "web-scraping": {
      "type": "stdio",
      "command": "scpr",
      "args": [
        "mcp"
      ],
      "env": {}
    }
  }
}
```

> _The above JSON snippet is reported as used by Claude Code, adapt it to your agent before using it_
