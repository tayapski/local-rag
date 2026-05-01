package ingest

import (
	"os"
	"github.com/ledongthuc/pdf"
)

type Ingestor struct {
	Topic string
}

type Chunk struct {
	Content string
	PageNumber int
	ChunkIndex int
	FilePath string
}

func (i *Ingestor) ExtractChunk(path string) ([]Chunk, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fileMeta, err := f.Stat()
	if err != nil {
		return nil, err
	}
	
	reader, err := pdf.NewReader(f, fileMeta.Size())
	if err != nil {
		return nil, err
	}

	var allChunks []Chunk;
	
	for p := 1; p<= reader.NumPage(); p++ {
		page := reader.Page(p)
		content, err := page.GetPlainText(nil)
		if err != nil {
			return nil, err
		}

		runes := []rune(content)
		for i := 0; i < len(runes); i += 1000 {
			end := min(i+1000, len(runes))
			stringChunk := string(runes[i:end])

			newChunk := Chunk {
				Content: stringChunk,
				PageNumber: p,
				ChunkIndex: len(allChunks),
				FilePath: path,
			}

			allChunks = append(allChunks, newChunk);
		}

	}
	
	return allChunks, nil
	
}