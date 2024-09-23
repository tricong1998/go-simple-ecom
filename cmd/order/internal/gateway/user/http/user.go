package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	dto "github.com/tricong1998/go-ecom/cmd/user/pkg/dto"
)

type Gateway struct {
	addr string
}

func New(addr string) *Gateway {
	return &Gateway{addr: addr}
}

func (g *Gateway) Get(ctx context.Context, movieId string) (float64, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/ratings/%s", g.addr, movieId), nil)
	if err != nil {
		return 0, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to get rating: %s", resp.Status)
	}

	var rating float64
	err = json.NewDecoder(resp.Body).Decode(&rating)
	if err != nil {
		return 0, err
	}

	return rating, nil
}

func (g *Gateway) Create(ctx context.Context, rating *dto.UserResponse) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/ratings", g.addr), nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create rating: %s", resp.Status)
	}

	return nil
}
