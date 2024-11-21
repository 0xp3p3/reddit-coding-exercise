
package stats

import (
	"testing"

	"subreddit-exercise/internal/models"
)

func TestTracker_AddPost(t *testing.T) {
	tracker := NewTracker()

	// Add a post
	post := models.Post{ID: "1", Title: "Post 1", Author: "User1", Upvotes: 10}
	tracker.AddPost(post)

	if tracker.UserPosts["User1"] != 1 {
		t.Errorf("Expected 1 post for User1, got %d", tracker.UserPosts["User1"])
	}

	if len(tracker.TopPosts) != 1 || tracker.TopPosts[0].Title != "Post 1" {
		t.Errorf("Expected top post to be 'Post 1', got '%s'", tracker.TopPosts[0].Title)
	}
}

func TestTracker_TopPosts(t *testing.T) {
	tracker := NewTracker()

	// Add multiple posts
	tracker.AddPost(models.Post{ID: "1", Title: "Post 1", Author: "User1", Upvotes: 10})
	tracker.AddPost(models.Post{ID: "2", Title: "Post 2", Author: "User2", Upvotes: 20})
	tracker.AddPost(models.Post{ID: "3", Title: "Post 3", Author: "User3", Upvotes: 5})

	if len(tracker.TopPosts) != 3 {
		t.Errorf("Expected 3 top posts, got %d", len(tracker.TopPosts))
	}

	if tracker.TopPosts[0].Title != "Post 2" {
		t.Errorf("Expected top post to be 'Post 2', got '%s'", tracker.TopPosts[0].Title)
	}
}
