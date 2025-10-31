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

// MongoDB Models - User Profiles and Preferences

type User struct {
	ID                string             `json:"id" bson:"_id,omitempty"`
	Email             string             `json:"email" bson:"email"`
	Password          string             `json:"-" bson:"password"` // Never expose in JSON
	Username          string             `json:"username" bson:"username"`
	DisplayName       string             `json:"display_name" bson:"display_name"`
	ProfilePictureURL string             `json:"profile_picture_url" bson:"profile_picture_url"`
	Preferences       UserPreferences    `json:"preferences" bson:"preferences"`
	ListeningHistory  []ListeningHistory `json:"listening_history" bson:"listening_history"`
	FavoriteArtists   []int              `json:"favorite_artists" bson:"favorite_artists"`
	FavoriteGenres    []string           `json:"favorite_genres" bson:"favorite_genres"`
	CreatedAt         time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at" bson:"updated_at"`
}

type UserPreferences struct {
	Theme           string   `json:"theme" bson:"theme"` // light, dark
	Language        string   `json:"language" bson:"language"`
	ExplicitContent bool     `json:"explicit_content" bson:"explicit_content"`
	PreferredGenres []string `json:"preferred_genres" bson:"preferred_genres"`
}

type ListeningHistory struct {
	TrackID   int       `json:"track_id" bson:"track_id"`
	PlayedAt  time.Time `json:"played_at" bson:"played_at"`
	Duration  int       `json:"duration" bson:"duration"` // How long they listened in seconds
	Completed bool      `json:"completed" bson:"completed"`
}

type Playlist struct {
	ID          string    `json:"id" bson:"_id,omitempty"`
	UserID      string    `json:"user_id" bson:"user_id"`
	Name        string    `json:"name" bson:"name"`
	Description string    `json:"description" bson:"description"`
	TrackIDs    []int     `json:"track_ids" bson:"track_ids"`
	IsPublic    bool      `json:"is_public" bson:"is_public"`
	CoverURL    string    `json:"cover_url" bson:"cover_url"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
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
