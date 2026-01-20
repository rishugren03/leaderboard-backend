package models

// User represents a leaderboard user with their rating.
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Rating   int    `json:"rating"`
}

// LeaderboardEntry represents a user with their computed rank.
type LeaderboardEntry struct {
	Rank     int    `json:"rank"`
	Username string `json:"username"`
	Rating   int    `json:"rating"`
	ID       int    `json:"id"`
}

// SearchResult represents a search result with global rank.
type SearchResult struct {
	GlobalRank int    `json:"globalRank"`
	Username   string `json:"username"`
	Rating     int    `json:"rating"`
	ID         int    `json:"id"`
}
