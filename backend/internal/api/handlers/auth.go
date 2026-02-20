package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/social-media-lead/backend/internal/api/middleware"
	"github.com/social-media-lead/backend/internal/models"
	"github.com/social-media-lead/backend/internal/store"
	"golang.org/x/crypto/bcrypt"
)

// AuthHandler handles authentication endpoints.
type AuthHandler struct {
	Store     store.Store
	JWTSecret string
}

// SignupRequest is the expected body for user registration.
type SignupRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
	FullName    string `json:"full_name" binding:"required"`
	CompanyName string `json:"company_name"`
}

// LoginRequest is the expected body for user login.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Signup creates a new user account.
func (h *AuthHandler) Signup(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password"})
		return
	}

	user := &models.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FullName:     req.FullName,
		CompanyName:  req.CompanyName,
		Plan:         "starter",
		IsActive:     true,
	}

	if err := h.Store.CreateUser(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists or database error"})
		return
	}

	token, err := middleware.GenerateToken(user.ID, user.Email, h.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Account created successfully",
		"token":   token,
		"user": gin.H{
			"id":           user.ID,
			"email":        user.Email,
			"full_name":    user.FullName,
			"company_name": user.CompanyName,
			"plan":         user.Plan,
		},
	})
}

// Login authenticates a user and returns a JWT.
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Store.GetUserByEmail(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{"error": "Account is deactivated"})
		return
	}

	token, err := middleware.GenerateToken(user.ID, user.Email, h.JWTSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user": gin.H{
			"id":           user.ID,
			"email":        user.Email,
			"full_name":    user.FullName,
			"company_name": user.CompanyName,
			"plan":         user.Plan,
		},
	})
}

// Me returns the current authenticated user's profile.
func (h *AuthHandler) Me(c *gin.Context) {
	userID, _ := c.Get("user_id")
	user, err := h.Store.GetUserByID(c.Request.Context(), userID.(int64))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":           user.ID,
			"email":        user.Email,
			"full_name":    user.FullName,
			"company_name": user.CompanyName,
			"plan":         user.Plan,
			"created_at":   user.CreatedAt,
		},
	})
}

// UpdateProfileRequest is the expected body for profile updates.
type UpdateProfileRequest struct {
	FullName    string `json:"full_name"`
	Email       string `json:"email"`
	CompanyName string `json:"company_name"`
}

// UpdateProfile updates the current user's profile info.
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get current user to fill in unchanged fields
	currentUser, err := h.Store.GetUserByID(c.Request.Context(), userID.(int64))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	fullName := currentUser.FullName
	email := currentUser.Email
	company := currentUser.CompanyName

	if req.FullName != "" {
		fullName = req.FullName
	}
	if req.Email != "" {
		email = req.Email
	}
	if req.CompanyName != "" {
		company = req.CompanyName
	}

	updated, err := h.Store.UpdateUserProfile(c.Request.Context(), userID.(int64), fullName, email, company)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated",
		"user": gin.H{
			"id":           updated.ID,
			"email":        updated.Email,
			"full_name":    updated.FullName,
			"company_name": updated.CompanyName,
			"plan":         updated.Plan,
			"created_at":   updated.CreatedAt,
		},
	})
}

// ChangePasswordRequest is the expected body for password changes.
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

// ChangePassword changes the current user's password.
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch current user to verify old password
	user, err := h.Store.GetUserByID(c.Request.Context(), userID.(int64))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Verify current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.CurrentPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Current password is incorrect"})
		return
	}

	// Hash new password
	newHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process new password"})
		return
	}

	if err := h.Store.UpdateUserPassword(c.Request.Context(), userID.(int64), string(newHash)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}
