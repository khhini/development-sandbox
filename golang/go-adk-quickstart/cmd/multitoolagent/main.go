package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher/adk"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/server/restapi/services"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
	"google.golang.org/genai"
)

type ToolResult struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func NewToolResult(status string, message string) ToolResult {
	return ToolResult{
		Status:  status,
		Message: message,
	}
}

type ToolInput struct {
	City string `json:"city"`
}

func getWeather(ctx tool.Context, input ToolInput) ToolResult {
	if strings.ToLower(input.City) != "new york" {
		return NewToolResult(
			"error",
			fmt.Sprintf("Weather information for %s is not available", input.City),
		)
	}

	return NewToolResult(
		"success",
		"The weather in New York is sunny with a temperature of 25 degrees Celcius (77 degrees Fahrenheit).",
	)
}

func getCurrentTime(ctx tool.Context, input ToolInput) ToolResult {
	if strings.ToLower(input.City) != "new york" {
		return NewToolResult(
			"error",
			fmt.Sprintf("Timezone information for %s is not available", input.City),
		)
	}
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		return NewToolResult(
			"error",
			fmt.Sprintf("Error loading %s Timezone", input.City),
		)
	}

	currentTime := time.Now().In(loc)

	return NewToolResult(
		"success",
		fmt.Sprintf("The current time in %s is now %s", input.City, currentTime.Format("2006-01-02 15:04:05 MST")),
	)
}

func main() {
	ctx := context.Background()

	model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
		APIKey: os.Getenv("GEMINI_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	getWeatherTool, _ := functiontool.New(functiontool.Config{
		Name:        "get_weather_tool",
		Description: "Retrieves the current weather report for a specified city.",
	}, getWeather)

	getCurrentTimeTool, _ := functiontool.New(functiontool.Config{
		Name:        "get_current_time_tool",
		Description: "Returns the current time in a specified city.",
	}, getCurrentTime)

	weatherTimeAgent, err := llmagent.New(llmagent.Config{
		Name:        "weather_time_agent",
		Model:       model,
		Description: "Agent to answer questions about the time and weather in a city.",
		Instruction: "You are a helpful agent who can answer user questions about the time and weather in a city.",
		Tools: []tool.Tool{
			getWeatherTool,
			getCurrentTimeTool,
		},
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	config := &adk.Config{
		AgentLoader: services.NewSingleAgentLoader(weatherTimeAgent),
	}

	l := full.NewLauncher()
	if err = l.Execute(ctx, config, os.Args[1:]); err != nil {
		log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}
