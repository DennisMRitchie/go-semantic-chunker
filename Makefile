.PHONY: up down build test benchmark tidy lint

up:
	docker compose up --build -d

down:
	docker compose down

build:
	go build -o bin/semantic-chunker .

tidy:
	go mod tidy

lint:
	go vet ./...

test:
	@echo "── Semantic chunking ──────────────────────────────"
	curl -s -X POST http://localhost:8083/chunk \
	  -H "Content-Type: application/json" \
	  -d '{ \
	    "text": "Go is a statically typed, compiled language. It was designed at Google. The language prioritizes simplicity and performance. Goroutines make concurrency extremely lightweight. Channels allow safe communication between goroutines. The standard library is rich and batteries-included. Go compiles to a single binary with no dependencies. Docker and Kubernetes are both written in Go.", \
	    "max_chunk_size": 300, \
	    "overlap": 60, \
	    "similarity_threshold": 0.80, \
	    "use_semantic": true \
	  }' | python3 -m json.tool

benchmark:
	@echo "── Benchmark: naive vs recursive vs semantic ──────"
	curl -s -X POST http://localhost:8083/benchmark \
	  -H "Content-Type: application/json" \
	  -d '{ \
	    "text": "Go is a statically typed, compiled language. It was designed at Google. The language prioritizes simplicity and performance. Goroutines make concurrency extremely lightweight. Channels allow safe communication between goroutines. The standard library is rich and batteries-included.", \
	    "max_chunk_size": 200, \
	    "overlap": 40, \
	    "similarity_threshold": 0.80 \
	  }' | python3 -m json.tool

health:
	curl -s http://localhost:8083/health | python3 -m json.tool
