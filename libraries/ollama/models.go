package ollama

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type LocalModels struct {
	Models []*RunningModel `json:"models"`
}

type RunningModel struct {
	Name    string               `json:"name"`
	Model   string               `json:"model"`
	Size    int64                `json:"size"`
	Digest  string               `json:"digest"`
	Details *RunningModelDetails `json:"details"`
}

type RunningModelDetails struct {
	Format            string `json:"format"`
	Family            string `json:"family"`
	ParameterSize     string `json:"parameter_size"`
	QuantizationLevel string `json:"quantization_level"`
}

func (c *client) ListRunningModels(ctx context.Context) ([]*RunningModel, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.endpointURL+"/api/tags", nil)
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
