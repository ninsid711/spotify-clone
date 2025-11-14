package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	MySQL   *sql.DB
	MongoDB *mongo.Database
	Neo4j   neo4j.DriverWithContext
)

// InitMySQL initializes MySQL connection for music catalog
func InitMySQL() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"),
	)

	var err error
	MySQL, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("error opening MySQL connection: %v", err)
	}

	// Set connection pool settings
	MySQL.SetMaxOpenConns(25)
	MySQL.SetMaxIdleConns(5)
	MySQL.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = MySQL.PingContext(ctx); err != nil {
		return fmt.Errorf("error pinging MySQL: %v", err)
	}

	log.Println("✅ MySQL connected successfully")

	// Initialize schema
	if err := initMySQLSchema(); err != nil {
		return fmt.Errorf("error initializing MySQL schema: %v", err)
	}

	// Initialize triggers, procedures, and functions
	if err := InitTriggersProceduresFunctions(); err != nil {
		log.Printf("⚠️  Warning: Could not initialize triggers/procedures/functions: %v", err)
		// Don't fail the entire initialization if this fails
	}

	return nil
}

// initMySQLSchema creates tables if they don't exist
func initMySQLSchema() error {
	schemas := []string{
		`CREATE TABLE IF NOT EXISTS artists (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			bio TEXT,
			image_url VARCHAR(500),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			INDEX idx_name (name)
		);`,
		`CREATE TABLE IF NOT EXISTS albums (
			id INT AUTO_INCREMENT PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			artist_id INT NOT NULL,
			release_date DATE,
			cover_url VARCHAR(500),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (artist_id) REFERENCES artists(id) ON DELETE CASCADE,
			INDEX idx_artist (artist_id),
			INDEX idx_title (title)
		);`,
		`CREATE TABLE IF NOT EXISTS tracks (
			id INT AUTO_INCREMENT PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			artist_id INT NOT NULL,
			album_id INT NOT NULL,
			duration INT NOT NULL,
			genre VARCHAR(100),
			release_date DATE,
			file_url VARCHAR(500),
			cover_url VARCHAR(500),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (artist_id) REFERENCES artists(id) ON DELETE CASCADE,
			FOREIGN KEY (album_id) REFERENCES albums(id) ON DELETE CASCADE,
			INDEX idx_title (title),
			INDEX idx_artist (artist_id),
			INDEX idx_album (album_id),
			INDEX idx_genre (genre)
		);`,
		`CREATE TABLE IF NOT EXISTS genres (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(100) NOT NULL UNIQUE
		);`,
		`CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			email VARCHAR(255) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			username VARCHAR(100) NOT NULL UNIQUE,
			display_name VARCHAR(255) NOT NULL,
			profile_picture_url VARCHAR(500),
			theme VARCHAR(20) DEFAULT 'dark',
			language VARCHAR(10) DEFAULT 'en',
			explicit_content BOOLEAN DEFAULT TRUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			INDEX idx_email (email),
			INDEX idx_username (username)
		);`,
		`CREATE TABLE IF NOT EXISTS user_favorite_genres (
			id INT AUTO_INCREMENT PRIMARY KEY,
			user_id INT NOT NULL,
			genre VARCHAR(100) NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			UNIQUE KEY unique_user_genre (user_id, genre),
			INDEX idx_user (user_id)
		);`,
		`CREATE TABLE IF NOT EXISTS user_favorite_artists (
			id INT AUTO_INCREMENT PRIMARY KEY,
			user_id INT NOT NULL,
			artist_id INT NOT NULL,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (artist_id) REFERENCES artists(id) ON DELETE CASCADE,
			UNIQUE KEY unique_user_artist (user_id, artist_id),
			INDEX idx_user (user_id)
		);`,
		`CREATE TABLE IF NOT EXISTS playlists (
			id INT AUTO_INCREMENT PRIMARY KEY,
			user_id INT NOT NULL,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			is_public BOOLEAN DEFAULT TRUE,
			cover_url VARCHAR(500),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			INDEX idx_user (user_id),
			INDEX idx_name (name)
		);`,
		`CREATE TABLE IF NOT EXISTS playlist_tracks (
			id INT AUTO_INCREMENT PRIMARY KEY,
			playlist_id INT NOT NULL,
			track_id INT NOT NULL,
			position INT NOT NULL DEFAULT 0,
			added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (playlist_id) REFERENCES playlists(id) ON DELETE CASCADE,
			FOREIGN KEY (track_id) REFERENCES tracks(id) ON DELETE CASCADE,
			UNIQUE KEY unique_playlist_track (playlist_id, track_id),
			INDEX idx_playlist (playlist_id),
			INDEX idx_track (track_id)
		);`,
		`CREATE TABLE IF NOT EXISTS plays (
			id INT AUTO_INCREMENT PRIMARY KEY,
			user_id INT NULL,
			track_id INT NOT NULL,
			played_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			duration_played INT DEFAULT 0,
			completed BOOLEAN DEFAULT FALSE,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
			FOREIGN KEY (track_id) REFERENCES tracks(id) ON DELETE CASCADE,
			INDEX idx_user (user_id),
			INDEX idx_track (track_id),
			INDEX idx_played_at (played_at)
		);`,
	}

	for _, schema := range schemas {
		_, err := MySQL.Exec(schema)
		if err != nil {
			return err
		}
	}

	log.Println("✅ MySQL schema initialized")
	return nil
}

// Close closes all database connections
func Close() {
	if MySQL != nil {
		MySQL.Close()
		log.Println("MySQL connection closed")
	}
}
