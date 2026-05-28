package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"local-rag/internal/ingest"
	"local-rag/internal/utils"
	"net/http"
	"sync"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

type EmbeddingRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type EmbeddingResponse struct {
	Embedding []float64 `json:"embedding"`
}

type GenerateRequest struct {
	Model      string `json:"model"`
	Prompt     string `json:"prompt"`
	StreamMode bool   `json:"stream"`
}

type GenerateResponse struct {
	Response string `json:"response"`
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}
}

func (c *Client) BatchEmbed(ctx context.Context, model string, chunks []ingest.Chunk, concurrency int) error {
	jobs := make(chan int, len(chunks))
	errs := make(chan error, len(chunks))
	ctx, cancel := context.WithCancel(ctx)

	defer cancel()

	var wg sync.WaitGroup

	for w := 1; w <= concurrency; w++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for index := range jobs {
				select {
				case <-ctx.Done():
					return
				default:
				}
				vector, err := c.GetEmbedding(ctx, model, chunks[index].Content)
				if err != nil {
					errs <- err
					cancel()
					return
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
	select {
	case err := <-errs:
		return err
	default:
	}
	
	return nil
}

func (c *Client) GetEmbedding(ctx context.Context, model string, prompt string) ([]float64, error) {
	reqBody := EmbeddingRequest{
		Model:  model,
		Prompt: prompt,
	}

	encodedJsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	postRequest, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL+"/api/embeddings", bytes.NewReader(encodedJsonBody))
	if err != nil {
		return nil, err
	}
	postRequest.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(postRequest);
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

func (c *Client) Generate(ctx context.Context, model string, prompt string) (string, error) {

	reqBody := GenerateRequest{
		Model:      model,
		Prompt:     prompt,
		StreamMode: false,
	}

	encodedJsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	postRequest, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL+"/api/generate", bytes.NewReader(encodedJsonBody))
	if err != nil {
		return "", err
	}
	postRequest.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(postRequest);
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
