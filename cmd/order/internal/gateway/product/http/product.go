package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	dto "github.com/tricong1998/go-ecom/cmd/product/pkg/dto"
)

type Gateway struct {
	addr string
}

func New(addr string) *Gateway {
	return &Gateway{addr: addr}
}

func (m *Gateway) Get(ctx context.Context, id string) (*dto.ProductResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/%s", m.addr, id), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get metadata: %s", resp.Status)
	}

	var metadata dto.ProductResponse
	err = json.NewDecoder(resp.Body).Decode(&metadata)
	if err != nil {
		return nil, err
	}

	return &metadata, nil
}
