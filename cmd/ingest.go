package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/spf13/cobra"
	"local-rag/internal/ingest"
)

var bookPath string
var topic	string

var ingestCmd = &cobra.Command{
    Use: "ingest",
    Short: "Ingest files (text or otherwise) into the knowledge base",
    RunE: func(cmd *cobra.Command, args []string) error {
        fmt.Printf("Ingesting books within %s\n", bookPath)
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
			content, err := ingestor.ExtractText(fullPath)
			if err != nil {
				return err
			}

			fmt.Printf("Successfully read %d characters from %s\n", len(content), fileName)

		}

		return nil
    },
}

func init() {
	
	rootCmd.AddCommand(ingestCmd)
	ingestCmd.Flags().StringVarP(&bookPath, "path", "p", "", "Path to the books directory")
	ingestCmd.Flags().StringVarP(&topic, "topic", "t", "", "Topic covered by books in the directory")	
}