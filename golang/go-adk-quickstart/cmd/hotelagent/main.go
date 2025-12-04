package main

import (
	"context"
	"log"
	"os"

	"github.com/googleapis/mcp-toolbox-sdk-go/tbadk"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher/adk"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/server/restapi/services"
	"google.golang.org/adk/tool"
	"google.golang.org/genai"
)

const systemPrompt = `
	You're a helpful hotel assistant. You handle hotel searching, booking, and 
	cancellations. When the user searches for a hotel, mention its name, id, location 
	and price tier. Always mension hotel ids while performing any searches. This is 
	very important for any operations. For any bookings or cancellations, please 
	provide the appropriate confirmation. Be sure to update checkin or checkout dates 
	if mentioned by the user.
	Don't ask for confirmations from the user.
	`

var queriesAdk = []string{
	"Find hotels in Basel.",
	"Find hotels with Base in its name.",
	"Can you book the hotel Hilton Basel for me?",
	"Oh wait, this is too expensive. Please cancel it.",
	"Please book the Hyatt Regency instead.",
	"My check in dates would be from April 10, 2024 to April 19, 2024.",
}

func main() {
	genaiKey := os.Getenv("GOOGLE_API_KEY")
	toolboxURL := "http://127.0.0.1:5000"
	ctx := context.Background()

	toolboxClient, err := tbadk.NewToolboxClient(toolboxURL)
	if err != nil {
		log.Fatalf("Failed to create MCP Toolbox client: %v", err)
	}

	toolsetName := "my-toolsets"
	mcpTools, err := toolboxClient.LoadToolset(toolsetName, ctx)
	if err != nil {
		log.Fatalf("Failed to load MCP toolset '%s': %v\nMake sure your Toolbox sever is running.", toolsetName, err)
	}

	model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
		APIKey: genaiKey,
	})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	tools := make([]tool.Tool, len(mcpTools))
	for i := range mcpTools {
		tools[i] = &mcpTools[i]
	}

	llmagent, err := llmagent.New(llmagent.Config{
		Name:        "hotel_assistant",
		Model:       model,
		Description: "Agent to answer questions about hotels.",
		Instruction: systemPrompt,
		Tools:       tools,
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	config := &adk.Config{
		AgentLoader: services.NewSingleAgentLoader(llmagent),
	}

	l := full.NewLauncher()
	err = l.Execute(ctx, config, os.Args[1:])
	if err != nil {
		log.Fatalf("run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}
