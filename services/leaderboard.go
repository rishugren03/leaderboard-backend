package services

import (
	"leaderboard-backend/models"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// LeaderboardService manages the in-memory leaderboard with thread-safe operations.
type LeaderboardService struct {
	mu       sync.RWMutex
	users    []models.User          // Sorted by rating descending
	userByID map[int]*models.User   // O(1) lookup by ID
	userByName map[string]*models.User // O(1) lookup by username
}

// NewLeaderboardService creates a new leaderboard service instance.
func NewLeaderboardService() *LeaderboardService {
	return &LeaderboardService{
		users:      make([]models.User, 0),
		userByID:   make(map[int]*models.User),
		userByName: make(map[string]*models.User),
	}
}

// sortUsers sorts users by rating descending. Must be called with lock held.
func (s *LeaderboardService) sortUsers() {
	sort.Slice(s.users, func(i, j int) bool {
		return s.users[i].Rating > s.users[j].Rating
	})
	// Update map pointers after sort
	for i := range s.users {
		s.userByID[s.users[i].ID] = &s.users[i]
		s.userByName[s.users[i].Username] = &s.users[i]
	}
}

// SeedUsers generates and adds n random users with ratings between minRating and maxRating.
func (s *LeaderboardService) SeedUsers(n, minRating, maxRating int) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Clear existing data
	s.users = make([]models.User, 0, n)
	s.userByID = make(map[int]*models.User, n)
	s.userByName = make(map[string]*models.User, n)

	firstNames := []string{"rahul", "priya", "amit", "neha", "vikram", "ananya", "rohan", "sneha", "arjun", "kavita", "sanjay", "meera", "karan", "pooja", "raj", "divya", "nikhil", "riya", "varun", "shweta", "aditya", "anjali", "manish", "nisha", "suresh", "deepika", "akash", "shreya", "vikas", "kritika"}
	lastNames := []string{"sharma", "kumar", "singh", "patel", "gupta", "joshi", "verma", "rao", "reddy", "iyer", "nair", "menon", "pillai", "shah", "mehta", "agarwal", "jain", "mishra", "pandey", "trivedi", "chopra", "malhotra", "kapoor", "khanna", "bhatia", "sethi", "arora", "bansal", "garg", "saxena"}

	for i := 1; i <= n; i++ {
		firstName := firstNames[rand.Intn(len(firstNames))]
		lastName := lastNames[rand.Intn(len(lastNames))]
		suffix := rand.Intn(1000)
		
		username := firstName
		if rand.Float32() > 0.3 {
			username = firstName + "_" + lastName
		}
		if suffix > 500 {
			username = username + "_" + string(rune('0'+suffix%10))
		}
		
		// Ensure unique username
		baseUsername := username
		counter := 1
		for {
			if _, exists := s.userByName[username]; !exists {
				break
			}
			username = baseUsername + "_" + strconv.Itoa(counter)
			counter++
		}

		user := models.User{
			ID:       i,
			Username: username,
			Rating:   minRating + rand.Intn(maxRating-minRating+1),
		}
		s.users = append(s.users, user)
		s.userByID[user.ID] = &s.users[len(s.users)-1]
		s.userByName[user.Username] = &s.users[len(s.users)-1]
	}

	s.sortUsers()
	return len(s.users)
}

// GetTieAwareRank calculates the rank for a given rating.
// Users with the same rating get the same rank.
// Rank = count of users with strictly higher rating + 1
func (s *LeaderboardService) GetTieAwareRank(rating int) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Binary search to find first user with rating <= given rating
	// All users before that index have higher ratings
	idx := sort.Search(len(s.users), func(i int) bool {
		return s.users[i].Rating <= rating
	})
	return idx + 1
}

// GetLeaderboard returns a paginated list of leaderboard entries with tie-aware ranks.
func (s *LeaderboardService) GetLeaderboard(limit, offset int) ([]models.LeaderboardEntry, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	total := len(s.users)
	if offset >= total {
		return []models.LeaderboardEntry{}, total
	}

	end := offset + limit
	if end > total {
		end = total
	}

	entries := make([]models.LeaderboardEntry, 0, end-offset)
	
	for i := offset; i < end; i++ {
		user := s.users[i]
		rank := s.getTieAwareRankUnlocked(user.Rating)
		entries = append(entries, models.LeaderboardEntry{
			Rank:     rank,
			Username: user.Username,
			Rating:   user.Rating,
			ID:       user.ID,
		})
	}

	return entries, total
}

// getTieAwareRankUnlocked calculates rank without acquiring lock (caller must hold lock).
func (s *LeaderboardService) getTieAwareRankUnlocked(rating int) int {
	idx := sort.Search(len(s.users), func(i int) bool {
		return s.users[i].Rating <= rating
	})
	return idx + 1
}

// SearchUsers finds users whose username contains the query string.
func (s *LeaderboardService) SearchUsers(query string) []models.SearchResult {
	s.mu.RLock()
	defer s.mu.RUnlock()

	query = strings.ToLower(query)
	results := make([]models.SearchResult, 0)

	for _, user := range s.users {
		if strings.Contains(strings.ToLower(user.Username), query) {
			rank := s.getTieAwareRankUnlocked(user.Rating)
			results = append(results, models.SearchResult{
				GlobalRank: rank,
				Username:   user.Username,
				Rating:     user.Rating,
				ID:         user.ID,
			})
		}
	}

	return results
}

// GetUserByID returns a user by their ID with their global rank.
func (s *LeaderboardService) GetUserByID(id int) (*models.SearchResult, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.userByID[id]
	if !exists {
		return nil, false
	}

	rank := s.getTieAwareRankUnlocked(user.Rating)
	return &models.SearchResult{
		GlobalRank: rank,
		Username:   user.Username,
		Rating:     user.Rating,
		ID:         user.ID,
	}, true
}

// UpdateRating updates a user's rating and re-sorts the leaderboard.
func (s *LeaderboardService) UpdateRating(id int, newRating int) (*models.SearchResult, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.userByID[id]
	if !exists {
		return nil, false
	}

	user.Rating = newRating
	s.sortUsers()

	// Get updated user reference after sort
	user = s.userByID[id]
	rank := s.getTieAwareRankUnlocked(user.Rating)
	
	return &models.SearchResult{
		GlobalRank: rank,
		Username:   user.Username,
		Rating:     user.Rating,
		ID:         user.ID,
	}, true
}

// GetTotalUsers returns the total number of users.
func (s *LeaderboardService) GetTotalUsers() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.users)
}
