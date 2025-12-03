import os

from google.adk.artifacts.in_memory_artifact_service import InMemoryArtifactService
from google.adk.runners import Runner
from google.adk.sessions.in_memory_session_service import InMemorySessionService
from google.adk.memory.in_memory_memory_service import InMemoryMemoryService
from a2a.types import AgentCapabilities, AgentSkill, AgentCard
from google.adk.a2a.utils.agent_to_a2a import to_a2a
from . import agent

host = os.environ.get("A2A_HOST", "localhost")
port = int(os.environ.get("A2A_PORT", 10003))
PUBLIC_URL = os.environ.get("PUBLIC_URL", f"http://{host}:{port}")
SUPPORTED_CONTENT_TYPES = ["text", "text/plain"]


class PlannerAgent:
    """An agent to help user planning a event with its desire location."""

    SUPPORTED_CONTENT_TYPES = ["text", "text/plain"]

    def __init__(self):
        self._agent = self._build_agent()
        self.runner = Runner(
            app_name=self._agent.name,
            agent=self._agent,
            artifact_service=InMemoryArtifactService(),
            session_service=InMemorySessionService(),
            memory_service=InMemoryMemoryService(),
        )

        capabilities = AgentCapabilities(streaming=True)

        skill = AgentSkill(
            id="event_planner",
            name="Event Planner",
            description="""
                This agent generates multiple fun plan lsuggestions tailored to your specified location, dates, and interests,
                all designed for a moderate budget. It delivers detailed itineraries,
                including precise venue information (name, latitude, longitude, and description), in a structured JSON format.
            """,
            tags=["event_planner"],
            examples=[
                "What are some fun events to do in San Francisco on the weekend of May 26th 2025, and something on coffee and art?"
            ],
        )

        self.aget_card = AgentCard(
            name="Event Planner Agent",
            description="""
                This agent generates multiple fun plan lsuggestions tailored to your specified location, dates, and interests,
                all designed for a moderate budget. It delivers detailed itineraries,
                including precise venue information (name, latitude, longitude, and description), in a structured JSON format.
            """,
            url=f"{PUBLIC_URL}",
            version="1.0.0",
            default_input_modes=PlannerAgent.SUPPORTED_CONTENT_TYPES,
            default_output_modes=PlannerAgent.SUPPORTED_CONTENT_TYPES,
            capabilities=capabilities,
            skills=[skill],
        )

    def _build_agent(self):
        """Builds the LLM agent for the night out planning agent."""
        return agent.root_agent


planner = PlannerAgent()

a2a_app = to_a2a(
    planner._agent,
    host=host,
    port=port,
    agent_card=planner.aget_card,
)
