package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/googleapis/mcp-toolbox-sdk-go/tbadk"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/session"
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

	appName := "hotel_assistant"
	userID := "user-123"

	sessionService := session.InMemoryService()
	resp, err := sessionService.Create(ctx, &session.CreateRequest{
		AppName: appName,
		UserID:  userID,
	})
	if err != nil {
		log.Fatalf("Failed to create the sessiokn service: %v", err)
	}
	session := resp.Session

	r, err := runner.New(runner.Config{
		AppName:        appName,
		Agent:          llmagent,
		SessionService: sessionService,
	})
	if err != nil {
		log.Fatalf("Failed to create runner: %v", err)
	}

	for i, query := range queriesAdk {
		fmt.Printf("\n=== Query %d: %s ===\n", i+1, query)
		userMsg := genai.NewContentFromText(query, genai.RoleUser)

		streamingMode := agent.StreamingModeSSE
		for event, err := range r.Run(ctx, userID, session.ID(), userMsg, agent.RunConfig{
			StreamingMode: streamingMode,
		}) {
			if err != nil {
				fmt.Printf("\nAGENT_ERROR: %v\n", err)
			} else {
				if event.LLMResponse.Content != nil {
					for _, p := range event.LLMResponse.Content.Parts {
						// if its running in streaming mode, don't print the non partial llmResponses
						if streamingMode != agent.StreamingModeSSE || event.LLMResponse.Partial {
							fmt.Print(p.Text)
						}
					}
				}
			}
		}
	}

	fmt.Println("\n", strings.Repeat("-", 80)+"\n")
}
