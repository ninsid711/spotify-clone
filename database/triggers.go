package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

// InitTriggersProceduresFunctions loads and executes the SQL file containing triggers, procedures, and functions
func InitTriggersProceduresFunctions() error {
	log.Println("⚙️  Initializing triggers, procedures, and functions...")

	// Execute statements directly instead of reading from file
	// This avoids DELIMITER issues with the Go MySQL driver

	ctx := context.Background()

	// Create utility tables
	if err := createUtilityTables(ctx); err != nil {
		log.Printf("⚠️  Warning creating utility tables: %v", err)
	}

	// Create triggers
	if err := createTriggers(ctx); err != nil {
		log.Printf("⚠️  Warning creating triggers: %v", err)
	}

	// Create function
	if err := createFunctions(ctx); err != nil {
		log.Printf("⚠️  Warning creating functions: %v", err)
	}

	// Create procedures
	if err := createProcedures(ctx); err != nil {
		log.Printf("⚠️  Warning creating procedures: %v", err)
	}

	// Initialize existing data
	if err := initializeExistingData(ctx); err != nil {
		log.Printf("⚠️  Warning initializing data: %v", err)
	}

	log.Println("✅ Triggers, procedures, and functions initialized successfully")
	return nil
}

// createUtilityTables creates supporting tables for triggers
func createUtilityTables(ctx context.Context) error {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS album_stats (
			album_id INT PRIMARY KEY,
			track_count INT DEFAULT 0,
			total_duration INT DEFAULT 0,
			last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (album_id) REFERENCES albums(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS track_stats (
			track_id INT PRIMARY KEY,
			play_count INT DEFAULT 0,
			last_played TIMESTAMP NULL,
			FOREIGN KEY (track_id) REFERENCES tracks(id) ON DELETE CASCADE
		)`,
	}

	for _, table := range tables {
		if _, err := MySQL.ExecContext(ctx, table); err != nil {
			return fmt.Errorf("error creating utility table: %v", err)
		}
	}
	return nil
}

// createTriggers creates the triggers
func createTriggers(ctx context.Context) error {
	triggers := []string{
		`DROP TRIGGER IF EXISTS after_track_insert`,
		`CREATE TRIGGER after_track_insert
		AFTER INSERT ON tracks
		FOR EACH ROW
		BEGIN
			INSERT INTO album_stats (album_id, track_count, total_duration)
			VALUES (NEW.album_id, 1, NEW.duration)
			ON DUPLICATE KEY UPDATE
				track_count = track_count + 1,
				total_duration = total_duration + NEW.duration;
			
			INSERT INTO track_stats (track_id, play_count)
			VALUES (NEW.id, 0);
		END`,
		`DROP TRIGGER IF EXISTS after_track_delete`,
		`CREATE TRIGGER after_track_delete
		AFTER DELETE ON tracks
		FOR EACH ROW
		BEGIN
			UPDATE album_stats
			SET track_count = track_count - 1,
				total_duration = total_duration - OLD.duration
			WHERE album_id = OLD.album_id;
			
			DELETE FROM track_stats WHERE track_id = OLD.id;
		END`,
		// New trigger: update track_stats on play insert
		`DROP TRIGGER IF EXISTS after_play_insert`,
		`CREATE TRIGGER after_play_insert
		AFTER INSERT ON plays
		FOR EACH ROW
		UPDATE track_stats
		SET play_count = play_count + 1,
			last_played = NEW.played_at
		WHERE track_id = NEW.track_id`,
	}

	for _, trigger := range triggers {
		if _, err := MySQL.ExecContext(ctx, trigger); err != nil {
			return fmt.Errorf("error creating trigger: %v", err)
		}
	}
	return nil
}

// createFunctions creates the SQL functions
func createFunctions(ctx context.Context) error {
	functions := []string{
		`DROP FUNCTION IF EXISTS get_album_duration`,
		`CREATE FUNCTION get_album_duration(album_id_param INT)
		RETURNS INT
		DETERMINISTIC
		READS SQL DATA
		BEGIN
			DECLARE total INT;
			SELECT COALESCE(SUM(duration), 0)
			INTO total
			FROM tracks
			WHERE album_id = album_id_param;
			RETURN total;
		END`,
	}

	for _, fn := range functions {
		if _, err := MySQL.ExecContext(ctx, fn); err != nil {
			return fmt.Errorf("error creating function: %v", err)
		}
	}
	return nil
}

// createProcedures creates the stored procedures
func createProcedures(ctx context.Context) error {
	procedures := []string{
		`DROP PROCEDURE IF EXISTS add_track`,
		`CREATE PROCEDURE add_track(
			IN p_title VARCHAR(255),
			IN p_artist_id INT,
			IN p_album_id INT,
			IN p_duration INT,
			IN p_genre VARCHAR(100),
			IN p_release_date DATE,
			IN p_file_url VARCHAR(500),
			IN p_cover_url VARCHAR(500),
			OUT p_track_id INT,
			OUT p_status VARCHAR(100)
		)
		BEGIN
			DECLARE artist_exists INT;
			DECLARE album_exists INT;
			DECLARE album_artist_id INT;
			
			SELECT COUNT(*) INTO artist_exists
			FROM artists
			WHERE id = p_artist_id;
			
			IF artist_exists = 0 THEN
				SET p_status = 'ERROR: Artist does not exist';
				SET p_track_id = NULL;
			ELSE
				SELECT COUNT(*), MAX(artist_id) INTO album_exists, album_artist_id
				FROM albums
				WHERE id = p_album_id;
				
				IF album_exists = 0 THEN
					SET p_status = 'ERROR: Album does not exist';
					SET p_track_id = NULL;
				ELSEIF album_artist_id != p_artist_id THEN
					SET p_status = 'ERROR: Album does not belong to the specified artist';
					SET p_track_id = NULL;
				ELSE
					INSERT INTO tracks (title, artist_id, album_id, duration, genre, release_date, file_url, cover_url)
					VALUES (p_title, p_artist_id, p_album_id, p_duration, p_genre, p_release_date, p_file_url, p_cover_url);
					
					SET p_track_id = LAST_INSERT_ID();
					SET p_status = 'SUCCESS: Track added successfully';
				END IF;
			END IF;
		END`,
		`DROP PROCEDURE IF EXISTS get_artist_stats`,
		`CREATE PROCEDURE get_artist_stats(IN p_artist_id INT)
		BEGIN
			DECLARE artist_name_var VARCHAR(255);
			
			SELECT name INTO artist_name_var
			FROM artists
			WHERE id = p_artist_id;
			
			SELECT 
				p_artist_id AS artist_id,
				artist_name_var AS artist_name,
				COUNT(DISTINCT a.id) AS total_albums,
				COUNT(DISTINCT t.id) AS total_tracks,
				COALESCE(SUM(t.duration), 0) AS total_duration_seconds,
				ROUND(COALESCE(SUM(t.duration), 0) / 60.0, 2) AS total_duration_minutes,
				COUNT(DISTINCT t.genre) AS unique_genres,
				GROUP_CONCAT(DISTINCT t.genre ORDER BY t.genre SEPARATOR ', ') AS genres,
				MIN(t.release_date) AS first_release,
				MAX(t.release_date) AS latest_release,
				COALESCE(AVG(ts.play_count), 0) AS avg_plays_per_track
			FROM artists ar
			LEFT JOIN albums a ON ar.id = a.artist_id
			LEFT JOIN tracks t ON ar.id = t.artist_id
			LEFT JOIN track_stats ts ON t.id = ts.track_id
			WHERE ar.id = p_artist_id
			GROUP BY ar.id, ar.name;
		END`,
	}

	for _, proc := range procedures {
		if _, err := MySQL.ExecContext(ctx, proc); err != nil {
			return fmt.Errorf("error creating procedure: %v", err)
		}
	}
	return nil
}

// initializeExistingData initializes statistics for existing data
func initializeExistingData(ctx context.Context) error {
	statements := []string{
		`INSERT INTO album_stats (album_id, track_count, total_duration)
		SELECT album_id, COUNT(*) as track_count, SUM(duration) as total_duration
		FROM tracks
		GROUP BY album_id
		ON DUPLICATE KEY UPDATE
			track_count = VALUES(track_count),
			total_duration = VALUES(total_duration)`,
		`INSERT INTO track_stats (track_id, play_count)
		SELECT id, 0
		FROM tracks
		ON DUPLICATE KEY UPDATE play_count = play_count`,
	}

	for _, stmt := range statements {
		if _, err := MySQL.ExecContext(ctx, stmt); err != nil {
			return fmt.Errorf("error initializing data: %v", err)
		}
	}
	return nil
}

// AddTrackWithValidation calls the stored procedure to add a track with validation
func AddTrackWithValidation(title string, artistID, albumID, duration int, genre, fileURL, coverURL string, releaseDate time.Time) (int, string, error) {
	var trackID sql.NullInt64
	var status string

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := MySQL.ExecContext(ctx,
		"CALL add_track(?, ?, ?, ?, ?, ?, ?, ?, @track_id, @status)",
		title, artistID, albumID, duration, genre, releaseDate.Format("2006-01-02"), fileURL, coverURL,
	)
	if err != nil {
		return 0, "", fmt.Errorf("error calling add_track procedure: %v", err)
	}

	// Get the output parameters
	row := MySQL.QueryRowContext(ctx, "SELECT @track_id, @status")
	err = row.Scan(&trackID, &status)
	if err != nil {
		return 0, "", fmt.Errorf("error getting output parameters: %v", err)
	}

	var finalTrackID int
	if trackID.Valid {
		finalTrackID = int(trackID.Int64)
	}

	return finalTrackID, status, nil
}

// ArtistStats represents statistics for an artist
type ArtistStats struct {
	ArtistID             int       `json:"artist_id"`
	ArtistName           string    `json:"artist_name"`
	TotalAlbums          int       `json:"total_albums"`
	TotalTracks          int       `json:"total_tracks"`
	TotalDurationSeconds int       `json:"total_duration_seconds"`
	TotalDurationMinutes float64   `json:"total_duration_minutes"`
	UniqueGenres         int       `json:"unique_genres"`
	Genres               string    `json:"genres"`
	FirstRelease         time.Time `json:"first_release"`
	LatestRelease        time.Time `json:"latest_release"`
	AvgPlaysPerTrack     float64   `json:"avg_plays_per_track"`
}

// GetArtistStats calls the stored procedure to get artist statistics
func GetArtistStats(artistID int) (*ArtistStats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := MySQL.QueryContext(ctx, "CALL get_artist_stats(?)", artistID)
	if err != nil {
		return nil, fmt.Errorf("error calling get_artist_stats procedure: %v", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, fmt.Errorf("artist not found")
	}

	var stats ArtistStats
	var firstRelease, latestRelease sql.NullTime

	err = rows.Scan(
		&stats.ArtistID,
		&stats.ArtistName,
		&stats.TotalAlbums,
		&stats.TotalTracks,
		&stats.TotalDurationSeconds,
		&stats.TotalDurationMinutes,
		&stats.UniqueGenres,
		&stats.Genres,
		&firstRelease,
		&latestRelease,
		&stats.AvgPlaysPerTrack,
	)
	if err != nil {
		return nil, fmt.Errorf("error scanning artist stats: %v", err)
	}

	if firstRelease.Valid {
		stats.FirstRelease = firstRelease.Time
	}
	if latestRelease.Valid {
		stats.LatestRelease = latestRelease.Time
	}

	return &stats, nil
}

// AlbumStats represents statistics for an album
type AlbumStats struct {
	AlbumID         int     `json:"album_id"`
	Title           string  `json:"title"`
	TrackCount      int     `json:"track_count"`
	TotalDuration   int     `json:"total_duration_seconds"`
	DurationMinutes float64 `json:"duration_minutes"`
}

// GetAlbumStats retrieves album statistics using the album_stats table
func GetAlbumStats(albumID int) (*AlbumStats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var stats AlbumStats
	query := `
		SELECT 
			a.id,
			a.title,
			COALESCE(ast.track_count, 0),
			COALESCE(ast.total_duration, 0),
			ROUND(COALESCE(ast.total_duration, 0) / 60.0, 2)
		FROM albums a
		LEFT JOIN album_stats ast ON a.id = ast.album_id
		WHERE a.id = ?
	`

	err := MySQL.QueryRowContext(ctx, query, albumID).Scan(
		&stats.AlbumID,
		&stats.Title,
		&stats.TrackCount,
		&stats.TotalDuration,
		&stats.DurationMinutes,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting album stats: %v", err)
	}

	return &stats, nil
}

// GetAlbumDuration uses the SQL function to get album duration
func GetAlbumDuration(albumID int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var duration int
	err := MySQL.QueryRowContext(ctx, "SELECT get_album_duration(?)", albumID).Scan(&duration)
	if err != nil {
		return 0, fmt.Errorf("error getting album duration: %v", err)
	}

	return duration, nil
}
