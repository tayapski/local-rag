package cmd

import (
	"fmt"
	"local-rag/internal/ai"
	"local-rag/internal/config"
	"local-rag/internal/db"
	"local-rag/internal/store"
	"local-rag/internal/utils"
	"strings"

	"github.com/spf13/cobra"
)

var askCmd = &cobra.Command{
	Use:   "ask",
	Short: "Ask questions about a topic",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmdCtx := cmd.Context()
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

	
		fmt.Println("Processing query")
		prompt := args[0]
		promptVector, err := aiClient.GetEmbedding(cmdCtx, "nomic-embed-text", prompt)
		if err != nil {
			return err
		}

		fmt.Println("Building context")
		searchResults, err := storeClient.Search(
			cmdCtx, topic, utils.ConvertSlice(promptVector),
		)
		if err != nil {
			return err
		}

		var promptBuilder strings.Builder
		promptBuilder.WriteString("Answer the question based on the provided context below\n\n")
		promptBuilder.WriteString("Context:\n")

		for _, result := range searchResults {
			source, err := db.GetSource(cmdCtx, result.SourceID)
			if err != nil {
				return err
			}
			citation := fmt.Sprintf("[Source: %s, Page: %d]", source.Metadata["title"], result.PageNumber)

			promptBuilder.WriteString(fmt.Sprintf("%s\n%s\n\n", citation, result.Content))
			promptBuilder.WriteString("\n")
		}

		promptBuilder.WriteString("\nQuestion: ")
		promptBuilder.WriteString(prompt)
		promptBuilder.WriteString("\nAnswer: ")

		fmt.Println("Context and prompt completed")
		answer, err := aiClient.Generate(cmdCtx, "llama3", promptBuilder.String())
		if err != nil {
			return err
		}
		fmt.Printf("%s", answer)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(askCmd)
	askCmd.Flags().StringVarP(&topic, "topic", "t", "", "Topic covered by the question or prompt")
}
