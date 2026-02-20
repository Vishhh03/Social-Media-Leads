---
name: PostgreSQL pgvector Configuration
description: Upgrading standard postgres to pgvector for AI RAG capabilities.
---
# PostgreSQL pgvector Setup

## Key Findings & Caveats
1. **Docker Replacement:** You do not need a separate database container for vector databases. Replace the standard `postgres:16-alpine` image with the `pgvector/pgvector:pg16` image in `docker-compose.yml`. It is fully backward compatible with your standard Postgres configuration, passwords, ports, and volumes.
2. **Enable the Extension:** Before using `vector(1536)` types in tables or running creation schemas, you must run the DB migration directive natively: `CREATE EXTENSION IF NOT EXISTS vector;`.
3. **Cosine Similarity Magic:** The standard `pgx` driver in Go handles converting `[]float32` slices into Postgres vectors. To query the most relevant documents in Postgres natively for Retrieval-Augmented Generation, use the `<=>` Cosine Similarity operator natively supported by the extension and order by it:
   `ORDER BY embedding <=> $1 LIMIT 5`.
