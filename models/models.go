package models

import (
	"time"
)

// MySQL Models - Music Catalog

type Track struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	ArtistID    int       `json:"artist_id"`
	ArtistName  string    `json:"artist_name,omitempty"`
	AlbumID     int       `json:"album_id"`
	AlbumName   string    `json:"album_name,omitempty"`
	Duration    int       `json:"duration"` // in seconds
	Genre       string    `json:"genre"`
	ReleaseDate time.Time `json:"release_date"`
	FileURL     string    `json:"file_url"`
	CoverURL    string    `json:"cover_url"`
	CreatedAt   time.Time `json:"created_at"`
}

type Artist struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Bio       string    `json:"bio"`
	ImageURL  string    `json:"image_url"`
	CreatedAt time.Time `json:"created_at"`
}

type Album struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	ArtistID    int       `json:"artist_id"`
	ArtistName  string    `json:"artist_name,omitempty"`
	ReleaseDate time.Time `json:"release_date"`
	CoverURL    string    `json:"cover_url"`
	CreatedAt   time.Time `json:"created_at"`
}

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// MySQL Models - All data in MySQL

type User struct {
	ID                int       `json:"id"`
	Email             string    `json:"email"`
	Password          string    `json:"-"` // Never expose in JSON
	Username          string    `json:"username"`
	DisplayName       string    `json:"display_name"`
	ProfilePictureURL string    `json:"profile_picture_url"`
	Theme             string    `json:"theme"`
	Language          string    `json:"language"`
	ExplicitContent   bool      `json:"explicit_content"`
	FavoriteGenres    []string  `json:"favorite_genres"`
	FavoriteArtists   []int     `json:"favorite_artists"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type UserPreferences struct {
	Theme           string   `json:"theme"` // light, dark
	Language        string   `json:"language"`
	ExplicitContent bool     `json:"explicit_content"`
	PreferredGenres []string `json:"preferred_genres"`
}

type ListeningHistory struct {
	TrackID        int       `json:"track_id"`
	PlayedAt       time.Time `json:"played_at"`
	DurationPlayed int       `json:"duration_played"` // How long they listened in seconds
	Completed      bool      `json:"completed"`
}

type Playlist struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	TrackIDs    []int     `json:"track_ids"`
	IsPublic    bool      `json:"is_public"`
	CoverURL    string    `json:"cover_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Request/Response Models

type RegisterRequest struct {
	Email       string   `json:"email" binding:"required,email"`
	Password    string   `json:"password" binding:"required,min=6"`
	Username    string   `json:"username" binding:"required,min=3"`
	DisplayName string   `json:"display_name" binding:"required"`
	Genres      []string `json:"genres"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type CreatePlaylistRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
}

type AddTrackToPlaylistRequest struct {
	TrackID int `json:"track_id" binding:"required"`
}

type SearchResponse struct {
	Tracks  []Track  `json:"tracks"`
	Artists []Artist `json:"artists"`
	Albums  []Album  `json:"albums"`
}

type RecommendationRequest struct {
	Limit int `json:"limit"`
}

type RecommendationResponse struct {
	Tracks []Track `json:"tracks"`
	Reason string  `json:"reason"`
}
