package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/vlm326/golang-course-hw-78/internal/collector/domain"
)

const githubAPIBaseURL = "https://api.github.com"

type Client struct {
	httpClient *http.Client
	baseURL    string
}

type repositoryResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Stars       int64  `json:"stargazers_count"`
	Forks       int64  `json:"forks_count"`
	CreatedAt   string `json:"created_at"`
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}

	return &Client{
		httpClient: httpClient,
		baseURL:    githubAPIBaseURL,
	}
}

func (c *Client) GetRepository(ctx context.Context, owner, repo string) (domain.Repository, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/repos/%s/%s", c.baseURL, owner, repo), nil)
	if err != nil {
		return domain.Repository{}, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "golang-course-hw-78-collector")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return domain.Repository{}, fmt.Errorf("request github api: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return domain.Repository{}, domain.ErrRepositoryNotFound
	default:
		return domain.Repository{}, fmt.Errorf("github api returned status %d", resp.StatusCode)
	}

	var payload repositoryResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return domain.Repository{}, fmt.Errorf("decode github response: %w", err)
	}

	if payload.Name == "" {
		return domain.Repository{}, errors.New("github response is missing repository name")
	}

	return domain.Repository{
		Name:        payload.Name,
		Description: payload.Description,
		Stars:       payload.Stars,
		Forks:       payload.Forks,
		CreatedAt:   payload.CreatedAt,
	}, nil
}
