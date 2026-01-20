# Leaderboard Backend

This project implements a high-performance, in-memory leaderboard system using Go (Golang) and the Gin web framework. It supports real-time ranking, searching, and updates for 10,000+ users.

## Features
- **In-Memory Storage**: Uses optimized data structures for fast read/write access.
- **Tie-Aware Ranking**: Users with the same score share the same rank (e.g., 1, 2, 2, 4).
- **Thread Safety**: Uses `sync.RWMutex` to ensure safe concurrent operations.
- **Search**: Username search functionality.
- **REST API**: Clean endpoints for frontend integration.

## Tech Stack
- **Language**: Go 1.21+
- **Framework**: Gin Gonic
- **Storage**: In-Memory (Slice + Maps)

## Architecture
The system uses a Dual-Structure approach for O(log n) reads and efficient updates:
1.  **Slice (`[]User`)**: Kept sorted by rating (descending). Used for determining rank via Binary Search and paginated lists.
2.  **Maps (`map[int]*User`, `map[string]*User`)**: Provides O(1) direct access to user data by ID or Username.

### Operations Complexity
- **Get Rank**: `O(log N)` (Binary Search on sorted slice)
- **Update Rating**: `O(N log N)` (Re-sorts slice after update - *Note: Optimization possible with specialized trees*)
- **Get Leaderboard**: `O(1)` (Slice access)
- **Search**: `O(N)` (Linear scan)

## How to Run

### Prerequisites
- Go installed on your machine.

### Steps
1. Navigate to the backend directory:
   ```bash
   cd leaderboard-backend
   ```
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Run the server:
   ```bash
   go run main.go
   ```
   *The server starts on port `8080` and automatically seeds 10,000 users.*

## API Endpoints

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/api/leaderboard?limit=50&offset=0` | Get paginated leaderboard entries |
| `GET` | `/api/search?q=username` | Search users by name |
| `GET` | `/api/user/:id/rank` | Get a specific user's rank and details |
| `POST` | `/api/user/:id/rating` | Update a user's rating (Body: `{"rating": 1500}`) |
| `POST` | `/api/seed` | Re-seed the leaderboard with N users |
| `GET` | `/api/stats` | Get total user count |
| `GET` | `/health` | Health check |

## Testing
You can test the API using `curl` or Postman.

**Get Top 5 Users:**
```bash
curl "http://localhost:8080/api/leaderboard?limit=5"
```

**Search for "Rahul":**
```bash
curl "http://localhost:8080/api/search?q=rahul"
```
