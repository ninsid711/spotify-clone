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
	"go.mongodb.org/mongo-driver/mongo/options"
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

// InitMongoDB initializes MongoDB connection for user profiles
func InitMongoDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_URI"))
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("error connecting to MongoDB: %v", err)
	}

	// Test the connection
	if err = client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("error pinging MongoDB: %v", err)
	}

	MongoDB = client.Database(os.Getenv("MONGODB_DATABASE"))
	log.Println("✅ MongoDB connected successfully")

	// Create indexes
	if err := createMongoIndexes(); err != nil {
		return fmt.Errorf("error creating MongoDB indexes: %v", err)
	}

	return nil
}

// InitNeo4j initializes Neo4j connection for recommendations
func InitNeo4j() error {
	uri := os.Getenv("NEO4J_URI")
	username := os.Getenv("NEO4J_USERNAME")
	password := os.Getenv("NEO4J_PASSWORD")

	var err error
	Neo4j, err = neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return fmt.Errorf("error creating Neo4j driver: %v", err)
	}

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = Neo4j.VerifyConnectivity(ctx); err != nil {
		return fmt.Errorf("error verifying Neo4j connectivity: %v", err)
	}

	log.Println("✅ Neo4j connected successfully")

	// Initialize constraints and indexes
	if err := initNeo4jSchema(); err != nil {
		return fmt.Errorf("error initializing Neo4j schema: %v", err)
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
		// Minimal plays table to support trigger-based stats updates
		`CREATE TABLE IF NOT EXISTS plays (
		id INT AUTO_INCREMENT PRIMARY KEY,
		user_id VARCHAR(64) NULL,
		track_id INT NOT NULL,
		played_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
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

// createMongoIndexes creates indexes for MongoDB collections
func createMongoIndexes() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Users collection indexes
	usersCollection := MongoDB.Collection("users")
	_, err := usersCollection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    map[string]interface{}{"email": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    map[string]interface{}{"username": 1},
			Options: options.Index().SetUnique(true),
		},
	})
	if err != nil {
		return err
	}

	// Playlists collection indexes
	playlistsCollection := MongoDB.Collection("playlists")
	_, err = playlistsCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: map[string]interface{}{"user_id": 1},
	})
	if err != nil {
		return err
	}

	log.Println("✅ MongoDB indexes created")
	return nil
}

// initNeo4jSchema creates constraints and indexes in Neo4j
func initNeo4jSchema() error {
	ctx := context.Background()
	session := Neo4j.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	queries := []string{
		"CREATE CONSTRAINT IF NOT EXISTS FOR (u:User) REQUIRE u.id IS UNIQUE",
		"CREATE CONSTRAINT IF NOT EXISTS FOR (t:Track) REQUIRE t.id IS UNIQUE",
		"CREATE CONSTRAINT IF NOT EXISTS FOR (a:Artist) REQUIRE a.id IS UNIQUE",
		"CREATE CONSTRAINT IF NOT EXISTS FOR (g:Genre) REQUIRE g.name IS UNIQUE",
		"CREATE INDEX IF NOT EXISTS FOR (t:Track) ON (t.genre)",
	}

	for _, query := range queries {
		_, err := session.Run(ctx, query, nil)
		if err != nil {
			log.Printf("Warning: Could not execute query '%s': %v", query, err)
			// Continue even if constraint already exists
		}
	}

	log.Println("✅ Neo4j schema initialized")
	return nil
}

// Close closes all database connections
func Close() {
	if MySQL != nil {
		MySQL.Close()
		log.Println("MySQL connection closed")
	}
	if MongoDB != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		MongoDB.Client().Disconnect(ctx)
		log.Println("MongoDB connection closed")
	}
	if Neo4j != nil {
		ctx := context.Background()
		Neo4j.Close(ctx)
		log.Println("Neo4j connection closed")
	}
}
