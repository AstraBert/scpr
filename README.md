# scpr

`scpr` is a simple and straightforward webscraping CLI tool made to scrape page as markdown content, and mainly developed to be used by coding agents either as an [MCP server](./.mcp.json) or as a [skill](.claude/skills/web-scraping/SKILL.md).

`scpr` is written in Go and based on [colly](https://github.com/gocolly/colly) for web scraping and [`html-to-markdown`](https://github.com/JohannesKaufmann/html-to-markdown) for converting HTML pages to markdown.

## Installation

Install with Go (v1.24+ required):

```bash
go install github.com/AstraBert/scpr
```

## Usage

**As a CLI tool**

```bash
scpr --url https://example.com
```

This will print markdown text to stdout, which you can easily redirect to a file:

```bash
scpr --url https://example.com > example.md
```

**As a stdio MCP server**

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

## Roadmap

- [ ] Handle retries
- [ ] Add recursive scraping based on allowed domains
- [ ] Add an NPM-installable version
- [ ] Add comprehensive testing
