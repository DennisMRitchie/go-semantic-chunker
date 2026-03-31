package chunker

// Chunk represents a single semantic chunk of text with optional embedding and metadata.
type Chunk struct {
	ID        string            `json:"id"`
	Content   string            `json:"content"`
	Embedding []float32         `json:"embedding,omitempty"`
	Metadata  map[string]string `json:"metadata"`
	StartIndex int              `json:"start_index"`
	EndIndex   int              `json:"end_index"`
	TokenCount int              `json:"token_count"`
}

// ChunkOptions configures chunking behavior.
type ChunkOptions struct {
	MaxChunkSize        int     `json:"max_chunk_size"`        // max chars per chunk
	Overlap             int     `json:"overlap"`               // overlap in chars
	SimilarityThreshold float32 `json:"similarity_threshold"`  // 0..1 for semantic merge
	UseSemantic         bool    `json:"use_semantic"`          // enable semantic merging
}
