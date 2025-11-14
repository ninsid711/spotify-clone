package handlers

import (
	"database/sql"
	"net/http"
	"spotify-clone/database"
	"spotify-clone/models"

	"github.com/gin-gonic/gin"
)

// CreatePlaylist creates a new playlist
func CreatePlaylist(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req models.CreatePlaylistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := database.MySQL.Exec(`
		INSERT INTO playlists (user_id, name, description, is_public)
		VALUES (?, ?, ?, ?)`,
		userID, req.Name, req.Description, req.IsPublic)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create playlist"})
		return
	}

	playlistID, _ := result.LastInsertId()

	playlist := models.Playlist{
		ID:          int(playlistID),
		UserID:      userID.(int),
		Name:        req.Name,
		Description: req.Description,
		IsPublic:    req.IsPublic,
		TrackIDs:    []int{},
	}

	c.JSON(http.StatusCreated, playlist)
}

// GetUserPlaylists returns all playlists for the authenticated user
func GetUserPlaylists(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	rows, err := database.MySQL.Query(`
		SELECT id, name, description, is_public, cover_url, created_at, updated_at
		FROM playlists WHERE user_id = ?
		ORDER BY updated_at DESC`, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch playlists"})
		return
	}
	defer rows.Close()

	playlists := []models.Playlist{}
	for rows.Next() {
		var playlist models.Playlist
		rows.Scan(&playlist.ID, &playlist.Name, &playlist.Description, &playlist.IsPublic,
			&playlist.CoverURL, &playlist.CreatedAt, &playlist.UpdatedAt)
		playlist.UserID = userID.(int)

		// Get track count
		var trackIDs []int
		trackRows, _ := database.MySQL.Query("SELECT track_id FROM playlist_tracks WHERE playlist_id = ? ORDER BY position", playlist.ID)
		for trackRows.Next() {
			var trackID int
			trackRows.Scan(&trackID)
			trackIDs = append(trackIDs, trackID)
		}
		trackRows.Close()
		playlist.TrackIDs = trackIDs

		playlists = append(playlists, playlist)
	}

	c.JSON(http.StatusOK, gin.H{"playlists": playlists})
}

// GetPlaylistByID returns a single playlist with full track details
func GetPlaylistByID(c *gin.Context) {
	playlistID := c.Param("id")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var playlist models.Playlist
	err := database.MySQL.QueryRow(`
		SELECT id, user_id, name, description, is_public, cover_url, created_at, updated_at
		FROM playlists WHERE id = ?`, playlistID).Scan(
		&playlist.ID, &playlist.UserID, &playlist.Name, &playlist.Description,
		&playlist.IsPublic, &playlist.CoverURL, &playlist.CreatedAt, &playlist.UpdatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Playlist not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch playlist"})
		return
	}

	// Check if user has access (owner or public playlist)
	if playlist.UserID != userID.(int) && !playlist.IsPublic {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Fetch track details
	tracks := []models.Track{}
	rows, err := database.MySQL.Query(`
		SELECT t.id, t.title, t.artist_id, a.name as artist_name,
		       t.album_id, al.title as album_name, t.duration,
		       t.genre, t.release_date, t.file_url, t.cover_url, t.created_at
		FROM playlist_tracks pt
		JOIN tracks t ON pt.track_id = t.id
		JOIN artists a ON t.artist_id = a.id
		JOIN albums al ON t.album_id = al.id
		WHERE pt.playlist_id = ?
		ORDER BY pt.position`, playlistID)

	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var track models.Track
			rows.Scan(&track.ID, &track.Title, &track.ArtistID, &track.ArtistName,
				&track.AlbumID, &track.AlbumName, &track.Duration, &track.Genre,
				&track.ReleaseDate, &track.FileURL, &track.CoverURL, &track.CreatedAt)
			tracks = append(tracks, track)
			playlist.TrackIDs = append(playlist.TrackIDs, track.ID)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"playlist": playlist,
		"tracks":   tracks,
	})
}

// AddTrackToPlaylist adds a track to a playlist
func AddTrackToPlaylist(c *gin.Context) {
	playlistID := c.Param("id")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req models.AddTrackToPlaylistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify track exists
	var trackExists bool
	err := database.MySQL.QueryRow("SELECT EXISTS(SELECT 1 FROM tracks WHERE id = ?)", req.TrackID).Scan(&trackExists)
	if err != nil || !trackExists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Track not found"})
		return
	}

	// Check if user owns the playlist
	var ownerID int
	err = database.MySQL.QueryRow("SELECT user_id FROM playlists WHERE id = ?", playlistID).Scan(&ownerID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Playlist not found"})
		return
	}
	if ownerID != userID.(int) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Get current max position
	var maxPosition int
	database.MySQL.QueryRow("SELECT COALESCE(MAX(position), -1) FROM playlist_tracks WHERE playlist_id = ?", playlistID).Scan(&maxPosition)

	// Add track to playlist
	_, err = database.MySQL.Exec(`
		INSERT INTO playlist_tracks (playlist_id, track_id, position)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE position = VALUES(position)`,
		playlistID, req.TrackID, maxPosition+1)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add track to playlist"})
		return
	}

	// Update playlist updated_at
	database.MySQL.Exec("UPDATE playlists SET updated_at = NOW() WHERE id = ?", playlistID)

	c.JSON(http.StatusOK, gin.H{"message": "Track added to playlist successfully"})
}

// RemoveTrackFromPlaylist removes a track from a playlist
func RemoveTrackFromPlaylist(c *gin.Context) {
	playlistID := c.Param("id")
	trackID := c.Param("trackId")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Check if user owns the playlist
	var ownerID int
	err := database.MySQL.QueryRow("SELECT user_id FROM playlists WHERE id = ?", playlistID).Scan(&ownerID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Playlist not found"})
		return
	}
	if ownerID != userID.(int) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Remove track from playlist
	_, err = database.MySQL.Exec("DELETE FROM playlist_tracks WHERE playlist_id = ? AND track_id = ?", playlistID, trackID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove track from playlist"})
		return
	}

	// Update playlist updated_at
	database.MySQL.Exec("UPDATE playlists SET updated_at = NOW() WHERE id = ?", playlistID)

	c.JSON(http.StatusOK, gin.H{"message": "Track removed from playlist successfully"})
}

// DeletePlaylist deletes a playlist
func DeletePlaylist(c *gin.Context) {
	playlistID := c.Param("id")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Check if user owns the playlist
	var ownerID int
	err := database.MySQL.QueryRow("SELECT user_id FROM playlists WHERE id = ?", playlistID).Scan(&ownerID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Playlist not found"})
		return
	}
	if ownerID != userID.(int) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Delete playlist (CASCADE will delete playlist_tracks)
	_, err = database.MySQL.Exec("DELETE FROM playlists WHERE id = ?", playlistID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete playlist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Playlist deleted successfully"})
}

// UpdatePlaylist updates playlist details
func UpdatePlaylist(c *gin.Context) {
	playlistID := c.Param("id")
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req models.CreatePlaylistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user owns the playlist
	var ownerID int
	err := database.MySQL.QueryRow("SELECT user_id FROM playlists WHERE id = ?", playlistID).Scan(&ownerID)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Playlist not found"})
		return
	}
	if ownerID != userID.(int) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Update playlist
	_, err = database.MySQL.Exec(`
		UPDATE playlists 
		SET name = ?, description = ?, is_public = ?, updated_at = NOW()
		WHERE id = ?`,
		req.Name, req.Description, req.IsPublic, playlistID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update playlist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Playlist updated successfully"})
}
