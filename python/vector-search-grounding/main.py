from google import genai
from google.genai.types import (
    GenerateContentConfig,
    HttpOptions,
    Retrieval,
    Tool,
    VertexAISearch,
)

project = "ebc-cloud-dev-03"
collection = "haimeteo-new-sources_1746612352277"
data_store = "haimeteo-new-sources_1746612352277_gcs_store"
engine_id = "haimeteo-extended_1746615480468"

client = genai.Client(
  vertexai=True,
  http_options=HttpOptions(api_version="v1"),
  project=project,
  location="us-central1"
)

# Load Data Store ID from Vertex AI Search
fqdn_engine = f"projects/{project}/locations/global/collections/default_collection/engines/{engine_id}"

response = client.models.generate_content(
    model="gemini-2.0-flash-001",
    contents="What SOP of Maintenances works",
    config=GenerateContentConfig(
        tools=[
            # Use Vertex AI Search Tool
            Tool(
                retrieval=Retrieval(
                    vertex_ai_search=VertexAISearch(
                        engine=fqdn_engine,
                    )
                )
            )
        ],
    ),
)

print(response.text)
if response.candidates[0].grounding_metadata.grounding_supports:
    for s in response.candidates[0].grounding_metadata.grounding_supports:
        print(f"{s.segment.text} {s.grounding_chunk_indices}")
    
for i, chunk in enumerate(response.candidates[0].grounding_metadata.grounding_chunks):
    print(f"{i}.\n")
    print(chunk.retrieved_context.text)
    print(chunk.retrieved_context.uri)
    print("\n")
     
# Example response:
# 'The process for making an appointment to renew your driver's license varies depending on your location. To provide you with the most accurate instructions...'