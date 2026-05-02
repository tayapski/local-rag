package cmd

import (
	"fmt"
	"local-rag/internal/ai"
	"local-rag/internal/ingest"
	"local-rag/internal/store"
	"local-rag/internal/utils"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var bookPath string
var topic	string


var ingestCmd = &cobra.Command{
    Use: "ingest",
    Short: "Ingest files (text or otherwise) into the knowledge base",
    RunE: func(cmd *cobra.Command, args []string) error {
        fmt.Printf("Ingesting books within %s\n", bookPath)
		aiClient := ai.NewClient("http://localhost:11434")
		storeClient, err := store.NewStore()
		if err != nil {
			return err
		}

		err = storeClient.CreateCollection(topic)
		if err != nil {
			if strings.Contains(err.Error(), "already exists") {
				fmt.Printf("Warning: collection %s will be reused\n", topic)
			}else{
				return err
			}
		} else {
			fmt.Printf("Created collection %s\n", topic)
		}

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

			for idx, chunk := range chunks {
				embedding, err := aiClient.GetEmbedding("nomic-embed-text", chunk.Content)
				if err != nil {
					return err
				}
				chunks[idx].Embedding = utils.ConvertSlice(embedding)
				// fmt.Printf("Chunk %d embedded, Vector size: %d\n", chunk.ChunkIndex, len(embedding))

			}

			err = storeClient.UpsertChunks(topic, chunks)
			if err != nil {
				return err
			}
			fmt.Printf("Successfully saved chunks from %s to %s collection\n", fileName, topic)

		}

		return nil
    },
}

func init() {
	
	rootCmd.AddCommand(ingestCmd)
	ingestCmd.Flags().StringVarP(&bookPath, "path", "p", "", "Path to the books directory")
	ingestCmd.Flags().StringVarP(&topic, "topic", "t", "", "Topic covered by books in the directory")	
	ingestCmd.MarkFlagRequired("topic")
	ingestCmd.MarkFlagRequired("path")
}