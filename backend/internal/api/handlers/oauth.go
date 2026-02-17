package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/social-media-lead/backend/internal/api/middleware"
	"github.com/social-media-lead/backend/internal/models"
	"github.com/social-media-lead/backend/internal/store"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// OAuthHandler handles OAuth2 authentication flows.
type OAuthHandler struct {
	Store       *store.Storage
	JWTSecret   string
	GoogleCfg   *oauth2.Config
	FrontendURL string
}

// NewOAuthHandler creates a new OAuthHandler with Google config.
func NewOAuthHandler(storage *store.Storage, jwtSecret, googleClientID, googleClientSecret, redirectURL, frontendURL string) *OAuthHandler {
	var googleCfg *oauth2.Config
	if googleClientID != "" && googleClientSecret != "" {
		googleCfg = &oauth2.Config{
			ClientID:     googleClientID,
			ClientSecret: googleClientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"openid", "email", "profile"},
			Endpoint:     google.Endpoint,
		}
	}

	return &OAuthHandler{
		Store:       storage,
		JWTSecret:   jwtSecret,
		GoogleCfg:   googleCfg,
		FrontendURL: frontendURL,
	}
}

// GoogleLogin redirects the user to Google's consent screen.
func (h *OAuthHandler) GoogleLogin(c *gin.Context) {
	if h.GoogleCfg == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Google OAuth not configured"})
		return
	}

	state := generateState()
	c.SetCookie("oauth_state", state, 600, "/", "", false, true)

	url := h.GoogleCfg.AuthCodeURL(state, oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GoogleCallback handles the OAuth2 callback from Google.
func (h *OAuthHandler) GoogleCallback(c *gin.Context) {
	if h.GoogleCfg == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Google OAuth not configured"})
		return
	}

	// Verify state
	state := c.Query("state")
	savedState, err := c.Cookie("oauth_state")
	if err != nil || state != savedState {
		c.Redirect(http.StatusTemporaryRedirect, h.FrontendURL+"/login?error=invalid_state")
		return
	}

	// Exchange code for token
	code := c.Query("code")
	if code == "" {
		c.Redirect(http.StatusTemporaryRedirect, h.FrontendURL+"/login?error=no_code")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 15*time.Second)
	defer cancel()

	token, err := h.GoogleCfg.Exchange(ctx, code)
	if err != nil {
		log.Printf("[OAuth] Token exchange failed: %v", err)
		c.Redirect(http.StatusTemporaryRedirect, h.FrontendURL+"/login?error=exchange_failed")
		return
	}

	// Fetch user info from Google
	googleUser, err := fetchGoogleUserInfo(ctx, token.AccessToken)
	if err != nil {
		log.Printf("[OAuth] Failed to fetch user info: %v", err)
		c.Redirect(http.StatusTemporaryRedirect, h.FrontendURL+"/login?error=userinfo_failed")
		return
	}

	// Find or create user
	user, err := h.Store.GetOrCreateOAuthUser(ctx, &models.User{
		Email:      googleUser.Email,
		FullName:   googleUser.Name,
		GoogleID:   googleUser.ID,
		AvatarURL:  googleUser.Picture,
		Plan:       "starter",
		IsActive:   true,
	})
	if err != nil {
		log.Printf("[OAuth] Failed to upsert user: %v", err)
		c.Redirect(http.StatusTemporaryRedirect, h.FrontendURL+"/login?error=db_error")
		return
	}

	// Generate JWT
	jwt, err := middleware.GenerateToken(user.ID, user.Email, h.JWTSecret)
	if err != nil {
		log.Printf("[OAuth] JWT generation failed: %v", err)
		c.Redirect(http.StatusTemporaryRedirect, h.FrontendURL+"/login?error=token_failed")
		return
	}

	// Build user JSON for frontend
	userJSON := fmt.Sprintf(`{"id":%d,"email":"%s","full_name":"%s","company_name":"%s","plan":"%s","avatar_url":"%s"}`,
		user.ID, user.Email, user.FullName, user.CompanyName, user.Plan, user.AvatarURL)

	// Redirect to frontend with token
	redirectURL := fmt.Sprintf("%s/oauth/callback?token=%s&user=%s",
		h.FrontendURL, jwt, base64.URLEncoding.EncodeToString([]byte(userJSON)))
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// GoogleUserInfo is the response from Google's userinfo endpoint.
type GoogleUserInfo struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func fetchGoogleUserInfo(ctx context.Context, accessToken string) (*GoogleUserInfo, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("google API error %d: %s", resp.StatusCode, string(body))
	}

	var user GoogleUserInfo
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func generateState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
