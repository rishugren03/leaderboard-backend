package handlers

import (
	"leaderboard-backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// LeaderboardHandler handles HTTP requests for leaderboard operations.
type LeaderboardHandler struct {
	service *services.LeaderboardService
}

// NewLeaderboardHandler creates a new handler with the given service.
func NewLeaderboardHandler(service *services.LeaderboardService) *LeaderboardHandler {
	return &LeaderboardHandler{service: service}
}

// GetLeaderboard handles GET /api/leaderboard
// Query params: limit (default 50), offset (default 0)
func (h *LeaderboardHandler) GetLeaderboard(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	entries, total := h.service.GetLeaderboard(limit, offset)
	
	c.JSON(http.StatusOK, gin.H{
		"data":   entries,
		"total":  total,
		"limit":  limit,
		"offset": offset,
		"hasMore": offset+limit < total,
	})
}

// SearchUsers handles GET /api/search
// Query params: q (search query)
func (h *LeaderboardHandler) SearchUsers(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' is required"})
		return
	}

	results := h.service.SearchUsers(query)
	c.JSON(http.StatusOK, gin.H{
		"data":  results,
		"count": len(results),
		"query": query,
	})
}

// GetUserRank handles GET /api/user/:id/rank
func (h *LeaderboardHandler) GetUserRank(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	result, found := h.service.GetUserByID(id)
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// UpdateRating handles POST /api/user/:id/rating
// Body: { "rating": int }
func (h *LeaderboardHandler) UpdateRating(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var body struct {
		Rating int `json:"rating" binding:"required,min=100,max=5000"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "rating must be between 100 and 5000"})
		return
	}

	result, found := h.service.UpdateRating(id, body.Rating)
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// SeedUsers handles POST /api/seed
// Query params: count (default 10000)
func (h *LeaderboardHandler) SeedUsers(c *gin.Context) {
	count, _ := strconv.Atoi(c.DefaultQuery("count", "10000"))
	if count <= 0 {
		count = 10000
	}
	if count > 100000 {
		count = 100000
	}

	seeded := h.service.SeedUsers(count, 100, 5000)
	c.JSON(http.StatusOK, gin.H{
		"message": "users seeded successfully",
		"count":   seeded,
	})
}

// GetStats handles GET /api/stats
func (h *LeaderboardHandler) GetStats(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"totalUsers": h.service.GetTotalUsers(),
	})
}
