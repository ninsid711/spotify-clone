package handlers

import (
	"net/http"
	"strconv"
	"time"

	"spotify-clone/database"

	"github.com/gin-gonic/gin"
)

// GetArtistStats returns comprehensive statistics for an artist
// GET /api/artists/:id/stats
func GetArtistStats(c *gin.Context) {
	artistIDStr := c.Param("id")
	artistID, err := strconv.Atoi(artistIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid artist ID"})
		return
	}

	stats, err := database.GetArtistStats(artistID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

// GetAlbumStats returns statistics for an album (uses trigger-maintained data)
// GET /api/albums/:id/stats
func GetAlbumStats(c *gin.Context) {
	albumIDStr := c.Param("id")
	albumID, err := strconv.Atoi(albumIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid album ID"})
		return
	}

	stats, err := database.GetAlbumStats(albumID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

// AddTrackRequest represents the request body for adding a track
type AddTrackRequest struct {
	Title       string `json:"title" binding:"required"`
	ArtistID    int    `json:"artist_id" binding:"required"`
	AlbumID     int    `json:"album_id" binding:"required"`
	Duration    int    `json:"duration" binding:"required"`
	Genre       string `json:"genre" binding:"required"`
	ReleaseDate string `json:"release_date" binding:"required"`
	FileURL     string `json:"file_url" binding:"required"`
	CoverURL    string `json:"cover_url"`
}

// AddTrackWithValidation adds a new track using the stored procedure
// POST /api/tracks/add
func AddTrackWithValidation(c *gin.Context) {
	var req AddTrackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse release date
	releaseDate, err := time.Parse("2006-01-02", req.ReleaseDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid release date format. Use YYYY-MM-DD"})
		return
	}

	// Call the stored procedure
	trackID, status, err := database.AddTrackWithValidation(
		req.Title,
		req.ArtistID,
		req.AlbumID,
		req.Duration,
		req.Genre,
		req.FileURL,
		req.CoverURL,
		releaseDate,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if the procedure returned an error status
	if trackID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": status,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success":  true,
		"message":  status,
		"track_id": trackID,
	})
}

// GetAlbumDuration gets the total duration of an album using the SQL function
// GET /api/albums/:id/duration
func GetAlbumDuration(c *gin.Context) {
	albumIDStr := c.Param("id")
	albumID, err := strconv.Atoi(albumIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid album ID"})
		return
	}

	duration, err := database.GetAlbumDuration(albumID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert to minutes and format nicely
	minutes := duration / 60
	seconds := duration % 60

	c.JSON(http.StatusOK, gin.H{
		"success":            true,
		"album_id":           albumID,
		"duration_seconds":   duration,
		"duration_minutes":   float64(duration) / 60.0,
		"duration_formatted": formatDuration(duration),
		"duration_display": gin.H{
			"minutes": minutes,
			"seconds": seconds,
		},
	})
}

// Helper function to format duration
func formatDuration(seconds int) string {
	minutes := seconds / 60
	secs := seconds % 60
	return strconv.Itoa(minutes) + ":" + pad(secs)
}

func pad(n int) string {
	if n < 10 {
		return "0" + strconv.Itoa(n)
	}
	return strconv.Itoa(n)
}
