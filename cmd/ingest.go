package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/spf13/cobra"
	"github.com/ledongthuc/pdf"
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

		for _, file := range dirEntries {
			fileName := file.Name()
			fileExt := filepath.Ext(fileName)
			if fileExt != ".pdf" {
				continue
			}
			fmt.Println("Reading", fileName)

			fullPath := bookPath + "/" + fileName
			f, err := os.Open(fullPath)
			if err != nil {
				return err
			}
			defer f.Close()

			fileMeta, err := f.Stat()
			if err != nil {
				return err
			}
			
			reader, err := pdf.NewReader(f, fileMeta.Size())
			if err != nil {
				return err
			}

			page := reader.Page(1)
			if page.V.IsNull() {
				return fmt.Errorf("Page is null")
			}

			content, err := page.GetPlainText(nil)
			if err != nil {
				return err
			}

			endCharacterCount := min(100, len(content))
			fmt.Println(fileName, content[:endCharacterCount])


		}

		

		return nil
    },
}

func init() {
	
	rootCmd.AddCommand(ingestCmd)
	ingestCmd.Flags().StringVarP(&bookPath, "path", "p", "", "Path to the books directory")
	ingestCmd.Flags().StringVarP(&topic, "topic", "t", "", "Topic covered by books in the directory")	
}