package ollama

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type LocalModels struct {
	Models []*LocalModel `json:"models"`
}

type LocalModel struct {
	Name       string             `json:"name"`
	ModifiedAt time.Time          `json:"modified_at"`
	Size       int64              `json:"size"`
	Digest     string             `json:"digest"`
	Details    *LocalModelDetails `json:"details"`
}

type LocalModelDetails struct {
	Format            string `json:"format"`
	Family            string `json:"family"`
	ParameterSize     string `json:"parameter_size"`
	QuantizationLevel string `json:"quantization_level"`
}

func (c *client) ListLocalModels(ctx context.Context) ([]*LocalModel, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.endpointURL+"/models", nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	var jsonResponse *LocalModels
	if err := json.NewDecoder(resp.Body).Decode(&jsonResponse); err != nil {
		return nil, fmt.Errorf("failed to decode models: %w", err)
	}

	return jsonResponse.Models, nil
}
