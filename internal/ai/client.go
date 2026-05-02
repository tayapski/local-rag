package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{},
	}
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

