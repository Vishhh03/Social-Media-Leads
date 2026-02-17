package meta

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/social-media-lead/backend/internal/cache"
)

// TokenRefresher handles refreshing expired Meta API access tokens.
type TokenRefresher struct {
	AppID     string
	AppSecret string
	Redis     *cache.RedisClient
	HTTPClient *http.Client
}

// NewTokenRefresher creates a token refresher with Meta app credentials.
func NewTokenRefresher(appID, appSecret string, redis *cache.RedisClient) *TokenRefresher {
	return &TokenRefresher{
		AppID:      appID,
		AppSecret:  appSecret,
		Redis:      redis,
		HTTPClient: &http.Client{Timeout: 15 * time.Second},
	}
}

// TokenResponse is the Meta API token exchange response.
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"` // seconds
}

// ExchangeForLongLivedToken exchanges a short-lived token for a long-lived one (60 days).
// This should be called when a user first connects a channel.
func (tr *TokenRefresher) ExchangeForLongLivedToken(ctx context.Context, shortLivedToken string) (*TokenResponse, error) {
	params := url.Values{
		"grant_type":        {"fb_exchange_token"},
		"client_id":         {tr.AppID},
		"client_secret":     {tr.AppSecret},
		"fb_exchange_token": {shortLivedToken},
	}

	reqURL := fmt.Sprintf("%s/oauth/access_token?%s", graphAPIBase, params.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create token exchange request: %w", err)
	}

	resp, err := tr.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token exchange request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("token exchange failed (%d): %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	return &tokenResp, nil
}

// RefreshLongLivedToken refreshes an existing long-lived token before it expires.
// Long-lived tokens can be refreshed once per day, and only if not yet expired.
func (tr *TokenRefresher) RefreshLongLivedToken(ctx context.Context, currentToken string) (*TokenResponse, error) {
	params := url.Values{
		"grant_type":    {"fb_exchange_token"},
		"client_id":     {tr.AppID},
		"client_secret": {tr.AppSecret},
		"fb_exchange_token": {currentToken},
	}

	reqURL := fmt.Sprintf("%s/oauth/access_token?%s", graphAPIBase, params.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh request: %w", err)
	}

	resp, err := tr.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token refresh request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("token refresh failed (%d): %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token refresh response: %w", err)
	}

	return &tokenResp, nil
}

// GetValidToken returns a valid access token for a channel, refreshing if needed.
// It checks Redis cache first, then falls back to the stored token.
func (tr *TokenRefresher) GetValidToken(ctx context.Context, channelID int64, storedToken string, tokenExpiry time.Time) (string, error) {
	// 1. Check Redis cache
	if tr.Redis != nil {
		cached, err := tr.Redis.GetCachedAccessToken(ctx, channelID)
		if err == nil && cached != "" {
			return cached, nil
		}
	}

	// 2. If token hasn't expired yet, cache and return it
	if !tokenExpiry.IsZero() && time.Until(tokenExpiry) > 24*time.Hour {
		if tr.Redis != nil {
			_ = tr.Redis.CacheAccessToken(ctx, channelID, storedToken, time.Until(tokenExpiry)-1*time.Hour)
		}
		return storedToken, nil
	}

	// 3. Token is expired or expiring soon — try to refresh
	if tr.AppID != "" && tr.AppSecret != "" && storedToken != "" {
		tokenResp, err := tr.RefreshLongLivedToken(ctx, storedToken)
		if err != nil {
			log.Printf("[TokenRefresh] Failed to refresh token for channel #%d: %v", channelID, err)
			// Fall back to stored token — it might still work
			return storedToken, nil
		}

		expiry := time.Duration(tokenResp.ExpiresIn) * time.Second
		if tr.Redis != nil {
			_ = tr.Redis.CacheAccessToken(ctx, channelID, tokenResp.AccessToken, expiry-1*time.Hour)
		}

		log.Printf("[TokenRefresh] ✅ Refreshed token for channel #%d (expires in %s)", channelID, expiry)
		return tokenResp.AccessToken, nil
	}

	// 4. No refresh possible — return stored token as-is
	return storedToken, nil
}
