-- Add up migration script here
CREATE TABLE document_chunks (
    id SERIAL PRIMARY KEY,
    document_name TEXT,     -- e.g., 'ikn_news.pdf'
    content TEXT,           -- The specific paragraph text
    page_number INTEGER, 
    embedding vector(768)   -- Gemini embedding
);

CREATE INDEX ON document_chunks USING hnsw (embedding vector_cosine_ops);
