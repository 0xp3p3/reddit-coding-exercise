package reddit

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"subreddit-exercise/internal/models"
)

type Client struct {
	HTTPClient *http.Client
	BaseURL    string
}

func NewClient() *Client {
	return &Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		BaseURL:    "https://www.reddit.com",
	}
}

// Fetcher interface for fetching posts
type Fetcher interface {
	FetchPosts(subreddit string) ([]models.Post, error)
}

func (c *Client) FetchPosts(subreddit string) ([]models.Post, error) {
	url := fmt.Sprintf("%s/r/%s/new.json", c.BaseURL, subreddit)
	fmt.Println(url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		resetTime := resp.Header.Get("X-Ratelimit-Reset")
		resetSeconds, _ := strconv.Atoi(resetTime)
		log.Printf("Rate limit exceeded. Waiting for %d seconds to reset.", resetSeconds)
		time.Sleep(time.Duration(resetSeconds) * time.Second)
		return nil, fmt.Errorf("rate limit exceeded, retrying after reset")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Adjust delay based on rate limit headers
	remaining, err := strconv.ParseFloat(resp.Header.Get("X-Ratelimit-Remaining"), 64)
	if err == nil && remaining < 5 { // Add a buffer for safety
		resetTime := resp.Header.Get("X-Ratelimit-Reset")
		resetSeconds, _ := strconv.Atoi(resetTime)
		delay := time.Duration(resetSeconds) * time.Second / time.Duration(remaining+1)
		log.Printf("Rate limiting: delaying requests by %v seconds.", delay)
		time.Sleep(delay)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data struct {
			Children []struct {
				Data models.Post `json:"data"`
			} `json:"children"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	var posts []models.Post
	for _, child := range response.Data.Children {
		posts = append(posts, child.Data)
	}

	return posts, nil
}
