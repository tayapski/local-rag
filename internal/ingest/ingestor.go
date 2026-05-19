package ingest

import (
	"context"
	"crypto/md5"
	"fmt"
	"local-rag/internal/db"
	"local-rag/internal/utils"
	"os"
	"path/filepath"
	"github.com/google/uuid"
	"github.com/ledongthuc/pdf"
)

type Ingestor struct {
	Topic string
	DB *db.MetadataDB
}

type Chunk struct {
	ID         string
	Content    string
	PageNumber int
	ChunkIndex int
	FilePath   string
	SourceID   int64
	Embedding  []float32
}

func (i *Ingestor) GenerateChunkID(filePath string, page int, chunkIndex int) string {
	fileName := filepath.Base(filePath)
	idSource := fmt.Sprintf("%s-%d-%d", fileName, page, chunkIndex)
	hash := md5.Sum([]byte(idSource))
	chunkID, _ := uuid.FromBytes(hash[:])
	return chunkID.String()
}

func (i *Ingestor) ProcessFile(ctx context.Context, path string) ([]Chunk, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	checksum, err := utils.ComputeHash(f)
	if err != nil {
		return nil, err
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	fi, _ := f.Stat()
	pdfReader, err := pdf.NewReader(f, fi.Size())

	// getting the metadata
	info := pdfReader.Trailer().Key("Info")
	metadata := make(map[string]any)
	for _, key := range info.Keys() {
		metadata[key] = info.Key(key).String()
	}

	docId, err := i.DB.SaveSource(
		ctx,
		&db.Source{
			FilePath: path,
			CheckSum: checksum,
			Metadata: metadata,
		},
	)

	if err != nil {
		return nil, err
	}
	
	return i.ExtractChunk(pdfReader, path, docId)
}

func (i *Ingestor) ExtractChunk(pReader *pdf.Reader, path string, sourceID int64) ([]Chunk, error) {

	var allChunks []Chunk

	for p := 1; p <= pReader.NumPage(); p++ {
		page := pReader.Page(p)
		content, err := page.GetPlainText(nil)
		if err != nil {
			return nil, err
		}

		runes := []rune(content)
		for idx := 0; idx < len(runes); idx += 1000 {
			end := min(idx+1000, len(runes))
			stringChunk := string(runes[idx:end])
			nextChunkIndex := len(allChunks)

			newChunk := Chunk{
				ID:         i.GenerateChunkID(path, p, nextChunkIndex),
				Content:    stringChunk,
				PageNumber: p,
				ChunkIndex: len(allChunks),
				FilePath:   path,
				SourceID: 	sourceID,
			}

			allChunks = append(allChunks, newChunk)
		}

	}

	return allChunks, nil

}
