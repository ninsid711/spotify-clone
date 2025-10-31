package handlers

import (
	"context"
	"net/http"
	"spotify-clone/database"
	"spotify-clone/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	playlist := models.Playlist{
		UserID:      userID.(string),
		Name:        req.Name,
		Description: req.Description,
		TrackIDs:    []int{},
		IsPublic:    req.IsPublic,
		CoverURL:    "",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	collection := database.MongoDB.Collection("playlists")
	result, err := collection.InsertOne(ctx, playlist)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create playlist"})
		return
	}

	playlist.ID = result.InsertedID.(primitive.ObjectID).Hex()
	c.JSON(http.StatusCreated, playlist)
}

// GetUserPlaylists returns all playlists for the authenticated user
func GetUserPlaylists(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection := database.MongoDB.Collection("playlists")
	cursor, err := collection.Find(ctx, bson.M{"user_id": userID.(string)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch playlists"})
		return
	}
	defer cursor.Close(ctx)

	var playlists []models.Playlist
	if err = cursor.All(ctx, &playlists); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode playlists"})
		return
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(playlistID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid playlist ID"})
		return
	}

	collection := database.MongoDB.Collection("playlists")
	var playlist models.Playlist
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&playlist)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Playlist not found"})
		return
	}

	// Check if user has access (owner or public playlist)
	if playlist.UserID != userID.(string) && !playlist.IsPublic {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Fetch track details from MySQL
	tracks := []models.Track{}
	if len(playlist.TrackIDs) > 0 {
		tracks = getTracksByIDs(playlist.TrackIDs)
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

	// Verify track exists in MySQL
	var trackExists bool
	err := database.MySQL.QueryRow("SELECT EXISTS(SELECT 1 FROM tracks WHERE id = ?)", req.TrackID).Scan(&trackExists)
	if err != nil || !trackExists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Track not found"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(playlistID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid playlist ID"})
		return
	}

	collection := database.MongoDB.Collection("playlists")

	// Check if user owns the playlist
	var playlist models.Playlist
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&playlist)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Playlist not found"})
		return
	}

	if playlist.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Add track to playlist (avoid duplicates)
	update := bson.M{
		"$addToSet": bson.M{"track_ids": req.TrackID},
		"$set":      bson.M{"updated_at": time.Now()},
	}

	_, err = collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add track to playlist"})
		return
	}

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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(playlistID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid playlist ID"})
		return
	}

	collection := database.MongoDB.Collection("playlists")

	// Check if user owns the playlist
	var playlist models.Playlist
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&playlist)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Playlist not found"})
		return
	}

	if playlist.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Remove track from playlist
	update := bson.M{
		"$pull": bson.M{"track_ids": trackID},
		"$set":  bson.M{"updated_at": time.Now()},
	}

	_, err = collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove track from playlist"})
		return
	}

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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(playlistID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid playlist ID"})
		return
	}

	collection := database.MongoDB.Collection("playlists")

	// Check if user owns the playlist
	var playlist models.Playlist
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&playlist)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Playlist not found"})
		return
	}

	if playlist.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Delete playlist
	_, err = collection.DeleteOne(ctx, bson.M{"_id": objectID})
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(playlistID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid playlist ID"})
		return
	}

	collection := database.MongoDB.Collection("playlists")

	// Check if user owns the playlist
	var playlist models.Playlist
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&playlist)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Playlist not found"})
		return
	}

	if playlist.UserID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Update playlist
	update := bson.M{
		"$set": bson.M{
			"name":        req.Name,
			"description": req.Description,
			"is_public":   req.IsPublic,
			"updated_at":  time.Now(),
		},
	}

	_, err = collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update playlist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Playlist updated successfully"})
}

// getTracksByIDs fetches multiple tracks by their IDs from MySQL
func getTracksByIDs(trackIDs []int) []models.Track {
	if len(trackIDs) == 0 {
		return []models.Track{}
	}

	query := `
		SELECT t.id, t.title, t.artist_id, a.name as artist_name, 
		       t.album_id, al.title as album_name, t.duration, 
		       t.genre, t.release_date, t.file_url, t.cover_url, t.created_at
		FROM tracks t
		JOIN artists a ON t.artist_id = a.id
		JOIN albums al ON t.album_id = al.id
		WHERE t.id IN (?` + string(make([]byte, len(trackIDs)-1)) + `)
	`

	// Build placeholders
	args := make([]interface{}, len(trackIDs))
	for i, id := range trackIDs {
		args[i] = id
		if i > 0 {
			query = query[:len(query)-1] + ",?)"
		}
	}

	rows, err := database.MySQL.Query(query, args...)
	if err != nil {
		return []models.Track{}
	}
	defer rows.Close()

	tracks := []models.Track{}
	for rows.Next() {
		var track models.Track
		err := rows.Scan(
			&track.ID, &track.Title, &track.ArtistID, &track.ArtistName,
			&track.AlbumID, &track.AlbumName, &track.Duration,
			&track.Genre, &track.ReleaseDate, &track.FileURL,
			&track.CoverURL, &track.CreatedAt,
		)
		if err == nil {
			tracks = append(tracks, track)
		}
	}

	return tracks
}
