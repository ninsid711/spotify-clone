package handlers

import (
	"context"
	"net/http"
	"spotify-clone/database"
	"spotify-clone/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// GetRecommendations returns personalized track recommendations based on user's listening history
func GetRecommendations(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	ctx := context.Background()
	session := database.Neo4j.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	// Get recommendations based on user's listening patterns
	// This query finds tracks that:
	// 1. Are by artists the user likes
	// 2. Are in genres the user prefers
	// 3. Haven't been played by the user yet (or played less frequently)
	// 4. Are popular among users with similar tastes
	query := `
		MATCH (u:User {id: $userId})
		
		// Find tracks by artists the user likes
		OPTIONAL MATCH (u)-[la:LIKES_ARTIST]->(a:Artist)<-[:BY_ARTIST]-(t1:Track)
		WHERE NOT (u)-[:PLAYED]->(t1) OR (u)-[p:PLAYED]->(t1) WHERE p.count < 3
		WITH u, t1, la.count as artistScore
		
		// Find tracks in genres the user likes
		OPTIONAL MATCH (u)-[lg:LIKES_GENRE]->(g:Genre)<-[:HAS_GENRE]-(t2:Track)
		WHERE NOT (u)-[:PLAYED]->(t2) OR (u)-[p2:PLAYED]->(t2) WHERE p2.count < 3
		WITH u, COLLECT(DISTINCT t1) as artistTracks, COLLECT(DISTINCT t2) as genreTracks, 
		     MAX(artistScore) as maxArtistScore
		
		// Combine and score tracks
		UNWIND (artistTracks + genreTracks) as track
		WITH DISTINCT track, u, maxArtistScore
		
		// Calculate collaborative filtering score
		OPTIONAL MATCH (other:User)-[:PLAYED]->(track)
		WHERE other <> u
		WITH track, COUNT(DISTINCT other) as popularity, maxArtistScore
		
		RETURN track.id as trackId, 
		       (popularity * 0.5 + maxArtistScore * 0.5) as score
		ORDER BY score DESC
		LIMIT $limit
	`

	params := map[string]interface{}{
		"userId": userID.(string),
		"limit":  limit,
	}

	result, err := session.Run(ctx, query, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get recommendations"})
		return
	}

	trackIDs := []int{}
	for result.Next(ctx) {
		record := result.Record()
		if trackIDValue, ok := record.Get("trackId"); ok {
			trackID := int(trackIDValue.(int64))
			trackIDs = append(trackIDs, trackID)
		}
	}

	if err = result.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing recommendations"})
		return
	}

	// If no recommendations found, get popular tracks
	if len(trackIDs) == 0 {
		trackIDs = getPopularTracks(limit)
	}

	// Fetch track details from MySQL
	tracks := []models.Track{}
	if len(trackIDs) > 0 {
		tracks = getTrackDetailsByIDs(trackIDs)
	}

	c.JSON(http.StatusOK, models.RecommendationResponse{
		Tracks: tracks,
		Reason: "Based on your listening history and preferences",
	})
}

// GetSimilarTracks returns tracks similar to a given track
func GetSimilarTracks(c *gin.Context) {
	trackID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid track ID"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	ctx := context.Background()
	session := database.Neo4j.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	// Find similar tracks based on:
	// 1. Same artist
	// 2. Same genre
	// 3. Users who played this track also played...
	query := `
		MATCH (t:Track {id: $trackId})
		
		// Same artist tracks
		OPTIONAL MATCH (t)-[:BY_ARTIST]->(a:Artist)<-[:BY_ARTIST]-(similar1:Track)
		WHERE similar1 <> t
		
		// Same genre tracks
		OPTIONAL MATCH (t)-[:HAS_GENRE]->(g:Genre)<-[:HAS_GENRE]-(similar2:Track)
		WHERE similar2 <> t
		
		// Collaborative filtering - users who played this also played
		OPTIONAL MATCH (u:User)-[:PLAYED]->(t)
		WITH t, COLLECT(DISTINCT similar1) as artistTracks, 
		     COLLECT(DISTINCT similar2) as genreTracks, 
		     COLLECT(u) as users
		
		UNWIND users as user
		MATCH (user)-[:PLAYED]->(similar3:Track)
		WHERE similar3 <> t
		
		WITH (artistTracks + genreTracks + COLLECT(DISTINCT similar3)) as allSimilar
		UNWIND allSimilar as similar
		WITH DISTINCT similar, COUNT(*) as score
		
		RETURN similar.id as trackId, score
		ORDER BY score DESC
		LIMIT $limit
	`

	params := map[string]interface{}{
		"trackId": trackID,
		"limit":   limit,
	}

	result, err := session.Run(ctx, query, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get similar tracks"})
		return
	}

	trackIDs := []int{}
	for result.Next(ctx) {
		record := result.Record()
		if trackIDValue, ok := record.Get("trackId"); ok {
			tid := int(trackIDValue.(int64))
			trackIDs = append(trackIDs, tid)
		}
	}

	if err = result.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing similar tracks"})
		return
	}

	// Fetch track details from MySQL
	tracks := []models.Track{}
	if len(trackIDs) > 0 {
		tracks = getTrackDetailsByIDs(trackIDs)
	}

	c.JSON(http.StatusOK, models.RecommendationResponse{
		Tracks: tracks,
		Reason: "Tracks similar to what you're listening to",
	})
}

// GetTrendingTracks returns currently trending tracks
func GetTrendingTracks(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	ctx := context.Background()
	session := database.Neo4j.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	// Get tracks with most plays in recent time
	query := `
		MATCH (u:User)-[p:PLAYED]->(t:Track)
		WHERE datetime() - p.lastPlayed < duration({days: 7})
		WITH t, SUM(p.count) as playCount
		ORDER BY playCount DESC
		LIMIT $limit
		RETURN t.id as trackId, playCount
	`

	params := map[string]interface{}{
		"limit": limit,
	}

	result, err := session.Run(ctx, query, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get trending tracks"})
		return
	}

	trackIDs := []int{}
	for result.Next(ctx) {
		record := result.Record()
		if trackIDValue, ok := record.Get("trackId"); ok {
			trackID := int(trackIDValue.(int64))
			trackIDs = append(trackIDs, trackID)
		}
	}

	if err = result.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing trending tracks"})
		return
	}

	// If no trending tracks found, get popular tracks
	if len(trackIDs) == 0 {
		trackIDs = getPopularTracks(limit)
	}

	// Fetch track details from MySQL
	tracks := []models.Track{}
	if len(trackIDs) > 0 {
		tracks = getTrackDetailsByIDs(trackIDs)
	}

	c.JSON(http.StatusOK, models.RecommendationResponse{
		Tracks: tracks,
		Reason: "Trending this week",
	})
}

// GetGenreRecommendations returns tracks from a specific genre
func GetGenreRecommendations(c *gin.Context) {
	genre := c.Param("genre")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	userID, exists := c.Get("user_id")

	ctx := context.Background()
	session := database.Neo4j.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	var query string
	params := map[string]interface{}{
		"genre": genre,
		"limit": limit,
	}

	if exists {
		// Personalized genre recommendations
		query = `
			MATCH (g:Genre {name: $genre})<-[:HAS_GENRE]-(t:Track)
			OPTIONAL MATCH (u:User {id: $userId})-[p:PLAYED]->(t)
			WITH t, COALESCE(p.count, 0) as playCount
			WHERE playCount < 3
			
			OPTIONAL MATCH (other:User)-[:PLAYED]->(t)
			WITH t, COUNT(other) as popularity
			ORDER BY popularity DESC
			LIMIT $limit
			RETURN t.id as trackId
		`
		params["userId"] = userID.(string)
	} else {
		// General popular tracks in genre
		query = `
			MATCH (g:Genre {name: $genre})<-[:HAS_GENRE]-(t:Track)
			OPTIONAL MATCH (u:User)-[:PLAYED]->(t)
			WITH t, COUNT(u) as popularity
			ORDER BY popularity DESC
			LIMIT $limit
			RETURN t.id as trackId
		`
	}

	result, err := session.Run(ctx, query, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get genre recommendations"})
		return
	}

	trackIDs := []int{}
	for result.Next(ctx) {
		record := result.Record()
		if trackIDValue, ok := record.Get("trackId"); ok {
			trackID := int(trackIDValue.(int64))
			trackIDs = append(trackIDs, trackID)
		}
	}

	if err = result.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing genre recommendations"})
		return
	}

	// Fetch track details from MySQL
	tracks := []models.Track{}
	if len(trackIDs) > 0 {
		tracks = getTrackDetailsByIDs(trackIDs)
	}

	c.JSON(http.StatusOK, models.RecommendationResponse{
		Tracks: tracks,
		Reason: "Popular tracks in " + genre,
	})
}

// getTrackDetailsByIDs fetches track details from MySQL
func getTrackDetailsByIDs(trackIDs []int) []models.Track {
	if len(trackIDs) == 0 {
		return []models.Track{}
	}

	// Build query with proper number of placeholders
	query := `
		SELECT t.id, t.title, t.artist_id, a.name as artist_name, 
		       t.album_id, al.title as album_name, t.duration, 
		       t.genre, t.release_date, t.file_url, t.cover_url, t.created_at
		FROM tracks t
		JOIN artists a ON t.artist_id = a.id
		JOIN albums al ON t.album_id = al.id
		WHERE t.id IN (?`

	args := make([]interface{}, len(trackIDs))
	args[0] = trackIDs[0]
	for i := 1; i < len(trackIDs); i++ {
		query += ",?"
		args[i] = trackIDs[i]
	}
	query += ")"

	rows, err := database.MySQL.Query(query, args...)
	if err != nil {
		return []models.Track{}
	}
	defer rows.Close()

	// Create a map to maintain order
	trackMap := make(map[int]models.Track)
	for rows.Next() {
		var track models.Track
		err := rows.Scan(
			&track.ID, &track.Title, &track.ArtistID, &track.ArtistName,
			&track.AlbumID, &track.AlbumName, &track.Duration,
			&track.Genre, &track.ReleaseDate, &track.FileURL,
			&track.CoverURL, &track.CreatedAt,
		)
		if err == nil {
			trackMap[track.ID] = track
		}
	}

	// Return tracks in the order of input IDs
	tracks := []models.Track{}
	for _, id := range trackIDs {
		if track, exists := trackMap[id]; exists {
			tracks = append(tracks, track)
		}
	}

	return tracks
}

// getPopularTracks returns popular tracks from MySQL as fallback
func getPopularTracks(limit int) []int {
	query := "SELECT id FROM tracks ORDER BY created_at DESC LIMIT ?"
	rows, err := database.MySQL.Query(query, limit)
	if err != nil {
		return []int{}
	}
	defer rows.Close()

	trackIDs := []int{}
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err == nil {
			trackIDs = append(trackIDs, id)
		}
	}

	return trackIDs
}
