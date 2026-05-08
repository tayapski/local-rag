package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"local-rag/internal/ingest"
	"local-rag/internal/utils"
	"net/http"
	"sync"
)

type Client struct {
	BaseURL 	string
	HTTPClient *http.Client
}

type EmbeddingRequest struct {
	Model 	string `json:"model"`
	Prompt	string `json:"prompt"`
}

type EmbeddingResponse struct {
	Embedding []float64 `json:"embedding"`
}

type GenerateRequest struct {
	Model		string 	`json:"model"`
	Prompt		string	`json:"prompt"`
	StreamMode	bool	`json:"stream"`

}

type GenerateResponse struct {
	Response 	string	`json:"response"`
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{},
	}
}

func (c *Client) BatchEmbed(model string, chunks []ingest.Chunk, concurrency int) error {
	jobs := make(chan int, len(chunks))
	var wg sync.WaitGroup

	for w := 1; w<= concurrency; w++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for index := range jobs {
				vector, err := c.GetEmbedding(model, chunks[index].Content)
				if err != nil {
					fmt.Printf("Error embedding chunk %d: %v\n", index, err)
					continue
				}

				chunks[index].Embedding = utils.ConvertSlice(vector)
			}
		}()
	}

	for i := range chunks {
		jobs <- i
	}

	close(jobs)

	wg.Wait()
	return nil
}

func (c *Client) GetEmbedding(model string, prompt string) ([]float64, error) {
	reqBody := EmbeddingRequest{
		Model: model,
		Prompt: prompt,
	}

	encodedJsonBody, err := json.Marshal(reqBody);
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Post(c.BaseURL+"/api/embeddings", "application/json", bytes.NewReader(encodedJsonBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()


	var response EmbeddingResponse
	responseCode := resp.StatusCode
	if responseCode != http.StatusOK {
		return nil, fmt.Errorf("ollama returned status: %d", responseCode)
	}

	responseErr := json.NewDecoder(resp.Body).Decode(&response)
	if responseErr != nil {
		return nil, responseErr
	}

	return response.Embedding, nil

}

func (c *Client) Generate(model string, prompt string) (string, error){

	reqBody := GenerateRequest{
		Model: model,
		Prompt: prompt,
		StreamMode: false,
	}

	encodedJsonBody, err := json.Marshal(reqBody);
	if err != nil {
		return "", err
	}

	resp, err := c.HTTPClient.Post(c.BaseURL+"/api/generate", "application/json", bytes.NewReader(encodedJsonBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()


	var response GenerateResponse
	responseCode := resp.StatusCode
	if responseCode != http.StatusOK {
		return "", fmt.Errorf("ollama returned status: %d", responseCode)
	}

	responseErr := json.NewDecoder(resp.Body).Decode(&response)
	if responseErr != nil {
		return "", responseErr
	}

	return response.Response, nil


}
