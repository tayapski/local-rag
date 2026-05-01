package ingest

import (
	"os"
	"strings"

	"github.com/ledongthuc/pdf"
)

type Ingestor struct {
	Topic string
}

func (i *Ingestor) ExtractText(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	fileMeta, err := f.Stat()
	if err != nil {
		return "", err
	}
	
	reader, err := pdf.NewReader(f, fileMeta.Size())
	if err != nil {
		return "", err
	}

	var builder strings.Builder

	for p := 1; p<= reader.NumPage(); p++ {
		page := reader.Page(p)
		content, err := page.GetPlainText(nil)
		if err != nil {
			return path, err
		}
		builder.WriteString(content)
		if err != nil {
			return path, err
		}
	}

	
	return builder.String(), nil
}


