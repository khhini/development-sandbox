import asyncio
from typing import Optional

from google.adk.agents import Agent
from google.adk.sessions import InMemorySessionService
from google.adk.runners import Runner
from google.genai.types import Tool


# @title Define the get_weather Tool
def get_weather(city: str) -> dict:
    """Retrieves the current weather report for a specified city.

    Args:
        city (str): The name of the city (e.g., "New York", "London", "Tokyo").

    Returns:
        dict: A dictionary containing the wewather information.
            Inlcudes a 'status' key ('success' or 'error').
            If 'success', includes a 'report' key with weather details.
            If 'error', includes an 'error_message' key.
    """

    print(f"--- Tool: get_weather called for city: {city} ---")
    city_normalized = city.lower().replace(" ", "")

    mock_weather_db = {
        "newyork": {
            "status": "success",
            "report": "The weather in New York is suny with a temperature of 25C.",
        },
        "london": {
            "status": "success",
            "report": "It's cloudy in London with a temperature of 15C.",
        },
        "tokyo": {
            "status": "success",
            "report": "Tokyo is experiencing light rain and a temperature of 18C.",
        },
    }

    if city_normalized in mock_weather_db:
        return mock_weather_db[city_normalized]

    return {
        "status": "error",
        "error_message": f"Sorry, I don't have weather information for '{city}'",
    }


AGENT_MODEL = "gemini-2.0-flash"

weather_agent_2_0 = Agent(
    name="weather_agent_2_0",
    model=AGENT_MODEL,
    description="Provides weather information for specific citites.",
    instruction="You are a helpful weather assistant. "
    "When hte user asks for the weather in a specific city. "
    "use the 'get_weather tool to find the information. "
    "If the tool returns an error, inform the user politely. "
    "If the tool is successful, present the weather report clearly.",
    tools=[get_weather],
)


def say_hello(name: Optional[str] = None) -> str:
    """Provides a simple greeting. IF a name is provided, it will be used.

    Args:
        name (str, optional): The name of the person to greet. Defaults to a generic greeting if not provided.

    returns:
        str: A friendly greeting message.
    """

    greeting = "Hello there!"  # Default greeting

    if name:
        greeting = f"Hello, {name}!"

    return greeting


def say_goodbye() -> str:
    """Provides a simple farewell message to conclude the conversation."""

    return "Goodbye! Have a great day."


greeting_agent = Agent(
    model=AGENT_MODEL,
    name="greeting_agent",
    instruction="You are the Greeting Agent. Your ONLY task is to provide a friendly greeting to the user. "
    "Use the 'say_hello' tool to generate the greeting. "
    "If the user provides their name, make sure to pass it to the tool. "
    "Do not engage in any other conversation or tasks.",
    description="Handles simple greetings and hellos using the 'say_hello' tool",
    tools=[say_hello],
)

farewell_agent = Agent(
    model=AGENT_MODEL,
    name="farewell_agent",
    instruction="YOu are the Farewell Agent. Your ONLY task is to provide a polite goodbye message. "
    "Use the 'say_goodbye' tool when the user indicates they are leaving or ending the conversation "
    "(e.g., using words like 'bye', 'goodbye', 'thanks bye', 'see you'). "
    "Do not perform any other actions.",
    description="Handles simple farewells and goodbyes using the 'say_goodbye' tool.",
    tools=[say_goodbye],
)

# weather_agent_team = Agent(
#     name="weather_agent_v2",
#     model=AGENT_MODEL,
#     description="The main coordinator agent. Handles weather requests and delegates greetings/farewells to specialists.",
#     instruction="You are the main Weather Agent coordinating a team. Your primary responsibility is to provide weather information. "
#     "Use the 'get_weather' tool ONLY for specific weather reqeusts (e.g., 'weather in London'). "
#     "You have specialized sub-agents: "
#     "1. 'greeting_agent': Handles simple greetings like 'Hi', 'Hello'. Delegate to it for these. "
#     "2. 'farewell_agent': Handles simple farewells like 'Bye', 'See you'. Delegate to it for these. "
#     "Analyze the user's query. If it's a greeting, delegate to 'greeting_agent'. If it's a farewell, delegate to 'farewell_agent'. "
#     "If it's a weather request, handle it yourself using 'get_weather'. "
#     "For anything else, response appropriately or state you cannot handle it.",
#     tools=[get_weather],
#     sub_agents=[greeting_agent, farewell_agent],
# )
#
session_service_stateful = InMemorySessionService()

APP_NAME = "weather_tutorial_app"
SESSION_ID_STATEFUL = "session_state_demo_001"
USER_ID_STATEFUL = "user_state_demo"

initial_state = {"user_preference_temperature_unit": "Fahrenheit"}

session_stateful = session_service_stateful.create_session(
    app_name=APP_NAME,
    user_id=USER_ID_STATEFUL,
    session_id=SESSION_ID_STATEFUL,
    state=initial_state,
)

from google.adk.tools.tool_context import ToolContext


def get_weather_stateful(city: str, tool_context: ToolContext) -> dict:
    """Retrieves weather, converts temp unit based on session state."""
    preferred_unit = tool_context.state.get(
        "user_preference_temperature_unit", "Fahrenheit"
    )

    city_normalized = city.lower().replace(" ", "")
    mock_weather_db = {
        "newyork": {"temp_c": 25, "condition": "sunny"},
        "london": {"temp_c": 15, "condition": "cloudy"},
        "tokyo": {"temp_c": 18, "condition": "light rain"},
    }

    if city_normalized in mock_weather_db:
        data = mock_weather_db[city_normalized]
        temp_c = data["temp_c"]
        condition = data["condition"]

        temp_value = temp_c
        temp_unit = "°C"

        if preferred_unit == "Fahrenheit":
            temp_value = (temp_c * 9 / 5) + 32
            temp_unit = "°F"

        report = f"The weather in {city.capitalize()} is {condition} with a temperature  of {temp_value:.0f}{temp_unit}."
        result = {"status": "success", "report": report}

        tool_context.state["last_city_checked_stateful"] = city

        return result


root_agent_stateful = Agent(
    name="weather_agent_v4_stateful",
    model=AGENT_MODEL,
    description="Main agent: Provides weather (state-aware unit), delegates greetings / farewells, saves report to state.",
    instruction="You are the main Weather Agent. Your job is to provide weather using 'get_weather_stateful'. "
    "The tool will format the temperature based on user preference stored in state. "
    "Delegate simple greetings to 'greeting_agent' and farewells to 'farewell_agent'. "
    "Handle only weather requests, greetings and farewells.",
    tools=[get_weather_stateful],
    sub_agents=[greeting_agent, farewell_agent],
    output_key="last_weather_report",
)

root_agent = root_agent_stateful

root_runner = Runner(
    agent=root_agent,
    app_name=APP_NAME,
    session_service=session_service_stateful,
)
