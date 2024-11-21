package reddit

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_FetchPosts(t *testing.T) {
	// Mock Reddit API response
	mockResponse := `{
		"data": {
			"children": [
				{"data": {"id": "1", "title": "Post 1", "author": "User1", "score": 10}},
				{"data": {"id": "2", "title": "Post 2", "author": "User2", "score": 20}}
			]
		}
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	client := &Client{
		HTTPClient: http.DefaultClient,
		BaseURL:    server.URL,
	}

	posts, err := client.FetchPosts("test")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(posts) != 2 {
		t.Errorf("Expected 2 posts, got %d", len(posts))
	}

	if posts[0].Title != "Post 1" {
		t.Errorf("Expected first post to be 'Post 1', got '%s'", posts[0].Title)
	}
}
