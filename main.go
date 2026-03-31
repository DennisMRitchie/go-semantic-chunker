package main

import (
	"log"
	"net/http"
	"os"

	"github.com/DennisMRitchie/go-semantic-chunker/internal/chunker"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "go-semantic-chunker"})
	})

	// POST /chunk — main chunking endpoint
	r.POST("/chunk", func(c *gin.Context) {
		var req struct {
			Text string `json:"text" binding:"required"`
			chunker.ChunkOptions
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Apply sensible defaults
		if req.MaxChunkSize == 0 {
			req.MaxChunkSize = 800
		}
		if req.Overlap == 0 {
			req.Overlap = 100
		}
		if req.SimilarityThreshold == 0 {
			req.SimilarityThreshold = 0.85
		}

		var chunks []chunker.Chunk
		if req.UseSemantic {
			chunks = chunker.SemanticSplit(req.Text, req.ChunkOptions)
		} else {
			chunks = chunker.RecursiveSplit(req.Text, req.ChunkOptions)
		}

		c.JSON(http.StatusOK, gin.H{
			"strategy":     map[bool]string{true: "semantic", false: "recursive"}[req.UseSemantic],
			"chunks_count": len(chunks),
			"chunks":       chunks,
		})
	})

	// POST /benchmark — compare all three strategies
	r.POST("/benchmark", func(c *gin.Context) {
		var req struct {
			Text string `json:"text" binding:"required"`
			chunker.ChunkOptions
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if req.MaxChunkSize == 0 {
			req.MaxChunkSize = 800
		}
		if req.Overlap == 0 {
			req.Overlap = 100
		}
		if req.SimilarityThreshold == 0 {
			req.SimilarityThreshold = 0.85
		}

		results := chunker.RunBenchmark(req.Text, req.ChunkOptions)
		c.JSON(http.StatusOK, gin.H{"benchmark": results})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	log.Printf("🚀 go-semantic-chunker running on :%s", port)
	log.Fatal(r.Run(":" + port))
}
