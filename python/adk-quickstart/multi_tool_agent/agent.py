from typing import Dict
from google.adk.agents import Agent


def get_weather(city: str) -> Dict[str, str]:
    """Retrieves the current weather report for a specified city.

    Returns:
        dict: A dictionary containing the weather information with a 'status' key ('sucess' or 'error') and a 'report' key with the weather details if successfull, or an 'error_message' if an error occured.
    """
    if city.lower() != "new york":
        return {
            "status": "error",
            "error_message": f"Weather information for {city} is not available.",
        }

    return {
        "status": "success",
        "report": "The weather in New York is sunny with a temperature of 25 degress Celsius (77 degress Fahrenheit).",
    }


def get_current_time(city: str) -> Dict[str, str]:
    """Returns the current time in a specified city.
    Args:
        dict: A dictionary containing the current time for a specified city information with a 'status' key ('success' or 'error') and a 'report' key with the current time details in a city if successful, or an 'error_message' if an error occured.
    """
    import datetime
    from zoneinfo import ZoneInfo

    if city.lower() != "new york":
        return {
            "status": "error",
            "error_message": f"Sorry, I don't have timezone information for {city}.",
        }
    tz_identifier = "America/New_York"
    tz = ZoneInfo(tz_identifier)
    now = datetime.datetime.now(tz)

    return {
        "status": "success",
        "report": f"""The current time in {city} is {now.strftime("%Y-%m-%d %H:%M:%S %Z%z")}""",
    }


root_agent = Agent(
    name="weather_time_agent",
    model="gemini-2.5-flash",
    description="Agent to answer questions about the time and weather in a city.",
    instruction="I can answer your questions about the time an weather in a city.",
    tools=[get_weather, get_current_time],
)
