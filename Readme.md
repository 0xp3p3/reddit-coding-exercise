
# Reddit Coding Exercise

## Overview

This application fetches posts from one or more subreddits in near real-time and tracks statistics such as:
- **Top posts by upvotes**.
- **Users with the most posts**.

It is designed to process data as concurrently as possible using **Goroutines** and **channels**, while respecting Reddit's rate limits.

---

## Features

1. **Multi-Subreddit Tracking**:
   - Fetches posts from multiple subreddits concurrently.
   - Each subreddit runs in its own Goroutine.

2. **Concurrency**:
   - A worker pool processes posts in parallel to maximize throughput.

3. **Rate Limiting**:
   - Adapts to Reddit's rate limits dynamically based on API headers (`X-Ratelimit-*`).
   - Prevents excessive API calls while maintaining high request rates.

4. **Statistics**:
   - Tracks:
     - Top users by the number of posts.
     - Top posts by upvotes.
   - Logs aggregated statistics periodically.

5. **Thread-Safe Design**:
   - Uses synchronization (e.g., `sync.Mutex`) to handle concurrent updates to shared data.

---

## Requirements

- **Go 1.19** or later
- Internet connection

---

## Setup Instructions

### 1. Clone the Repository
```bash
git clone https://github.com/0xp3p3/reddit-coding-exercise.git
cd reddit-coding-exercise
```

### 2. Run the Application
No API token is required since the app uses Reddit's unauthenticated public API.

Simply run:
```bash
go run main.go
```

Run Tests:
```bash
go test ./...
```
---

## Configuration

The application allows you to customize subreddits and other settings by editing `main.go`:

### Set Subreddits
In `main.go`, define the subreddits you want to track:
```go
subreddits := []string{"programming", "golang", "technology"} // Add more as needed
```

---

## How It Works

### Architecture

1. **Subreddit Fetchers**:
   - Each subreddit runs in its own Goroutine, fetching posts every 10 seconds.
2. **Post Processing**:
   - Posts from all subreddits are sent to a shared `postsChannel`.
   - A worker pool processes these posts in parallel, updating statistics.
3. **Rate Limit Management**:
   - Dynamically adjusts API call frequency based on Reddit's `X-Ratelimit-*` headers.
4. **Statistics Reporting**:
   - Aggregated statistics are logged every 30 seconds.

---

## Example Output

```
Starting Reddit Coding Exercise Application...

Worker 0 started
Worker 1 started
Worker 2 started
Worker 3 started

https://www.reddit.com/r/technology/new.json
https://www.reddit.com/r/programming/new.json
https://www.reddit.com/r/golang/new.json

Top Users by Posts:
a_Ninja_b0y: 5 posts
ketralnis: 5 posts
MetaKnowing: 3 posts
stackoverflooooooow: 2 posts

Top Posts by Upvotes:

Joe Biden Just Trump-Proofed His Hallmark CHIPS Act by Mynameis__--__ (23936 upvotes)
OpenAI accidentally deleted potential evidence in NY Times copyright lawsuit by likwitsnake (1720 upvotes)
Boeing CEO to Employees: We Can’t Afford Another Mistake by uhhhwhatok (1469 upvotes)
The Right’s Triumph Over Social Media by Majano57 (498 upvotes)

Shutting down application...
Worker 0 stopping
Worker 1 stopping
Worker 2 stopping
Worker 3 stopping

Stopping fetcher for subreddit: programming
Stopping fetcher for subreddit: technology
Stopping fetcher for subreddit: golang

Application stopped.
```

---

## Future Enhancements

1. **Persistent Storage**:
   - Store fetched data in a database (e.g., PostgreSQL) for long-term analysis.

2. **Support for Authentication**:
   - Use Reddit's OAuth2 for accessing private or more detailed data.

3. **API or Web UI**:
   - Serve statistics via a REST API or a user-friendly web interface.

4. **Improved Rate Limit Handling**:
   - Implement exponential backoff for more efficient retries.

---

## Design Principles

The application adheres to:
- **SOLID principles**:
  - Modular and maintainable code.
- **Concurrency best practices**:
  - Efficient use of Goroutines and channels.
  - Thread-safe data handling.
- **Scalability**:
  - Designed to handle multiple subreddits without significant changes.
