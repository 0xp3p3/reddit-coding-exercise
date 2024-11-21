package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"subreddit-exercise/internal/models"
	"subreddit-exercise/internal/reddit"
	"subreddit-exercise/internal/stats"
)

func runApp(ctx context.Context, fetcher reddit.Fetcher, tracker *stats.Tracker, subreddits []string) {
	var wg sync.WaitGroup
	postsChannel := make(chan models.Post, 100)

	// Goroutines for fetching posts from subreddits
	for _, subreddit := range subreddits {
		wg.Add(1)
		go func(sub string) {
			defer wg.Done()
			ticker := time.NewTicker(10 * time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					log.Printf("Stopping fetcher for subreddit: %s", sub)
					return
				case <-ticker.C:
					posts, err := fetcher.FetchPosts(sub)
					if err != nil {
						log.Printf("Error fetching posts from %s: %v", sub, err)
						continue
					}
					for _, post := range posts {
						postsChannel <- post
					}
				}
			}
		}(subreddit)
	}

	// Worker pool for processing posts
	numWorkers := 4
	for i := 0; i < numWorkers; i++ {
		go func(workerID int) {
			log.Printf("Worker %d started", workerID)
			for {
				select {
				case <-ctx.Done():
					log.Printf("Worker %d stopping", workerID)
					return
				case post := <-postsChannel:
					tracker.AddPost(post)
				}
			}
		}(i)
	}

	// Periodically log statistics
	statsTicker := time.NewTicker(30 * time.Second)
	defer statsTicker.Stop()

	go func() {
		for range statsTicker.C {
			tracker.LogStats()
		}
	}()

	// Wait for all fetchers to complete
	wg.Wait()
	close(postsChannel)
}

func main() {
	log.Println("Starting Reddit Coding Exercise Application...")

	// Context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Capture OS signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		log.Println("Shutting down application...")
		cancel()
	}()

	// Initialize Reddit client and statistics tracker
	client := reddit.NewClient()
	tracker := stats.NewTracker()

	// Define the subreddits to track
	subreddits := []string{"programming", "golang", "technology"}

	// Run the application
	runApp(ctx, client, tracker, subreddits)

	log.Println("Application stopped.")
}
