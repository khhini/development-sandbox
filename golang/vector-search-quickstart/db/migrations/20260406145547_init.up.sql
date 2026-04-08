-- Add up migration script here

CREATE EXTENSION IF NOT EXISTS vector;

CREATE TABLE indonesian_foods (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    island_origin TEXT,
    embedding vector(768) -- Matches Google Gemini dimensions
);

CREATE INDEX ON indonesian_foods USING hnsw (embedding vector_cosine_ops);
