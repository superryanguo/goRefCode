package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/ThinkInAIXYZ/go-mcp/server"
	"github.com/ThinkInAIXYZ/go-mcp/transport"
)

type TimeRequest struct {
	Timezone string `json:"timezone" description:"timezone" required:"true"` // Use field tag to describe input schema
}

func main() {
	// Create SSE transport server
	transportServer, err := transport.NewSSEServerTransport("127.0.0.1:8071")
	if err != nil {
		log.Fatalf("Failed to create transport server: %v", err)
	}

	// Initialize MCP server
	mcpServer, err := server.NewServer(transportServer)
	if err != nil {
		log.Fatalf("Failed to create MCP server: %v", err)
	}

	// Register time query tool
	tool, err := protocol.NewTool("current_time", "Get current time for specified timezone", TimeRequest{})
	if err != nil {
		log.Fatalf("Failed to create tool: %v", err)
		return
	}
	mcpServer.RegisterTool(tool, handleTimeRequest)

	// Start server
	if err = mcpServer.Run(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func handleTimeRequest(cxt context.Context, req *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	var timeReq TimeRequest
	if err := protocol.VerifyAndUnmarshal(req.RawArguments, &timeReq); err != nil {
		return nil, err
	}

	loc, err := time.LoadLocation(timeReq.Timezone)
	if err != nil {
		return nil, fmt.Errorf("invalid timezone: %v", err)
	}

	return &protocol.CallToolResult{
		Content: []protocol.Content{
			&protocol.TextContent{
				Type: "text",
				Text: time.Now().In(loc).String(),
			},
		},
	}, nil
}
