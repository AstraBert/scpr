package cmd

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestMcpServer(t *testing.T) {
	server := GetMcpServer()
	ctx := context.Background()
	serverTransport, clientTransport := mcp.NewInMemoryTransports()
	serverSession, err := server.Connect(ctx, serverTransport, nil)
	if err != nil {
		t.Errorf("Expecting server to be able to connect, got error %s", err.Error())
		log.Fatal(err)
	}
	client := mcp.NewClient(&mcp.Implementation{Name: "client"}, nil)
	clientSession, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Errorf("Expecting client to be able to connect, got error %s", err.Error())
		log.Fatal(err)
	}
	res, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      "scpr",
		Arguments: map[string]any{"url": "https://example-files.online-convert.com/document/txt/example.txt", "parallel": 1, "max": 1, "allowed": []string{}, "recursive": false, "output": "test"},
	})
	if err != nil {
		t.Errorf("Expecting tool call to succeed, got error %s", err.Error())
		log.Fatal(err)
	}
	if len(res.Content) != 1 {
		t.Errorf("Expecting the returned content to have length 1, got %d", len(res.Content))
		log.Fatal("Wrong content length")
	}
	content := res.Content[0]
	typedContent, ok := content.(*mcp.TextContent)
	if !ok {
		t.Error("Expecting content to be of type TextContent, but it is not")
		log.Fatal("Wrong content type")
	}
	text := typedContent.Text
	if text != fmt.Sprintf("Scraping results saved in the following files:\n- %s", strings.Join([]string{"test/document_txt_example.md"}, "\n- ")) {
		t.Errorf("Got an unexpected result: %s", text)
		log.Fatal("Wrong result")
	}

	_ = clientSession.Close()
	_ = serverSession.Wait()
}
