package chunker

import (
	"fmt"
	"strings"
)

// RecursiveSplit splits text using a hierarchy of separators: paragraph → sentence → word.
// This is the fallback strategy when semantic splitting is disabled or text is too short.
func RecursiveSplit(text string, opts ChunkOptions) []Chunk {
	separators := []string{"\n\n", "\n", ". ", " "}
	return recursiveSplitWith(text, separators, opts, 0)
}

func recursiveSplitWith(text string, separators []string, opts ChunkOptions, depth int) []Chunk {
	if len(text) <= opts.MaxChunkSize || len(separators) == 0 {
		chunk := createChunk(text, depth)
		chunk.ID = fmt.Sprintf("rchunk-%04d", depth)
		return []Chunk{chunk}
	}

	sep := separators[0]
	parts := strings.Split(text, sep)

	var chunks []Chunk
	current := ""
	idx := 0

	for _, part := range parts {
		candidate := current
		if candidate != "" {
			candidate += sep
		}
		candidate += part

		if len(candidate) > opts.MaxChunkSize && current != "" {
			// Recurse on remaining separators if still too large
			if len(current) > opts.MaxChunkSize {
				sub := recursiveSplitWith(current, separators[1:], opts, idx)
				chunks = append(chunks, sub...)
			} else {
				c := createChunk(current, idx)
				c.ID = fmt.Sprintf("rchunk-%04d", idx)
				chunks = append(chunks, c)
			}
			idx++
			current = part
		} else {
			current = candidate
		}
	}

	if strings.TrimSpace(current) != "" {
		if len(current) > opts.MaxChunkSize {
			sub := recursiveSplitWith(current, separators[1:], opts, idx)
			chunks = append(chunks, sub...)
		} else {
			c := createChunk(current, idx)
			c.ID = fmt.Sprintf("rchunk-%04d", idx)
			chunks = append(chunks, c)
		}
	}

	return chunks
}
