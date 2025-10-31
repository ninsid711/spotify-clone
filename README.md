# Spotify Clone - Complete Project Documentation

> A full-stack music streaming application built with Go, MySQL, MongoDB, and Neo4j

**Date:** October 31, 2025

---

## Table of Contents

1. [Project Overview](#project-overview)
2. [Architecture](#architecture)
3. [Database Schema](#database-schema)
4. [API Endpoints](#api-endpoints)
5. [Database Features](#database-features)
6. [Project Structure](#project-structure)
7. [Setup & Installation](#setup--installation)
8. [Testing Guide](#testing-guide)

---

## Project Overview

### Description

A Spotify-inspired music streaming platform that demonstrates the use of multiple database systems working together. The application provides music catalog management, user authentication, playlist creation, and personalized recommendations.

### Key Features

- 🎵 **Music Catalog** - Browse tracks, artists, and albums
- 👤 **User Management** - Authentication and profile management
- 📋 **Playlists** - Create and manage custom playlists
- 🎯 **Recommendations** - Personalized track suggestions
- 📊 **Statistics** - Artist and album analytics
- 🔍 **Search** - Find tracks, artists, and albums
- ⚡ **Database Triggers** - Automatic statistics maintenance
- ✅ **Data Validation** - Stored procedures for data integrity

### Technology Stack

| Component | Technology |
|-----------|------------|
| **Backend** | Go (Golang) 1.21+ |
| **Web Framework** | Gin |
| **Music Catalog DB** | MySQL 8.0+ |
| **User Data DB** | MongoDB 6.0+ |
| **Recommendations** | Neo4j 5.0+ |
| **Authentication** | JWT (JSON Web Tokens) |

---

## Architecture

### System Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                           CLIENT LAYER                           │
│                  (Browser / Mobile App / Postman)                │
└────────────────────────────────┬────────────────────────────────┘
                                 │
                                 │ HTTP/REST
                                 ▼
┌─────────────────────────────────────────────────────────────────┐
│                      APPLICATION LAYER (Go)                      │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │                    Gin Web Framework                      │  │
│  │  • Routing • Middleware • JSON Handling • CORS           │  │
│  └──────────────────────────────────────────────────────────┘  │
│                                                                  │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │                       Handlers                            │  │
│  │  • auth.go           - Authentication & Registration      │  │
│  │  • tracks.go         - Track management                   │  │
│  │  • playlists.go      - Playlist CRUD operations          │  │
│  │  • recommendations.go - Recommendation engine            │  │
│  │  • database_features.go - DB procedures & functions      │  │
│  └──────────────────────────────────────────────────────────┘  │
│                                                                  │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │                      Middleware                           │  │
│  │  • auth.go - JWT Authentication & Authorization          │  │
│  │  • CORS - Cross-Origin Resource Sharing                  │  │
│  └──────────────────────────────────────────────────────────┘  │
│                                                                  │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │                   Database Layer                          │  │
│  │  • db.go      - Connection management                    │  │
│  │  • triggers.go - Triggers, procedures, functions         │  │
│  └──────────────────────────────────────────────────────────┘  │
└────────────────────────────────┬────────────────────────────────┘
                                 │
                    ┌────────────┼────────────┐
                    │            │            │
                    ▼            ▼            ▼
         ┌─────────────┐ ┌─────────────┐ ┌──────────┐
         │    MySQL    │ │   MongoDB   │ │  Neo4j   │
         │             │ │             │ │          │
         │  • Tracks   │ │  • Users    │ │  • Graph │
         │  • Artists  │ │  • Profiles │ │  • Links │
         │  • Albums   │ │  • Playlists│ │  • Recs  │
         │  • Genres   │ │  • History  │ │          │
         └─────────────┘ └─────────────┘ └──────────┘
```

### Database Responsibilities

| Database | Purpose | Data Types |
|----------|---------|------------|
| **MySQL** | Music Catalog | Tracks, Artists, Albums, Genres, Statistics |
| **MongoDB** | User Data | Users, Profiles, Playlists, Listening History |
| **Neo4j** | Relationships | User-Track relationships, Recommendations |

### Data Flow Example

**Adding a Track:**
```
User Request
    ↓
POST /api/v1/tracks/add
    ↓
Handler: AddTrackWithValidation()
    ↓
Database Layer: Calls stored procedure
    ↓
MySQL: add_track() procedure validates data
    ↓
INSERT INTO tracks
    ↓
Trigger: after_track_insert FIRES
    ├─→ UPDATE album_stats (track_count +1)
    └─→ INSERT track_stats
    ↓
Response: {success: true, track_id: 31}
```

---

## Database Schema

### ER Diagram

```
┌─────────────────────┐
│      ARTISTS        │
│─────────────────────│
│ PK  id              │
│     name            │
│     bio             │
│     image_url       │
│     created_at      │
└──────────┬──────────┘
           │
           │ 1
           │
           │ N
┌──────────┴──────────┐         ┌─────────────────────┐
│      ALBUMS         │         │      GENRES         │
│─────────────────────│         │─────────────────────│
│ PK  id              │         │ PK  id              │
│ FK  artist_id       │         │     name            │
│     title           │         └─────────────────────┘
│     release_date    │
│     cover_url       │
│     created_at      │
└──────────┬──────────┘
           │
           │ 1
           │
           │ N
┌──────────┴──────────┐
│      TRACKS         │◄──────────┐
│─────────────────────│           │
│ PK  id              │           │
│ FK  artist_id       │           │
│ FK  album_id        │           │
│     title           │           │
│     duration        │           │
│     genre           │           │
│     release_date    │           │
│     file_url        │           │
│     cover_url       │           │
│     created_at      │           │
└──────────┬──────────┘           │
           │                      │
           │ 1                    │ 1
           │                      │
           │ 1                    │ 1
┌──────────┴──────────┐ ┌────────┴─────────┐
│    TRACK_STATS      │ │   ALBUM_STATS    │
│─────────────────────│ │──────────────────│
│ PK,FK track_id      │ │ PK,FK album_id   │
│     play_count      │ │     track_count  │
│     last_played     │ │  total_duration  │
└─────────────────────┘ │  last_updated    │
                        └──────────────────┘

┌────────────────────────────────────────────────────────┐
│                    MONGODB                             │
└────────────────────────────────────────────────────────┘

┌─────────────────────┐         ┌─────────────────────┐
│       USERS         │         │     PLAYLISTS       │
│─────────────────────│         │─────────────────────│
│ PK  _id             │         │ PK  _id             │
│     email (unique)  │         │ FK  user_id         │
│     username        │         │     name            │
│     password_hash   │         │     description     │
│     display_name    │         │     track_ids[]     │
│     profile_pic_url │         │     is_public       │
│     preferences {}  │         │     cover_url       │
│  listening_history[]│         │     created_at      │
│  favorite_artists[] │         │     updated_at      │
│  favorite_genres[]  │         └─────────────────────┘
│     created_at      │
│     updated_at      │
└─────────────────────┘

┌────────────────────────────────────────────────────────┐
│                     NEO4J                              │
└────────────────────────────────────────────────────────┘

    ┌─────────┐
    │  User   │
    └────┬────┘
         │
         │ :PLAYED
         │ {played_at, duration}
         ↓
    ┌─────────┐        ┌─────────┐
    │  Track  │───────→│  Genre  │
    └────┬────┘:HAS_GENRE└────────┘
         │
         │ :SIMILAR_TO
         │ {similarity_score}
         ↓
    ┌─────────┐
    │  Track  │
    └─────────┘
```

### MySQL Schema

#### Tables

**1. artists**
```sql
CREATE TABLE artists (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    bio TEXT,
    image_url VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_name (name)
);
```

**2. albums**
```sql
CREATE TABLE albums (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    artist_id INT NOT NULL,
    release_date DATE,
    cover_url VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (artist_id) REFERENCES artists(id) ON DELETE CASCADE,
    INDEX idx_artist (artist_id),
    INDEX idx_title (title)
);
```

**3. tracks**
```sql
CREATE TABLE tracks (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    artist_id INT NOT NULL,
    album_id INT NOT NULL,
    duration INT NOT NULL,              -- in seconds
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
);
```

**4. genres**
```sql
CREATE TABLE genres (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);
```

**5. album_stats** (Maintained by triggers)
```sql
CREATE TABLE album_stats (
    album_id INT PRIMARY KEY,
    track_count INT DEFAULT 0,
    total_duration INT DEFAULT 0,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (album_id) REFERENCES albums(id) ON DELETE CASCADE
);
```

**6. track_stats**
```sql
CREATE TABLE track_stats (
    track_id INT PRIMARY KEY,
    play_count INT DEFAULT 0,
    last_played TIMESTAMP NULL,
    FOREIGN KEY (track_id) REFERENCES tracks(id) ON DELETE CASCADE
);
```

### MongoDB Schema

**users collection**
```javascript
{
    _id: ObjectId,
    email: String (unique),
    username: String (unique),
    password: String (hashed),
    display_name: String,
    profile_picture_url: String,
    preferences: {
        theme: String,              // "light" or "dark"
        language: String,
        explicit_content: Boolean,
        preferred_genres: [String]
    },
    listening_history: [{
        track_id: Number,
        played_at: Date,
        duration: Number,           // seconds listened
        completed: Boolean
    }],
    favorite_artists: [Number],     // artist IDs from MySQL
    favorite_genres: [String],
    created_at: Date,
    updated_at: Date
}
```

**playlists collection**
```javascript
{
    _id: ObjectId,
    user_id: String,                // references users._id
    name: String,
    description: String,
    track_ids: [Number],            // track IDs from MySQL
    is_public: Boolean,
    cover_url: String,
    created_at: Date,
    updated_at: Date
}
```

### Neo4j Schema

**Nodes:**
- `(:User {user_id: String})`
- `(:Track {track_id: Number})`
- `(:Genre {name: String})`

**Relationships:**
- `(:User)-[:PLAYED {played_at: DateTime, duration: Number}]->(:Track)`
- `(:Track)-[:HAS_GENRE]->(:Genre)`
- `(:Track)-[:SIMILAR_TO {similarity_score: Float}]->(:Track)`

---

## API Endpoints

### Base URL
```
http://localhost:8080/api/v1
```

### Authentication Endpoints

#### Register User
```http
POST /api/v1/auth/register
```

**Request Body:**
```json
{
    "email": "user@example.com",
    "password": "password123",
    "username": "johndoe",
    "display_name": "John Doe",
    "genres": ["Pop", "Rock"]
}
```

**Response:**
```json
{
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
        "id": "507f1f77bcf86cd799439011",
        "email": "user@example.com",
        "username": "johndoe",
        "display_name": "John Doe"
    }
}
```

#### Login
```http
POST /api/v1/auth/login
```

**Request Body:**
```json
{
    "email": "user@example.com",
    "password": "password123"
}
```

**Response:**
```json
{
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
        "id": "507f1f77bcf86cd799439011",
        "email": "user@example.com",
        "username": "johndoe"
    }
}
```

---

### Track Endpoints

#### Get All Tracks
```http
GET /api/v1/tracks
```

**Query Parameters:**
- `limit` (optional) - Number of tracks to return (default: 20)
- `offset` (optional) - Pagination offset (default: 0)

**Response:**
```json
{
    "tracks": [
        {
            "id": 1,
            "title": "Blinding Lights",
            "artist_id": 1,
            "artist_name": "The Weeknd",
            "album_id": 1,
            "album_name": "After Hours",
            "duration": 200,
            "genre": "Pop",
            "release_date": "2020-03-20T00:00:00Z",
            "file_url": "https://audio.example.com/blinding-lights.mp3",
            "cover_url": "https://i.scdn.co/image/cover.jpg"
        }
    ]
}
```

#### Get Track by ID
```http
GET /api/v1/tracks/:id
```

**Response:**
```json
{
    "track": {
        "id": 1,
        "title": "Blinding Lights",
        "artist_name": "The Weeknd",
        "album_name": "After Hours",
        "duration": 200,
        "genre": "Pop"
    }
}
```

#### Add Track with Validation (Stored Procedure)
```http
POST /api/v1/tracks/add
```

**Request Body:**
```json
{
    "title": "New Song",
    "artist_id": 1,
    "album_id": 1,
    "duration": 240,
    "genre": "Pop",
    "release_date": "2025-10-31",
    "file_url": "https://audio.example.com/new-song.mp3",
    "cover_url": "https://i.scdn.co/image/cover.jpg"
}
```

**Response (Success):**
```json
{
    "success": true,
    "message": "SUCCESS: Track added successfully",
    "track_id": 31
}
```

**Response (Error):**
```json
{
    "success": false,
    "message": "ERROR: Album does not belong to the specified artist"
}
```

**Features:**
- ✅ Validates artist exists
- ✅ Validates album exists
- ✅ Validates album belongs to artist
- ✅ Triggers automatically update album_stats

#### Get Similar Tracks
```http
GET /api/v1/tracks/:id/similar
```

**Response:**
```json
{
    "similar_tracks": [
        {
            "id": 2,
            "title": "Save Your Tears",
            "artist_name": "The Weeknd",
            "similarity_score": 0.85
        }
    ]
}
```

---

### Artist Endpoints

#### Get All Artists
```http
GET /api/v1/artists
```

**Response:**
```json
{
    "artists": [
        {
            "id": 1,
            "name": "The Weeknd",
            "bio": "Canadian singer and songwriter...",
            "image_url": "https://i.scdn.co/image/artist.jpg"
        }
    ]
}
```

#### Get Artist by ID
```http
GET /api/v1/artists/:id
```

**Response:**
```json
{
    "artist": {
        "id": 1,
        "name": "The Weeknd",
        "bio": "Canadian singer...",
        "image_url": "https://..."
    }
}
```

#### Get Artist Statistics (Stored Procedure)
```http
GET /api/v1/artists/:id/stats
```

**Response:**
```json
{
    "success": true,
    "data": {
        "artist_id": 1,
        "artist_name": "The Weeknd",
        "total_albums": 1,
        "total_tracks": 3,
        "total_duration_seconds": 652,
        "total_duration_minutes": 10.87,
        "unique_genres": 2,
        "genres": "Pop, R&B",
        "first_release": "2020-03-20T00:00:00Z",
        "latest_release": "2020-03-20T00:00:00Z",
        "avg_plays_per_track": 0
    }
}
```

**Features:**
- Uses stored procedure `get_artist_stats()`
- Comprehensive analytics
- Performance optimized

---

### Album Endpoints

#### Get All Albums
```http
GET /api/v1/albums
```

**Response:**
```json
{
    "albums": [
        {
            "id": 1,
            "title": "After Hours",
            "artist_id": 1,
            "artist_name": "The Weeknd",
            "release_date": "2020-03-20T00:00:00Z",
            "cover_url": "https://..."
        }
    ]
}
```

#### Get Album Statistics (Trigger-maintained)
```http
GET /api/v1/albums/:id/stats
```

**Response:**
```json
{
    "success": true,
    "data": {
        "album_id": 1,
        "title": "After Hours",
        "track_count": 3,
        "total_duration_seconds": 652,
        "duration_minutes": 10.87
    }
}
```

**Features:**
- Data automatically maintained by triggers
- Fast query (no calculation needed)
- Always up-to-date

#### Get Album Duration (SQL Function)
```http
GET /api/v1/albums/:id/duration
```

**Response:**
```json
{
    "success": true,
    "album_id": 1,
    "duration_seconds": 652,
    "duration_minutes": 10.87,
    "duration_formatted": "10:52",
    "duration_display": {
        "minutes": 10,
        "seconds": 52
    }
}
```

**Features:**
- Uses SQL function `get_album_duration()`
- Calculates total track duration
- Multiple format options

---

### Playlist Endpoints (Protected)

All playlist endpoints require authentication. Include JWT token in header:
```
Authorization: Bearer <token>
```

#### Create Playlist
```http
POST /api/v1/playlists
```

**Request Body:**
```json
{
    "name": "My Awesome Playlist",
    "description": "Best tracks ever",
    "is_public": true
}
```

**Response:**
```json
{
    "playlist": {
        "id": "507f1f77bcf86cd799439012",
        "name": "My Awesome Playlist",
        "description": "Best tracks ever",
        "track_ids": [],
        "is_public": true
    }
}
```

#### Get User Playlists
```http
GET /api/v1/playlists
```

**Response:**
```json
{
    "playlists": [
        {
            "id": "507f1f77bcf86cd799439012",
            "name": "My Awesome Playlist",
            "track_ids": [1, 2, 3],
            "is_public": true
        }
    ]
}
```

#### Get Playlist by ID
```http
GET /api/v1/playlists/:id
```

**Response:**
```json
{
    "playlist": {
        "id": "507f1f77bcf86cd799439012",
        "name": "My Awesome Playlist",
        "tracks": [
            {
                "id": 1,
                "title": "Blinding Lights",
                "artist_name": "The Weeknd"
            }
        ]
    }
}
```

#### Update Playlist
```http
PUT /api/v1/playlists/:id
```

**Request Body:**
```json
{
    "name": "Updated Playlist Name",
    "description": "New description"
}
```

#### Delete Playlist
```http
DELETE /api/v1/playlists/:id
```

**Response:**
```json
{
    "message": "Playlist deleted successfully"
}
```

#### Add Track to Playlist
```http
POST /api/v1/playlists/:id/tracks
```

**Request Body:**
```json
{
    "track_id": 5
}
```

#### Remove Track from Playlist
```http
DELETE /api/v1/playlists/:id/tracks/:trackId
```

---

### Recommendation Endpoints

#### Get Personalized Recommendations (Protected)
```http
GET /api/v1/recommendations
```

**Response:**
```json
{
    "recommendations": [
        {
            "id": 5,
            "title": "Shape of You",
            "artist_name": "Ed Sheeran",
            "reason": "Based on your listening history"
        }
    ]
}
```

#### Get Trending Tracks
```http
GET /api/v1/recommendations/trending
```

**Response:**
```json
{
    "trending": [
        {
            "id": 1,
            "title": "Blinding Lights",
            "artist_name": "The Weeknd",
            "play_count": 1000
        }
    ]
}
```

#### Get Genre Recommendations
```http
GET /api/v1/recommendations/genre/:genre
```

**Response:**
```json
{
    "genre": "Pop",
    "tracks": [
        {
            "id": 1,
            "title": "Blinding Lights",
            "artist_name": "The Weeknd"
        }
    ]
}
```

---

### Search Endpoint

#### Search Tracks, Artists, Albums
```http
GET /api/v1/search?q=weeknd
```

**Response:**
```json
{
    "tracks": [...],
    "artists": [...],
    "albums": [...]
}
```

---

### User Profile Endpoints (Protected)

#### Get Profile
```http
GET /api/v1/profile
```

**Response:**
```json
{
    "user": {
        "id": "507f1f77bcf86cd799439011",
        "email": "user@example.com",
        "username": "johndoe",
        "display_name": "John Doe",
        "preferences": {
            "theme": "dark",
            "explicit_content": true
        }
    }
}
```

#### Update Preferences
```http
PUT /api/v1/profile/preferences
```

**Request Body:**
```json
{
    "theme": "dark",
    "language": "en",
    "preferred_genres": ["Pop", "Rock"]
}
```

---

### Play Tracking Endpoint (Protected)

#### Record Track Play
```http
POST /api/v1/tracks/:id/play
```

**Request Body:**
```json
{
    "duration": 180,
    "completed": true
}
```

**Features:**
- Updates listening history in MongoDB
- Creates relationships in Neo4j
- Used for recommendations

---

### Health Check

#### Check System Health
```http
GET /health
```

**Response:**
```json
{
    "status": "ok",
    "mysql": true,
    "mongodb": true,
    "neo4j": true
}
```

---

## Database Features

### Triggers (2)

#### 1. after_track_insert
**Purpose:** Automatically maintain album statistics when tracks are added.

**Triggered By:** `INSERT INTO tracks`

**Actions:**
- Updates `album_stats` table with incremented track count
- Adds track duration to album's total duration
- Initializes entry in `track_stats` table

**Example:**
```sql
-- Add a track
INSERT INTO tracks (title, artist_id, album_id, duration, ...)
VALUES ('New Song', 1, 1, 240, ...);

-- Trigger automatically:
-- 1. Updates album_stats: track_count +1, total_duration +240
-- 2. Creates track_stats entry with play_count = 0
```

#### 2. after_track_delete
**Purpose:** Automatically maintain album statistics when tracks are deleted.

**Triggered By:** `DELETE FROM tracks`

**Actions:**
- Decrements track count in `album_stats`
- Subtracts track duration from album's total duration
- Removes entry from `track_stats`

---

### Stored Procedures (2)

#### 1. add_track()
**Purpose:** Add tracks with comprehensive validation.

**Parameters:**
- IN: `p_title`, `p_artist_id`, `p_album_id`, `p_duration`, `p_genre`, `p_release_date`, `p_file_url`, `p_cover_url`
- OUT: `p_track_id`, `p_status`

**Validations:**
✅ Artist must exist
✅ Album must exist
✅ Album must belong to the specified artist

**Usage:**
```sql
CALL add_track(
    'New Song',
    1,                    -- artist_id
    1,                    -- album_id
    240,                  -- duration
    'Pop',
    '2025-10-31',
    'https://...',
    'https://...',
    @track_id,
    @status
);

SELECT @track_id, @status;
```

**API Endpoint:** `POST /api/v1/tracks/add`

#### 2. get_artist_stats()
**Purpose:** Retrieve comprehensive statistics for an artist.

**Parameters:**
- IN: `p_artist_id`

**Returns:**
- Total albums
- Total tracks
- Total duration (seconds and minutes)
- Unique genres count
- List of genres
- First and latest release dates
- Average plays per track

**Usage:**
```sql
CALL get_artist_stats(1);
```

**API Endpoint:** `GET /api/v1/artists/:id/stats`

---

### Functions (1)

#### get_album_duration()
**Purpose:** Calculate total duration of all tracks in an album.

**Parameters:**
- IN: `album_id_param` (INT)

**Returns:** Total duration in seconds (INT)

**Usage:**
```sql
SELECT get_album_duration(1) AS duration_seconds;

-- Use in queries
SELECT 
    id, 
    title, 
    get_album_duration(id) AS duration,
    ROUND(get_album_duration(id) / 60.0, 2) AS minutes
FROM albums;
```

**API Endpoint:** `GET /api/v1/albums/:id/duration`

---

### Benefits of Database Features

| Feature | Benefit |
|---------|---------|
| **Triggers** | Automatic data consistency, no application code needed |
| **Procedures** | Data validation at database level, single source of truth |
| **Functions** | Reusable calculations, can be used in complex queries |

---

## Project Structure

```
spotify-clone/
├── main.go                      # Application entry point, routes
├── go.mod                       # Go module dependencies
├── go.sum                       # Dependency checksums
├── .env                         # Environment variables (not in repo)
├── start.bat                    # Windows startup script
│
├── handlers/                    # HTTP request handlers
│   ├── auth.go                 # Registration & login
│   ├── tracks.go               # Track CRUD operations
│   ├── playlists.go            # Playlist management
│   ├── recommendations.go      # Recommendation engine
│   └── database_features.go    # DB procedures & functions
│
├── middleware/                  # HTTP middleware
│   └── auth.go                 # JWT authentication
│
├── database/                    # Database layer
│   ├── db.go                   # Connection management
│   └── triggers.go             # Triggers, procedures, functions
│
├── models/                      # Data models
│   └── models.go               # Struct definitions
│
├── utils/                       # Utility functions
│   └── jwt.go                  # JWT token generation/validation
│
└── seed/                        # Database seed data
    ├── seed_data.sql           # Sample music data
    └── triggers_procedures_functions.sql  # DB objects (reference)
```

### File Descriptions

#### main.go
- Application entry point
- Initializes database connections
- Sets up Gin router
- Defines all API routes
- Starts HTTP server

#### handlers/
**auth.go**
- User registration
- User login
- Password hashing
- JWT token generation

**tracks.go**
- Get all tracks
- Get track by ID
- Search tracks
- Get similar tracks (Neo4j)

**playlists.go**
- Create playlist
- Get user playlists
- Update playlist
- Delete playlist
- Add/remove tracks

**recommendations.go**
- Personalized recommendations
- Trending tracks
- Genre-based recommendations
- Uses Neo4j graph relationships

**database_features.go**
- Artist statistics (stored procedure)
- Album statistics (trigger-maintained)
- Album duration (SQL function)
- Add track with validation (stored procedure)

#### middleware/auth.go
- JWT token validation
- User authentication
- Protected route middleware

#### database/
**db.go**
- MySQL connection and initialization
- MongoDB connection and indexes
- Neo4j connection and schema
- Schema creation

**triggers.go**
- Creates utility tables (album_stats, track_stats)
- Creates triggers (after_track_insert, after_track_delete)
- Creates function (get_album_duration)
- Creates procedures (add_track, get_artist_stats)
- Initializes existing data

#### models/models.go
- Track, Artist, Album, Genre structs
- User, Playlist structs
- Request/Response models

#### utils/jwt.go
- Generate JWT tokens
- Validate JWT tokens
- Extract user info from tokens

---

## Setup & Installation

### Prerequisites

- **Go** 1.21 or higher
- **MySQL** 8.0 or higher
- **MongoDB** 6.0 or higher
- **Neo4j** 5.0 or higher

### Environment Variables

Create a `.env` file in the project root:

```env
# Server
PORT=8080

# MySQL
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_USER=root
MYSQL_PASSWORD=your_password
MYSQL_DATABASE=spotify_music

# MongoDB
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=spotify_users

# Neo4j
NEO4J_URI=bolt://localhost:7687
NEO4J_USERNAME=neo4j
NEO4J_PASSWORD=your_password

# JWT
JWT_SECRET=your_super_secret_key_change_this_in_production
```

### Installation Steps

1. **Clone the repository:**
```bash
git clone <repository-url>
cd spotify-clone
```

2. **Install Go dependencies:**
```bash
go mod download
```

3. **Start databases:**

**MySQL:**
```bash
# Using Docker
docker run -d --name mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=password mysql:8.0

# Create database
mysql -u root -p
CREATE DATABASE spotify_music;
```

**MongoDB:**
```bash
# Using Docker
docker run -d --name mongodb -p 27017:27017 mongo:6.0
```

**Neo4j:**
```bash
# Using Docker
docker run -d --name neo4j -p 7474:7474 -p 7687:7687 \
  -e NEO4J_AUTH=neo4j/password neo4j:5
```

4. **Configure environment:**
```bash
# Copy and edit .env file
cp .env.example .env
# Edit .env with your database credentials
```

5. **Build the application:**
```bash
go build -o spotify-clone.exe
```

6. **Run the application:**
```bash
.\spotify-clone.exe
# Or use the batch file
.\start.bat
```

7. **Verify installation:**
```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
    "status": "ok",
    "mysql": true,
    "mongodb": true,
    "neo4j": true
}
```

### Database Initialization

The application automatically:
- Creates MySQL tables
- Creates MongoDB indexes
- Creates Neo4j constraints
- Creates triggers, procedures, and functions
- Seeds sample data (optional)

To manually seed data:
```bash
mysql -u root -p spotify_music < seed/seed_data.sql
```

---

## Testing Guide

### Using curl (Windows CMD)

#### Test Health
```cmd
curl http://localhost:8080/health
```

#### Register User
```cmd
curl -X POST http://localhost:8080/api/v1/auth/register -H "Content-Type: application/json" -d "{\"email\":\"test@example.com\",\"password\":\"password123\",\"username\":\"testuser\",\"display_name\":\"Test User\",\"genres\":[\"Pop\"]}"
```

#### Login
```cmd
curl -X POST http://localhost:8080/api/v1/auth/login -H "Content-Type: application/json" -d "{\"email\":\"test@example.com\",\"password\":\"password123\"}"
```

#### Get Tracks
```cmd
curl http://localhost:8080/api/v1/tracks
```

#### Get Artist Stats (Stored Procedure)
```cmd
curl http://localhost:8080/api/v1/artists/1/stats
```

#### Get Album Duration (SQL Function)
```cmd
curl http://localhost:8080/api/v1/albums/1/duration
```

#### Get Album Stats (Trigger-maintained)
```cmd
curl http://localhost:8080/api/v1/albums/1/stats
```

#### Add Track with Validation (Stored Procedure)
```cmd
curl -X POST http://localhost:8080/api/v1/tracks/add -H "Content-Type: application/json" -d "{\"title\":\"Test Song\",\"artist_id\":1,\"album_id\":1,\"duration\":240,\"genre\":\"Pop\",\"release_date\":\"2025-10-31\",\"file_url\":\"https://audio.example.com/test.mp3\",\"cover_url\":\"https://i.scdn.co/image/test.jpg\"}"
```

### Testing Triggers

**Before adding track:**
```cmd
curl http://localhost:8080/api/v1/albums/1/stats
# Note the track_count and total_duration
```

**Add a track:**
```cmd
curl -X POST http://localhost:8080/api/v1/tracks/add -H "Content-Type: application/json" -d "{\"title\":\"Trigger Test\",\"artist_id\":1,\"album_id\":1,\"duration\":180,\"genre\":\"Pop\",\"release_date\":\"2025-10-31\",\"file_url\":\"https://test.mp3\",\"cover_url\":\"https://test.jpg\"}"
```

**After adding track:**
```cmd
curl http://localhost:8080/api/v1/albums/1/stats
# Verify track_count increased by 1 and total_duration increased by 180
```

This proves the trigger is working automatically!

### Using Postman

1. Import the following as a Postman collection
2. Set base URL: `http://localhost:8080/api/v1`
3. For protected endpoints, add header: `Authorization: Bearer <token>`

### Direct Database Testing

**MySQL:**
```bash
mysql -u root -p spotify_music

-- Test function
SELECT get_album_duration(1);

-- Test procedure
CALL get_artist_stats(1);

-- View trigger-maintained data
SELECT * FROM album_stats;

-- Test adding track with validation
CALL add_track('Test Song', 1, 1, 200, 'Pop', '2025-10-31', 'test.mp3', 'test.jpg', @id, @status);
SELECT @id, @status;
```

**MongoDB:**
```bash
mongosh spotify_users

-- View users
db.users.find()

-- View playlists
db.playlists.find()
```

**Neo4j:**
```cypher
// View all nodes
MATCH (n) RETURN n LIMIT 25;

// View track relationships
MATCH (u:User)-[r:PLAYED]->(t:Track)
RETURN u, r, t LIMIT 10;
```

---

## Summary

### Key Achievements

✅ **Multi-database Architecture** - MySQL, MongoDB, Neo4j working together
✅ **RESTful API** - 25+ endpoints with proper HTTP methods
✅ **Authentication** - JWT-based user authentication
✅ **Database Features** - 2 triggers, 2 procedures, 1 function
✅ **Automatic Updates** - Triggers maintain statistics
✅ **Data Validation** - Stored procedures validate data
✅ **Graph Recommendations** - Neo4j-powered suggestions
✅ **Production Ready** - Error handling, logging, CORS

### Technology Highlights

| Component | Implementation |
|-----------|----------------|
| **Web Framework** | Gin (lightweight, fast) |
| **Catalog DB** | MySQL with triggers & procedures |
| **User DB** | MongoDB with flexible schema |
| **Graph DB** | Neo4j for recommendations |
| **Auth** | JWT tokens with middleware |
| **API Design** | RESTful, JSON responses |

### Performance Features

- Database connection pooling
- Indexed queries for fast searches
- Denormalized statistics (trigger-maintained)
- Efficient graph traversal (Neo4j)
- Minimal data transfer (pagination)

---

## Contact & Support

For questions or issues, please refer to the source code or contact the development team.

**Project Status:** Production Ready ✅

**Last Updated:** October 31, 2025

---

*Built with ❤️ using Go, MySQL, MongoDB, and Neo4j*

