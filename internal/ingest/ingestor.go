package ingest

import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"github.com/ledongthuc/pdf"
	"github.com/google/uuid"
)

type Ingestor struct {
	Topic string
}

type Chunk struct {
	ID			string
	Content 	string
	PageNumber 	int
	ChunkIndex 	int
	FilePath 	string
	Embedding	[]float32
}

func (i *Ingestor) GenerateChunkID(filePath string, page int, chunkIndex int) string {
	fileName := filepath.Base(filePath)
	idSource := fmt.Sprintf("%s-%d-%d", fileName, page, chunkIndex)
	hash := md5.Sum([]byte(idSource))
	chunkID, _ := uuid.FromBytes(hash[:])
	return chunkID.String()
}

func (ing *Ingestor) ExtractChunk(path string) ([]Chunk, error) {
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
			nextChunkIndex := len(allChunks)

			newChunk := Chunk {
				ID: ing.GenerateChunkID(path, p, nextChunkIndex),
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