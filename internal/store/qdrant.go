package store

import (
	"context"
	"local-rag/internal/config"
	"local-rag/internal/ingest"
	"local-rag/internal/utils"
	"strconv"

	"github.com/qdrant/go-client/qdrant"
)

type Store struct {
	client *qdrant.Client
}

func NewStore() (*Store, error) {
	cfg := config.GetConfig()
	port, err := strconv.Atoi(cfg.QdrantPort)

	if err != nil {
		return nil, err
	}

	client, err := qdrant.NewClient(
		&qdrant.Config{
			Host: cfg.QdrantHost,
			Port: port,
		},
	)
	if err != nil {
		return nil, err
	}

	return &Store{
		client: client,
	}, nil
}

func (s *Store) CreateCollection(collectionName string) error {
	err := s.client.CreateCollection(
		context.Background(),
		&qdrant.CreateCollection{
			CollectionName: collectionName,
			VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
				Size: 768,
				Distance: qdrant.Distance_Cosine,
			}),
		},
	)

	if err != nil {
		return err
	}
	return nil
}

func (s *Store) UpsertChunks(collectionName string, chunks []ingest.Chunk) error {

	pointStructSlice := make([]*qdrant.PointStruct, len(chunks))

	for i, chunk := range chunks {
		
		pointStructSlice[i] = &qdrant.PointStruct{
			Id: qdrant.NewIDUUID(chunks[i].ID),
			Vectors: qdrant.NewVectors(chunk.Embedding...),
			Payload: qdrant.NewValueMap(map[string]any{
				"content": chunk.Content,
				"page": chunk.PageNumber,
				"file_path": chunk.FilePath,
			}),
		}
	}

	wait := true
	_, err := s.client.Upsert(
		context.Background(),
		&qdrant.UpsertPoints{
			CollectionName: collectionName,
			Points: pointStructSlice,
			Wait: &wait,
		},
	)

	if err != nil {
		return err
	}

	return nil

}

func (s *Store) Search(collectionName string, queryVector []float32) ([]string, error){
	searchResult, err := s.client.Query(
		context.Background(),
		&qdrant.QueryPoints{
			CollectionName: collectionName,
			Query: qdrant.NewQuery(queryVector...),
			Limit: utils.PtrUint64(uint64(5)),
			WithPayload: qdrant.NewWithPayloadEnable(true),
		},
	)

	if err != nil{
		return nil, err
	}
	
	payloadContents := make([]string, len(searchResult))
	for i := range searchResult {
		payloadContents[i] = searchResult[i].Payload["content"].GetStringValue()
	}

	return payloadContents, nil
}