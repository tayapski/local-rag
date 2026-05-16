package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"local-rag/internal/ai"
	"local-rag/internal/config"
	"local-rag/internal/ingest"
	"local-rag/internal/store"
	"local-rag/internal/db"
	"os"
	"path/filepath"
	"strings"
)

var bookPath string

var ingestCmd = &cobra.Command{
	Use:   "ingest",
	Short: "Ingest files (text or otherwise) into the knowledge base",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmdCtx := cmd.Context()
		fmt.Printf("Ingesting books within %s\n", bookPath)
		envConfig := config.GetConfig()
		aiClient := ai.NewClient(envConfig.OllamaURL)
		storeClient, err := store.NewStore()
		if err != nil {
			return err
		}

		db, err := db.NewMetadataDB()
		if err != nil {
			return err
		}
		defer db.Close()

		err = storeClient.CreateCollection(topic)
		if err != nil {
			if strings.Contains(err.Error(), "already exists") {
				fmt.Printf("Warning: collection %s will be reused\n", topic)
			} else {
				return err
			}
		} else {
			fmt.Printf("Created collection %s\n", topic)
		}

		dirEntries, err := os.ReadDir(bookPath)

		if err != nil {
			fmt.Printf("Got an error: %v\n", err)
			return err
		}

		ingestor := ingest.Ingestor{Topic: topic, DB: db}

		for _, file := range dirEntries {
			fileName := file.Name()
			fileExt := filepath.Ext(fileName)
			if fileExt != ".pdf" {
				continue
			}
			fmt.Println("Reading", fileName)

			fullPath := bookPath + "/" + fileName

			chunks, err := ingestor.ProcessFile(cmdCtx, fullPath)
			if err != nil {
				return err
			}

			fmt.Printf("Successfully read %d chunks from %s\n", len(chunks), fileName)

			err = aiClient.BatchEmbed("nomic-embed-text", chunks, 4)
			if err != nil {
				return err
			}

			fmt.Printf("Retrieved embeddings for %s\n", fileName)
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
