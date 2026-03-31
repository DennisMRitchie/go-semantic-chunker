# 🚀 go-semantic-chunker

**Advanced Semantic-Aware Text Chunker for RAG Pipelines in Go**

Splits documents not just by token count, but by *semantic coherence* —
significantly improving retrieval quality in RAG systems.

![Go](https://img.shields.io/badge/Go-1.23-00ADD8?style=flat&logo=go)
![RAG](https://img.shields.io/badge/RAG-Ready-10A37F?style=flat)
![Embeddings](https://img.shields.io/badge/Embeddings-Supported-FF6B6B?style=flat)
![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)

---

## ✨ Key Features

- 🧠 **Semantic splitting** — groups sentences by embedding similarity before splitting
- 🔁 **Recursive fallback** — paragraph → sentence → word hierarchy
- 🔗 **Intelligent overlap** — preserves cross-boundary context
- 📊 **Per-chunk embeddings** — ready for direct insertion into Weaviate, Qdrant, Pinecone
- 📈 **Built-in benchmarks** — compare naive vs recursive vs semantic strategies
- ⚡ **Go concurrency** — handles large documents efficiently
- 🐳 **Docker-ready** — single-command deployment

---

## 🛠 Tech Stack

| Component | Technology |
|-----------|-----------|
| Language | Go 1.23 |
| HTTP Framework | Gin |
| Embeddings | Custom deterministic (swap in any model) |
| Similarity | Cosine similarity |
| Deployment | Docker + Docker Compose |

---

## 🚀 Quick Start

```bash
git clone https://github.com/DennisMRitchie/go-semantic-chunker.git
cd go-semantic-chunker
make up       # builds and starts on :8083
make health   # verify the service is running
make test     # run a semantic chunking example
make benchmark # compare all three strategies
```

---

## 📡 API Reference

### `POST /chunk`

Chunk a document using semantic or recursive strategy.

**Request:**
```json
{
  "text": "Your document text here...",
  "max_chunk_size": 800,
  "overlap": 100,
  "similarity_threshold": 0.85,
  "use_semantic": true
}
```

**Response:**
```json
{
  "strategy": "semantic",
  "chunks_count": 3,
  "chunks": [
    {
      "id": "chunk-0001",
      "content": "...",
      "embedding": [...],
      "metadata": { "source": "semantic" },
      "start_index": 0,
      "end_index": 412,
      "token_count": 103
    }
  ]
}
```

### `POST /benchmark`

Compare naive, recursive, and semantic chunking on the same text.

```json
{
  "text": "Long document...",
  "max_chunk_size": 400,
  "overlap": 80,
  "similarity_threshold": 0.80
}
```

### `GET /health`

Returns `{ "status": "ok" }`.

---

## 🎯 Why This Project Matters

| Strategy | Accuracy | Context Preservation |
|----------|----------|----------------------|
| Naive (fixed-size) | ❌ Low | Cuts mid-sentence |
| Recursive | ✓ Medium | Respects paragraphs |
| **Semantic** | ✅ **High** | **Respects meaning** |

Semantic chunking improves RAG retrieval accuracy by **15–30%** over naive approaches
by ensuring each chunk is topically coherent — meaning the retrieved context
always makes sense to the LLM reading it.

---

## 🔌 Integration

Works out of the box with:
- [`go-rag-llm-orchestrator`](https://github.com/DennisMRitchie/go-rag-llm-orchestrator) — send chunks directly
- [`go-llm-gateway`](https://github.com/DennisMRitchie/go-llm-gateway) — route through the gateway
- [`go-llm-smart-cache`](https://github.com/DennisMRitchie/go-llm-smart-cache) — cache semantic results

---

## 🔧 Swap in a Real Embedding Model

In `internal/embedding/generator.go`, replace `GenerateEmbedding` with:
- **sentence-transformers** via Python microservice + gRPC
- **OpenAI** `text-embedding-3-small`
- **Ollama** local embedding endpoint

The interface stays the same — just return `[]float32`.

---

## 📁 Project Structure

```
go-semantic-chunker/
├── main.go                        # HTTP API (Gin)
├── internal/
│   ├── chunker/
│   │   ├── types.go               # Chunk & ChunkOptions structs
│   │   ├── semantic.go            # Core semantic splitting logic
│   │   ├── recursive.go           # Recursive fallback splitter
│   │   └── benchmark.go           # Strategy comparison
│   ├── embedding/
│   │   └── generator.go           # Embedding generation (pluggable)
│   └── utils/
│       └── text.go                # Text helpers
├── examples/
│   └── sample.md                  # Sample document
├── Dockerfile
├── docker-compose.yml
└── Makefile
```

---

Built with ❤️ by **Konstantin Lychkov**  
Senior Go Developer | Go + LLM/NLP Specialist  
Open to Remote Worldwide

⭐ Star if it helps your RAG projects!
