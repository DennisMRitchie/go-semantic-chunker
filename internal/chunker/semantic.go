package chunker

import (
	"fmt"
	"math"
	"strings"

	"github.com/DennisMRitchie/go-semantic-chunker/internal/embedding"
)

// SemanticSplit splits text into semantically coherent chunks.
// It first breaks the text into sentences, groups them by size,
// applies overlap, then optionally merges adjacent similar chunks.
func SemanticSplit(text string, opts ChunkOptions) []Chunk {
	sentences := splitIntoSentences(text)
	if len(sentences) == 0 {
		return nil
	}

	var chunks []Chunk
	current := ""

	for i, sent := range sentences {
		if len(current)+len(sent) > opts.MaxChunkSize && current != "" {
			chunk := createChunk(current, len(chunks))
			chunks = append(chunks, chunk)
			// Start next chunk with overlap from previous sentences
			current = overlapSentences(sentences, i, opts.Overlap)
		} else {
			if current != "" {
				current += " " + sent
			} else {
				current = sent
			}
		}
	}

	// Flush the last chunk
	if strings.TrimSpace(current) != "" {
		chunks = append(chunks, createChunk(current, len(chunks)))
	}

	// Optionally merge semantically similar adjacent chunks
	if opts.UseSemantic && opts.SimilarityThreshold > 0 {
		chunks = mergeSimilar(chunks, opts.SimilarityThreshold)
	}

	// Re-assign IDs and indexes after potential merging
	offset := 0
	for i := range chunks {
		chunks[i].ID = fmt.Sprintf("chunk-%04d", i+1)
		chunks[i].StartIndex = offset
		chunks[i].EndIndex = offset + len(chunks[i].Content)
		offset = chunks[i].EndIndex + 1
	}

	return chunks
}

// splitIntoSentences splits text on ". ", "! ", "? " boundaries.
func splitIntoSentences(text string) []string {
	// Normalize sentence-ending punctuation
	text = strings.ReplaceAll(text, "! ", ". ")
	text = strings.ReplaceAll(text, "? ", ". ")
	text = strings.ReplaceAll(text, "!\n", ". ")
	text = strings.ReplaceAll(text, "?\n", ". ")

	parts := strings.Split(text, ". ")
	var sentences []string
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			sentences = append(sentences, trimmed+".")
		}
	}
	return sentences
}

// overlapSentences returns a string of sentences ending just before index i
// that together are within overlapSize characters.
func overlapSentences(sentences []string, currentIndex int, overlapSize int) string {
	if overlapSize == 0 || currentIndex == 0 {
		return ""
	}
	var buf []string
	total := 0
	// Walk backwards from currentIndex-1
	for j := currentIndex - 1; j >= 0; j-- {
		total += len(sentences[j])
		if total > overlapSize {
			break
		}
		buf = append([]string{sentences[j]}, buf...)
	}
	return strings.Join(buf, " ")
}

// createChunk builds a Chunk from raw content and a positional index.
func createChunk(content string, index int) Chunk {
	content = strings.TrimSpace(content)
	emb := embedding.GenerateEmbedding(content)
	return Chunk{
		ID:         fmt.Sprintf("chunk-%04d", index+1),
		Content:    content,
		Embedding:  emb,
		Metadata:   map[string]string{"source": "semantic"},
		TokenCount: len(content) / 4, // ~4 chars per token heuristic
	}
}

// mergeSimilar merges adjacent chunks whose embeddings exceed the similarity threshold.
func mergeSimilar(chunks []Chunk, threshold float32) []Chunk {
	if len(chunks) < 2 {
		return chunks
	}

	var result []Chunk
	current := chunks[0]

	for i := 1; i < len(chunks); i++ {
		sim := cosineSimilarity(current.Embedding, chunks[i].Embedding)
		if sim >= threshold {
			// Merge: combine content and average embeddings
			current.Content += " " + chunks[i].Content
			current.TokenCount += chunks[i].TokenCount
			current.Embedding = averageEmbeddings(current.Embedding, chunks[i].Embedding)
			current.Metadata["merged"] = "true"
		} else {
			result = append(result, current)
			current = chunks[i]
		}
	}
	result = append(result, current)
	return result
}

// cosineSimilarity computes cosine similarity between two equal-length vectors.
func cosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}
	var dot, normA, normB float64
	for i := range a {
		dot += float64(a[i]) * float64(b[i])
		normA += float64(a[i]) * float64(a[i])
		normB += float64(b[i]) * float64(b[i])
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return float32(dot / (math.Sqrt(normA) * math.Sqrt(normB)))
}

// averageEmbeddings returns the element-wise average of two embedding vectors.
func averageEmbeddings(a, b []float32) []float32 {
	res := make([]float32, len(a))
	for i := range a {
		res[i] = (a[i] + b[i]) / 2
	}
	return res
}
