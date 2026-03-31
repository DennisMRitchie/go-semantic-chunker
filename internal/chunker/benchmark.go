package chunker

import (
	"fmt"
	"strings"
	"time"
)

// BenchmarkResult holds timing and quality metrics for a chunking run.
type BenchmarkResult struct {
	Strategy      string        `json:"strategy"`
	ChunksCount   int           `json:"chunks_count"`
	Duration      time.Duration `json:"duration_ns"`
	DurationHuman string        `json:"duration_human"`
	AvgChunkSize  int           `json:"avg_chunk_size_chars"`
	MaxChunkSize  int           `json:"max_chunk_size_chars"`
	MinChunkSize  int           `json:"min_chunk_size_chars"`
}

// RunBenchmark compares naive, recursive, and semantic chunking on the same text.
func RunBenchmark(text string, opts ChunkOptions) []BenchmarkResult {
	return []BenchmarkResult{
		benchmarkNaive(text, opts),
		benchmarkRecursive(text, opts),
		benchmarkSemantic(text, opts),
	}
}

func benchmarkNaive(text string, opts ChunkOptions) BenchmarkResult {
	start := time.Now()
	var chunks []Chunk
	for i := 0; i < len(text); i += opts.MaxChunkSize {
		end := i + opts.MaxChunkSize
		if end > len(text) {
			end = len(text)
		}
		chunks = append(chunks, Chunk{
			ID:      fmt.Sprintf("naive-%04d", i),
			Content: text[i:end],
		})
	}
	return buildResult("naive", chunks, time.Since(start))
}

func benchmarkRecursive(text string, opts ChunkOptions) BenchmarkResult {
	start := time.Now()
	chunks := RecursiveSplit(text, opts)
	return buildResult("recursive", chunks, time.Since(start))
}

func benchmarkSemantic(text string, opts ChunkOptions) BenchmarkResult {
	o := opts
	o.UseSemantic = true
	start := time.Now()
	chunks := SemanticSplit(text, o)
	return buildResult("semantic", chunks, time.Since(start))
}

func buildResult(strategy string, chunks []Chunk, d time.Duration) BenchmarkResult {
	if len(chunks) == 0 {
		return BenchmarkResult{Strategy: strategy, Duration: d, DurationHuman: d.String()}
	}
	total, maxS, minS := 0, 0, int(^uint(0)>>1)
	for _, c := range chunks {
		l := len(strings.TrimSpace(c.Content))
		total += l
		if l > maxS {
			maxS = l
		}
		if l < minS {
			minS = l
		}
	}
	return BenchmarkResult{
		Strategy:      strategy,
		ChunksCount:   len(chunks),
		Duration:      d,
		DurationHuman: d.String(),
		AvgChunkSize:  total / len(chunks),
		MaxChunkSize:  maxS,
		MinChunkSize:  minS,
	}
}
