package agents

import (
	"fmt"
	"os"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/agenttool"
	"google.golang.org/adk/tool/functiontool"
	"google.golang.org/adk/tool/mcptoolset"
)

func NewRootAgent(model model.LLM) (agent.Agent, error) {
	alphavantageMCPToolsets, err := mcptoolset.New(mcptoolset.Config{
		Client: mcp.NewClient(&mcp.Implementation{
			Name: "aplhavantage-client",
		}, nil),
		Transport: &mcp.StreamableClientTransport{
			Endpoint: fmt.Sprintf("https://mcp.alphavantage.co/mcp?apikey=%s", os.Getenv("ALPHAVANTAGE_API_KEY")),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create mcptoolset: %v", err)
	}

	greetingTool, err := functiontool.New(functiontool.Config{
		Name:        "greeting",
		Description: "Greeting to user",
	}, func(ctx tool.Context, args struct{ Msg string }) (string, error) {
		return fmt.Sprintf("Hi, %s you are awesome, have a nice day", args.Msg), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create greeting tool: %v", err)
	}

	tickerResolverAgent, err := NewTickerResolverAgent(model)
	if err != nil {
		return nil, fmt.Errorf("failed to create ticker_resolver_agent: %v", err)
	}

	tickerResolverAgentTool := agenttool.New(tickerResolverAgent, nil)

	getCurrentDatetimeTool, err := functiontool.New(functiontool.Config{
		Name:        "get_current_datetime",
		Description: "Returns the current UTC date and time",
	}, func(ctx tool.Context, args struct{}) (string, error) {
		return time.Now().Format("2006-01-02 15:04:05"), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create get_current_datetime: %v", err)
	}

	rootAgent, err := llmagent.New(llmagent.Config{
		Name:        "alpha_vantage_agent",
		Model:       model,
		Description: "An agent that provides financial data and stock analysist usign Alpha Vantage",
		Instruction: `
		You are a specialized Aplha Vantage Financial Analysist Assistant.
		Your goal is to provide accurate stock market data, fundamental insights, and financial trends using the provided Alpha Vantage tools.

		TOOLS:
		The Alpha Vantage MCP server provides tools across several categories:
		1. 'core_stock_apis': Real-time and historical stock quotes (e.g., 'GLOBAL_QUOTE'. 'TIME_SERIES_INTRADAY').
		2. 'fundamental_data': Company metrics, earnings, and financial statements (e.g. 'COMPANY_OVERVIEW', 'EARNINGS').
		3. 'alpha_intelligence': Market news and sentiment analysist (e.g., 'NEWS_SENTIMENT').
		4. 'forex' and 'cryptocurrenies': Exchange rates for traditional and digital currencies.
		5. 'economic_indicators' and 'technical_indicators': Broad market data and analytical indicators (e.g., 'SMA', 'REAL_GDP').
		6. 'get_current_datetime': Returns the current UTC date and time. Use this whenever you need today's date, the current year, or to calculate date ranges (e.g., 'last 30 days', 'year to date').

		GUIDELINES:
		- CORE QUOTES: For simple price checks, always use 'GLOBAL_QUOTE'.
		- FUNDAMENTAL ANALYSIST: when asked about a company's health or financial status, use 'COMPANY_OVERVIEW' and 'EARNINGS'.
		- SENTIMENT: Use 'NEWS_SENTIMENT' to gauge the current market mood for a specific ticker.
		- TICKER RESOLUTION: When a company name is given instead of ticker symbol, call the 'ticker_resolver' agent tool 
				to resolve it (e.g. 'Intel' -> 'INTC'), then use the retunred ticker with Alpha Vantage tools.
		- RATE LIMIT: The Alpha Vantage API enforces a rate limit of 1 request per 3 seconds on free-tier keys.
				Wait atleast 3 seconds between consecutive tool calls to avoid errors.
		- SCOPE: Only assist with financial and stock-market related queries. Politely decline all other topics.
		`,
		Toolsets: []tool.Toolset{
			alphavantageMCPToolsets,
		},
		Tools: []tool.Tool{
			greetingTool,
			getCurrentDatetimeTool,
			tickerResolverAgentTool,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create alpha_vantage_agent: %v", err)
	}

	return rootAgent, nil
}
