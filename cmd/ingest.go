package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/spf13/cobra"
	"local-rag/internal/ingest"
	"local-rag/internal/ai"
)

var bookPath string
var topic	string

var ingestCmd = &cobra.Command{
    Use: "ingest",
    Short: "Ingest files (text or otherwise) into the knowledge base",
    RunE: func(cmd *cobra.Command, args []string) error {
        fmt.Printf("Ingesting books within %s\n", bookPath)
		aiClient := ai.NewClient("http://localhost:11434")
		dirEntries, err := os.ReadDir(bookPath)

		if err != nil{
			fmt.Printf("Got an error: %v\n", err)
			return err
		}

		ingestor := ingest.Ingestor{Topic: topic}

		for _, file := range dirEntries {
			fileName := file.Name()
			fileExt := filepath.Ext(fileName)
			if fileExt != ".pdf" {
				continue
			}
			fmt.Println("Reading", fileName)

			fullPath := bookPath + "/" + fileName
			chunks, err := ingestor.ExtractChunk(fullPath)
			if err != nil {
				return err
			}

			fmt.Printf("Successfully read %d chunks from %s\n", len(chunks), fileName)

			for _, chunk := range chunks {
				embedding, err := aiClient.GetEmbedding("nomic-embed-text", chunk.Content)
				if err != nil {
					return err
				}

				fmt.Printf("Chunk %d embedded, Vector size: %d\n", chunk.ChunkIndex, len(embedding))

				if chunk.ChunkIndex == 4 {
					break
				}
			}

		}

		return nil
    },
}

func init() {
	
	rootCmd.AddCommand(ingestCmd)
	ingestCmd.Flags().StringVarP(&bookPath, "path", "p", "", "Path to the books directory")
	ingestCmd.Flags().StringVarP(&topic, "topic", "t", "", "Topic covered by books in the directory")	
}