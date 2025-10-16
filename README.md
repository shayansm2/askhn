# AskHN: Temporal-powered Agentic RAG over Hacker News

An Agentic Retrieval-Augmented Generation (RAG) system that explores and indexes Hacker News content into Elasticsearch, and orchestrates RAG workflows with Temporal.

---

## Problem Statement

People frequently ask questions about current tech topics that are actively discussed on Hacker News (HN). General-purpose LLMs may hallucinate or miss context from ongoing HN discussions. This project solves this by:

- Indexing Hacker News comments into an Elasticsearch Knowledge Base.
- Retrieving relevant context for user questions
- Constructing prompts with the retrieved context and sending them to the configured LLM
- Uses an agentic loop to explore and expand the knowledge base before answering

## Objective

Build an end-to-end, agentic RAG application that:

- Ingests a dataset into a knowledge base
- Implements retrieval and prompting to an LLM
- Evaluates retrieval and LLM quality
- Exposes an interface (API and UI)
<!--- Collects user feedback and supports monitoring-->


<!--### Checklist mapping-->

<!--- [x] Select dataset: Hacker News stories and comments-->
<!--- [x] Ingest data into knowledge base: Temporal workflows index HN into Elasticsearch-->
<!--- [x] Implement retrieval flow: ES search → prompt build → LLM call (Gin API)-->
<!--- [ ] Evaluate RAG performance: retrieval and LLM comparisons (plan included below)-->
<!--- [x] Create an interface: REST API and Vite/React UI-->
<!--- [ ] Collect user feedback and monitor: plan included below (Temporal UI, ES/Kibana, feedback endpoint)-->

<!------->

---

## Datasets :
[Hacker News](https://news.ycombinator.com/news) stories and comments via [official Firebase API](https://github.com/HackerNews/API?tab=readme-ov-file) and [Algolia Search](https://hn.algolia.com/).


## Tech Stack

- Language: Go 1.22+
- Orchestration: Temporal (server + UI)
- Storage/KB: Elasticsearch 8.4.x
- API: Gin
- UI: Vite + React + TypeScript (`ui/`)
- LLMs: Gemini, OpenAI, Ollama-compatible HTTP API
- Containers: Docker + docker-compose

## Architecture

- **Ingestion/Orchestration**: Temporal Workflows and Activities (`internal/temporal`) coordinate fetching HN items and indexing into Elasticsearch.
  - `IndexHackerNewsStoryWorkflow` recursively pulls a story and its comments and indexes them.
  - `RetrivalAugmentedGenerationWorkflow` and `ProsConsRagWorkflow` implement classic RAG.
  - `AgenticRAGWorkflow` runs an agent loop that can trigger new indexing when needed.
- **Knowledge Base**: Elasticsearch 8.x (`internal/elasticsearch`).
- **LLM Layer**: Pluggable LLM client (`internal/llm`).
  - Supports Gemini, OpenAI and Ollama.
  - System prompts are template-driven (`internal/llm/prompts`).
- **API**: Gin HTTP server exposing RAG endpoints (`cmd/api`, `internal/api`).
- **UI**: Vite/React app in `ui/` with chat experiences.
- **Containers**: `docker-compose.yml` spins up Temporal, Temporal UI, Elasticsearch, API, and Worker.

<!--High-level flow:-->

<!--1. User sends a query → API → Temporal workflow-->
<!--2. Workflow searches ES → builds contextual prompt → calls LLM-->
<!--3. Agentic mode may first expand the index by exploring HN → then answers-->

---

## Getting Started

### Prerequisites

- Docker and Docker Compose
- An LLM provider key (e.g., Gemini or OpenAI) or local Ollama with a model installed

### Environment

Copy and adapt the example environment file:

```bash
cp example.env .env
# Then edit .env to set:
# GEMINI_API_KEY / OPEN_AI_API_KEY
# LLM (gemini|openai|ollama) and LLM_MODEL
```

### Run with Docker Compose (recommended)

```bash
docker compose up --build
```

This starts:

- Postgres (Temporal persistence)
- Temporal server at `localhost:7233`
- Temporal UI at `http://localhost:8233`
- Elasticsearch at `http://localhost:9200`
- API at `http://localhost:8080`
- Worker process for workflows

<!--UI is developed locally (see below). CORS is configured for `http://localhost:5173`.-->

### Run services locally (Development mode)

- Install temporal, golang, yarn and vite
- Temporal (dev): `temporal server start-dev`
- Elasticsearch:

```bash
docker run -it --rm --name elasticsearch \
  -p 9200:9200 -p 9300:9300 \
  -e "discovery.type=single-node" \
  -e "xpack.security.enabled=false" \
  elasticsearch:8.4.3
```

- API: `go run cmd/api/main.go`
- Worker: `go run cmd/worker/main.go`
- UI: `cd ui; yarn dev`

---

## API Endpoints

Base URL: `http://localhost:8080`

- `GET /v1/chat?message=...` — Classic RAG, synchronous result
- `GET /v3/chat?message=...` — Agentic RAG, synchronous result

Example:

```bash
curl "http://localhost:8080/v1/chat?message=What did HN say about Kubernetes?"
```

---

## Data Ingestion

Temporal workflow `IndexHackerNewsStoryWorkflow` takes an HN item id, fetches the story, traverses comments (`kids`), and indexes documents into Elasticsearch with fields `{id, score, title, type, text}`.

Agentic RAG may trigger ingestion by:

- Calling Algolia HN Search to discover related stories
- Launching child workflows to index each story and its comments

Manual kick-off can be added via a CLI or Temporal UI signal; see `internal/temporal/workflows.go` for implementation details.

---

## Retrieval and Prompting

<!--todo vector search-->
<!--- Retrieval: `TextSearch(query, size)` boosts `title` over `text` and returns ES documents.-->
- Prompting: system prompts are generated from templates (`internal/llm/prompts/`) with injected context.
- Classic flows:
  - `RetrivalAugmentedGenerationWorkflow` builds context from top-K docs and calls `LLMActivities.Chat`.
  - `ProsConsRagWorkflow` builds a stance-aware system prompt (agree/disagree) then calls the LLM.
- Agentic flow:
  - Iteratively calls `LLMActivities.AgenticChat` to decide actions (answer vs explore), may index new data, then answers.

---

## Configuration

Environment variables (see `example.env`):

- `ELASTICSEARCH_URL`, `ELASTICSEARCH_USER`, `ELASTICSEARCH_PASS`
- `GEMINI_API_KEY`, `OPEN_AI_API_KEY`
- `LLM` (gemini|openai|ollama), `LLM_MODEL`
- `TEMPORAL_HOST`, `TEMPORALTASK_QUEUE_NAME`
- `OLLAMA_BASE_URL` (e.g., `http://host.docker.internal:11434/v1`)
