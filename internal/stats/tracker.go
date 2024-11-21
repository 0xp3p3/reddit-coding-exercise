package stats

import (
	"log"
	"sort"
	"sync"

	"subreddit-exercise/internal/models"
)

type Tracker struct {
	mu        sync.Mutex
	UserPosts map[string]int
	TopPosts  []models.Post
}

func NewTracker() *Tracker {
	return &Tracker{
		UserPosts: make(map[string]int),
		TopPosts:  []models.Post{},
	}
}

func (t *Tracker) SortPosts() {
	sort.SliceStable(t.TopPosts, func(i, j int) bool {
		return t.TopPosts[i].Upvotes > t.TopPosts[j].Upvotes
	})
}

func (t *Tracker) AddPost(post models.Post) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Track posts by user
	t.UserPosts[post.Author]++

	// Duplicate post check before add
	for _, p := range t.TopPosts {
		if p.Title == post.Title && p.Author == post.Author {
			if post.Upvotes > p.Upvotes {
				p.Upvotes = post.Upvotes
				t.SortPosts()
			}
			return
		}
	}

	// Track top posts
	t.TopPosts = append(t.TopPosts, post)

	t.SortPosts()
}

func (t *Tracker) LogStats() {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Sort by Posts
	type kv struct {
		Key   string
		Value int
	}
	var sorted []kv
	for k, v := range t.UserPosts {
		sorted = append(sorted, kv{k, v})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value > sorted[j].Value // Descending order
	})

	log.Println("Top Users by Posts:")
	for _, kv := range sorted {
		log.Printf("%s: %d posts\n", kv.Key, kv.Value)
	}

	log.Println("Top Posts by Upvotes:")
	sortPosts := t.TopPosts
	if len(sortPosts) > 10 {
		sortPosts = sortPosts[:10]
	}

	for _, post := range sortPosts {
		log.Printf("  %s by %s (%d upvotes)\n", post.Title, post.Author, post.Upvotes)
	}
}
