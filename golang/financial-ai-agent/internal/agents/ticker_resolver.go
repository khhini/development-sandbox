package agents

import (
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/geminitool"
)

func NewTickerResolverAgent(model model.LLM) (agent.Agent, error) {
	rootAgent, err := llmagent.New(llmagent.Config{
		Name:        "ticker_resolver",
		Model:       model,
		Description: "Resolves a company name to its stock ticker symbol using Google Search.",
		Instruction: `
		You are a stock ticker symbol resolver.
		When given a company name, use google_search to find its official stock ticker symbol.
		Search for '<company name> stock ticker symbol site:finance.yahoo.com OR site:google.com/finance'.
		Returun only the ticker symbol as plain text (e.g. 'INTC' for Intel, 'AAPL' for Apple).
		Do not include any explanation, punctuation, or extra text - just the ticker.
		`,
		Tools: []tool.Tool{
			&geminitool.GoogleSearch{},
		},
	})
	if err != nil {
		return nil, err
	}

	return rootAgent, nil
}
