package embedding

// GenerateEmbedding produces a deterministic pseudo-embedding for a given text.
// In production: replace with a call to sentence-transformers, OpenAI embeddings,
// or a dedicated gRPC embedding service.
func GenerateEmbedding(text string) []float32 {
	const dim = 384
	emb := make([]float32, dim)

	seed := 0
	for _, r := range text {
		seed = (seed*31 + int(r)) % 10007
	}
	for i := range emb {
		emb[i] = float32(seed%1000)/1000.0 + float32(i%10)/50.0
	}
	return emb
}
