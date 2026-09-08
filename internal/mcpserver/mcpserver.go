package mcpserver

import (
	"context"
	"ct-go-web-starter/internal/service"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const ServerName = "ct-go-web-starter"

type greetInput struct {
	Name string `json:"name" jsonschema:"the name of the person to greet"`
}

type greetOutput struct {
	Message string `json:"message"`
}

func New() *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    ServerName,
		Version: service.Version,
	}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "greet",
		Description: "Greet someone by name.",
	}, handleGreet)

	return server
}

// A returned error becomes a tool-level error the model can see and correct,
// not a protocol error — so domain refusals travel back as results.
func handleGreet(ctx context.Context, req *mcp.CallToolRequest, in greetInput) (*mcp.CallToolResult, greetOutput, error) {
	greeting, err := service.Greet(in.Name)
	if err != nil {
		return nil, greetOutput{}, err
	}

	return nil, greetOutput{Message: greeting.Message}, nil
}
