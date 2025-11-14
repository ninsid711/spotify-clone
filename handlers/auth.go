package handlers

import (
	"net/http"
	"spotify-clone/database"
	"spotify-clone/models"
	"spotify-clone/utils"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Register creates a new user account
func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists
	var existingID int
	err := database.MySQL.QueryRow("SELECT id FROM users WHERE email = ? OR username = ?", req.Email, req.Username).Scan(&existingID)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User with this email or username already exists"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create user
	result, err := database.MySQL.Exec(`
		INSERT INTO users (email, password, username, display_name, theme, language, explicit_content)
		VALUES (?, ?, ?, ?, 'dark', 'en', TRUE)`,
		req.Email, string(hashedPassword), req.Username, req.DisplayName)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	userID, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID"})
		return
	}

	// Add favorite genres
	if len(req.Genres) > 0 {
		for _, genre := range req.Genres {
			database.MySQL.Exec("INSERT INTO user_favorite_genres (user_id, genre) VALUES (?, ?)", userID, genre)
		}
	}

	// Generate token
	token, err := utils.GenerateToken(int(userID), req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Fetch the created user
	user := models.User{
		ID:              int(userID),
		Email:           req.Email,
		Username:        req.Username,
		DisplayName:     req.DisplayName,
		Theme:           "dark",
		Language:        "en",
		ExplicitContent: true,
		FavoriteGenres:  req.Genres,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	c.JSON(http.StatusCreated, models.AuthResponse{
		Token: token,
		User:  user,
	})
}

// Login authenticates a user
func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user
	var user models.User
	var hashedPassword string
	err := database.MySQL.QueryRow(`
		SELECT id, email, password, username, display_name, profile_picture_url, 
		       theme, language, explicit_content, created_at, updated_at
		FROM users WHERE email = ?`, req.Email).Scan(
		&user.ID, &user.Email, &hashedPassword, &user.Username, &user.DisplayName,
		&user.ProfilePictureURL, &user.Theme, &user.Language, &user.ExplicitContent,
		&user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Fetch favorite genres
	rows, _ := database.MySQL.Query("SELECT genre FROM user_favorite_genres WHERE user_id = ?", user.ID)
	defer rows.Close()
	for rows.Next() {
		var genre string
		rows.Scan(&genre)
		user.FavoriteGenres = append(user.FavoriteGenres, genre)
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, models.AuthResponse{
		Token: token,
		User:  user,
	})
}

// GetProfile returns the current user's profile
func GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	err := database.MySQL.QueryRow(`
		SELECT id, email, username, display_name, profile_picture_url, 
		       theme, language, explicit_content, created_at, updated_at
		FROM users WHERE id = ?`, userID).Scan(
		&user.ID, &user.Email, &user.Username, &user.DisplayName,
		&user.ProfilePictureURL, &user.Theme, &user.Language, &user.ExplicitContent,
		&user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Fetch favorite genres
	rows, _ := database.MySQL.Query("SELECT genre FROM user_favorite_genres WHERE user_id = ?", userID)
	defer rows.Close()
	for rows.Next() {
		var genre string
		rows.Scan(&genre)
		user.FavoriteGenres = append(user.FavoriteGenres, genre)
	}

	c.JSON(http.StatusOK, user)
}

// UpdatePreferences updates user preferences
func UpdatePreferences(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var prefs models.UserPreferences
	if err := c.ShouldBindJSON(&prefs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update user preferences
	_, err := database.MySQL.Exec(`
		UPDATE users 
		SET theme = ?, language = ?, explicit_content = ?, updated_at = NOW()
		WHERE id = ?`,
		prefs.Theme, prefs.Language, prefs.ExplicitContent, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update preferences"})
		return
	}

	// Update favorite genres
	database.MySQL.Exec("DELETE FROM user_favorite_genres WHERE user_id = ?", userID)
	for _, genre := range prefs.PreferredGenres {
		database.MySQL.Exec("INSERT INTO user_favorite_genres (user_id, genre) VALUES (?, ?)", userID, genre)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Preferences updated successfully"})
}
