package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"spotify-clone/database"
	"spotify-clone/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetTracks returns a list of tracks with pagination
func GetTracks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	genre := c.Query("genre")
	search := c.Query("search")

	offset := (page - 1) * limit

	query := `
		SELECT t.id, t.title, t.artist_id, a.name as artist_name, 
		       t.album_id, al.title as album_name, t.duration, 
		       t.genre, t.release_date, t.file_url, t.cover_url, t.created_at
		FROM tracks t
		JOIN artists a ON t.artist_id = a.id
		JOIN albums al ON t.album_id = al.id
		WHERE 1=1
	`

	args := []interface{}{}
	if genre != "" {
		query += " AND t.genre = ?"
		args = append(args, genre)
	}
	if search != "" {
		query += " AND (t.title LIKE ? OR a.name LIKE ?)"
		searchParam := "%" + search + "%"
		args = append(args, searchParam, searchParam)
	}

	query += " ORDER BY t.created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := database.MySQL.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tracks"})
		return
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
		if err != nil {
			continue
		}
		tracks = append(tracks, track)
	}

	c.JSON(http.StatusOK, gin.H{
		"tracks": tracks,
		"page":   page,
		"limit":  limit,
	})
}

// GetTrackByID returns a single track by ID
func GetTrackByID(c *gin.Context) {
	id := c.Param("id")

	query := `
		SELECT t.id, t.title, t.artist_id, a.name as artist_name, 
		       t.album_id, al.title as album_name, t.duration, 
		       t.genre, t.release_date, t.file_url, t.cover_url, t.created_at
		FROM tracks t
		JOIN artists a ON t.artist_id = a.id
		JOIN albums al ON t.album_id = al.id
		WHERE t.id = ?
	`

	var track models.Track
	err := database.MySQL.QueryRow(query, id).Scan(
		&track.ID, &track.Title, &track.ArtistID, &track.ArtistName,
		&track.AlbumID, &track.AlbumName, &track.Duration,
		&track.Genre, &track.ReleaseDate, &track.FileURL,
		&track.CoverURL, &track.CreatedAt,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Track not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch track"})
		return
	}

	c.JSON(http.StatusOK, track)
}

// GetArtists returns a list of artists
func GetArtists(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit

	query := "SELECT id, name, bio, image_url, created_at FROM artists ORDER BY name LIMIT ? OFFSET ?"
	rows, err := database.MySQL.Query(query, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch artists"})
		return
	}
	defer rows.Close()

	artists := []models.Artist{}
	for rows.Next() {
		var artist models.Artist
		err := rows.Scan(&artist.ID, &artist.Name, &artist.Bio, &artist.ImageURL, &artist.CreatedAt)
		if err != nil {
			continue
		}
		artists = append(artists, artist)
	}

	c.JSON(http.StatusOK, gin.H{
		"artists": artists,
		"page":    page,
		"limit":   limit,
	})
}

// GetArtistByID returns a single artist with their tracks
func GetArtistByID(c *gin.Context) {
	id := c.Param("id")

	// Get artist info
	var artist models.Artist
	err := database.MySQL.QueryRow(
		"SELECT id, name, bio, image_url, created_at FROM artists WHERE id = ?", id,
	).Scan(&artist.ID, &artist.Name, &artist.Bio, &artist.ImageURL, &artist.CreatedAt)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Artist not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch artist"})
		return
	}

	// Get artist's tracks
	query := `
		SELECT t.id, t.title, t.artist_id, a.name as artist_name, 
		       t.album_id, al.title as album_name, t.duration, 
		       t.genre, t.release_date, t.file_url, t.cover_url, t.created_at
		FROM tracks t
		JOIN artists a ON t.artist_id = a.id
		JOIN albums al ON t.album_id = al.id
		WHERE t.artist_id = ?
		ORDER BY t.release_date DESC
	`

	rows, err := database.MySQL.Query(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch artist tracks"})
		return
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
		if err != nil {
			continue
		}
		tracks = append(tracks, track)
	}

	c.JSON(http.StatusOK, gin.H{
		"artist": artist,
		"tracks": tracks,
	})
}

// GetAlbums returns a list of albums
func GetAlbums(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset := (page - 1) * limit

	query := `
		SELECT al.id, al.title, al.artist_id, a.name as artist_name, 
		       al.release_date, al.cover_url, al.created_at
		FROM albums al
		JOIN artists a ON al.artist_id = a.id
		ORDER BY al.release_date DESC
		LIMIT ? OFFSET ?
	`

	rows, err := database.MySQL.Query(query, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch albums"})
		return
	}
	defer rows.Close()

	albums := []models.Album{}
	for rows.Next() {
		var album models.Album
		err := rows.Scan(
			&album.ID, &album.Title, &album.ArtistID, &album.ArtistName,
			&album.ReleaseDate, &album.CoverURL, &album.CreatedAt,
		)
		if err != nil {
			continue
		}
		albums = append(albums, album)
	}

	c.JSON(http.StatusOK, gin.H{
		"albums": albums,
		"page":   page,
		"limit":  limit,
	})
}

// Search performs a global search across tracks, artists, and albums
func Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	searchParam := "%" + query + "%"

	// Search tracks
	tracksQuery := `
		SELECT t.id, t.title, t.artist_id, a.name as artist_name, 
		       t.album_id, al.title as album_name, t.duration, 
		       t.genre, t.release_date, t.file_url, t.cover_url, t.created_at
		FROM tracks t
		JOIN artists a ON t.artist_id = a.id
		JOIN albums al ON t.album_id = al.id
		WHERE t.title LIKE ? OR a.name LIKE ?
		LIMIT 10
	`

	tracks := []models.Track{}
	rows, err := database.MySQL.Query(tracksQuery, searchParam, searchParam)
	if err == nil {
		defer rows.Close()
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
	}

	// Search artists
	artistsQuery := "SELECT id, name, bio, image_url, created_at FROM artists WHERE name LIKE ? LIMIT 10"
	artists := []models.Artist{}
	rows, err = database.MySQL.Query(artistsQuery, searchParam)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var artist models.Artist
			err := rows.Scan(&artist.ID, &artist.Name, &artist.Bio, &artist.ImageURL, &artist.CreatedAt)
			if err == nil {
				artists = append(artists, artist)
			}
		}
	}

	// Search albums
	albumsQuery := `
		SELECT al.id, al.title, al.artist_id, a.name as artist_name, 
		       al.release_date, al.cover_url, al.created_at
		FROM albums al
		JOIN artists a ON al.artist_id = a.id
		WHERE al.title LIKE ? OR a.name LIKE ?
		LIMIT 10
	`
	albums := []models.Album{}
	rows, err = database.MySQL.Query(albumsQuery, searchParam, searchParam)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var album models.Album
			err := rows.Scan(
				&album.ID, &album.Title, &album.ArtistID, &album.ArtistName,
				&album.ReleaseDate, &album.CoverURL, &album.CreatedAt,
			)
			if err == nil {
				albums = append(albums, album)
			}
		}
	}

	c.JSON(http.StatusOK, models.SearchResponse{
		Tracks:  tracks,
		Artists: artists,
		Albums:  albums,
	})
}

// RecordPlay records a track play in user's listening history and Neo4j
func RecordPlay(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	trackID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid track ID"})
		return
	}

	// Get track info from MySQL
	var track models.Track
	err = database.MySQL.QueryRow(
		"SELECT id, title, artist_id, genre, duration FROM tracks WHERE id = ?", trackID,
	).Scan(&track.ID, &track.Title, &track.ArtistID, &track.Genre, &track.Duration)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Track not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch track"})
		return
	}

	// Record in MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, _ := primitive.ObjectIDFromHex(userID.(string))
	collection := database.MongoDB.Collection("users")

	historyEntry := models.ListeningHistory{
		TrackID:   trackID,
		PlayedAt:  time.Now(),
		Duration:  track.Duration,
		Completed: true,
	}

	update := bson.M{
		"$push": bson.M{
			"listening_history": bson.M{
				"$each":  []models.ListeningHistory{historyEntry},
				"$slice": -100, // Keep only last 100 plays
			},
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	_, err = collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record play"})
		return
	}

	// Record in Neo4j for recommendations
	go recordPlayInNeo4j(userID.(string), trackID, track.ArtistID, track.Genre)

	c.JSON(http.StatusOK, gin.H{"message": "Play recorded successfully"})
}

// recordPlayInNeo4j records the play relationship in Neo4j
func recordPlayInNeo4j(userID string, trackID, artistID int, genre string) {
	ctx := context.Background()
	session := database.Neo4j.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	query := `
		MERGE (u:User {id: $userId})
		MERGE (t:Track {id: $trackId})
		MERGE (a:Artist {id: $artistId})
		MERGE (g:Genre {name: $genre})
		
		MERGE (t)-[:BY_ARTIST]->(a)
		MERGE (t)-[:HAS_GENRE]->(g)
		
		MERGE (u)-[p:PLAYED]->(t)
		ON CREATE SET p.count = 1, p.lastPlayed = datetime()
		ON MATCH SET p.count = p.count + 1, p.lastPlayed = datetime()
		
		MERGE (u)-[l:LIKES_ARTIST]->(a)
		ON CREATE SET l.count = 1
		ON MATCH SET l.count = l.count + 1
		
		MERGE (u)-[lg:LIKES_GENRE]->(g)
		ON CREATE SET lg.count = 1
		ON MATCH SET lg.count = lg.count + 1
	`

	params := map[string]interface{}{
		"userId":   userID,
		"trackId":  trackID,
		"artistId": artistID,
		"genre":    genre,
	}

	session.Run(ctx, query, params)
}
