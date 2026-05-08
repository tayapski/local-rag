package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

var topic	string

var rootCmd = &cobra.Command{
    Use:   "local-rag",
    Short: "A local RAG CLI for Carlo's library",
}



func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}