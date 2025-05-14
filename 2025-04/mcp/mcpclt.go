package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/ThinkInAIXYZ/go-mcp/client"
	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/ThinkInAIXYZ/go-mcp/transport"
)

func main() {
	// Create SSE transport client
	transportClient, err := transport.NewSSEClientTransport("http://127.0.0.1:8071/sse")
	if err != nil {
		log.Fatalf("Failed to create transport client: %v", err)
	}

	// Initialize MCP client
	mcpClient, err := client.NewClient(transportClient)
	if err != nil {
		log.Fatalf("Failed to create MCP client: %v", err)
	}
	defer mcpClient.Close()

	// Get available tools
	ctx := context.Background()
	tools, err := mcpClient.ListTools(ctx)
	if err != nil {
		log.Fatalf("Failed to list tools: %v", err)
	}
	for _, tool := range tools.Tools {
		log.Printf("- %s: %s\n", tool.Name, tool.Description)
	}

	getTimeRequest := &protocol.CallToolRequest{
		Name:         tools.Tools[0].Name,
		RawArguments: json.RawMessage(`{"timezone": "Asia/Shanghai"}`),
	}
	result, err := mcpClient.CallTool(ctx, getTimeRequest)
	if err != nil {
		log.Fatalf("Failed to Get Srv Time Tool Calling: %v", err)
	}
	printToolResult(result)
}

func printToolResult(result *protocol.CallToolResult) {
	for _, content := range result.Content {
		if textContent, ok := content.(*protocol.TextContent); ok {
			log.Println(textContent.Text)
		} else {
			jsonBytes, _ := json.MarshalIndent(content, "", "  ")
			log.Println(string(jsonBytes))
		}
	}
}
