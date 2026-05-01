package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var ingestCmd = &cobra.Command{
    Use: "ingest",
    Short: "Ingest files (text or otherwise) into the knowledge base",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Ingestion in progress")
    },
}

func init() {
	rootCmd.AddCommand(ingestCmd)
}