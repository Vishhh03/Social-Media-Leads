package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient wraps the Redis connection for caching and rate limiting.
type RedisClient struct {
	Client *redis.Client
}

// New creates a new Redis client and verifies the connection.
func New(host, port, password string) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", host, port),
		Password:     password,
		DB:           0,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     20,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("unable to connect to Redis: %w", err)
	}

	return &RedisClient{Client: client}, nil
}

// Close closes the Redis connection.
func (r *RedisClient) Close() error {
	return r.Client.Close()
}

// ---- Caching ----

// Set stores a value with an expiration time.
func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.Client.Set(ctx, key, data, expiration).Err()
}

// Get retrieves a cached value and unmarshals it into dest.
func (r *RedisClient) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := r.Client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

// Delete removes a key from the cache.
func (r *RedisClient) Delete(ctx context.Context, key string) error {
	return r.Client.Del(ctx, key).Err()
}

// ---- Rate Limiting (Sliding Window) ----

// RateLimitResult contains the result of a rate limit check.
type RateLimitResult struct {
	Allowed   bool
	Remaining int64
	RetryAfter time.Duration
}

// CheckRateLimit uses a sliding window counter in Redis to rate-limit actions.
// key: action identifier (e.g., "meta_api:user:123")
// limit: max allowed requests in the window
// window: time window duration
func (r *RedisClient) CheckRateLimit(ctx context.Context, key string, limit int64, window time.Duration) (*RateLimitResult, error) {
	now := time.Now().UnixMilli()
	windowStart := now - window.Milliseconds()

	pipe := r.Client.Pipeline()

	// Remove old entries outside the window
	pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", windowStart))
	// Add current request
	pipe.ZAdd(ctx, key, redis.Z{Score: float64(now), Member: fmt.Sprintf("%d", now)})
	// Count entries in window
	countCmd := pipe.ZCard(ctx, key)
	// Set expiry on the key
	pipe.Expire(ctx, key, window)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	count := countCmd.Val()
	if count > limit {
		// Over limit — remove the entry we just added
		r.Client.ZRemRangeByRank(ctx, key, -1, -1)
		return &RateLimitResult{
			Allowed:    false,
			Remaining:  0,
			RetryAfter: window,
		}, nil
	}

	return &RateLimitResult{
		Allowed:   true,
		Remaining: limit - count,
	}, nil
}

// ---- Session / Token Caching ----

// CacheAccessToken stores a Meta access token in Redis with its expiry.
func (r *RedisClient) CacheAccessToken(ctx context.Context, channelID int64, token string, expiry time.Duration) error {
	key := fmt.Sprintf("token:channel:%d", channelID)
	return r.Client.Set(ctx, key, token, expiry).Err()
}

// GetCachedAccessToken retrieves a cached Meta access token.
func (r *RedisClient) GetCachedAccessToken(ctx context.Context, channelID int64) (string, error) {
	key := fmt.Sprintf("token:channel:%d", channelID)
	return r.Client.Get(ctx, key).Result()
}

// InvalidateAccessToken removes a cached token (e.g., after refresh or disconnect).
func (r *RedisClient) InvalidateAccessToken(ctx context.Context, channelID int64) error {
	key := fmt.Sprintf("token:channel:%d", channelID)
	return r.Client.Del(ctx, key).Err()
}

// ---- Broadcast Deduplication ----

// MarkBroadcastSent records that a broadcast was sent to a contact (prevents double-sends).
func (r *RedisClient) MarkBroadcastSent(ctx context.Context, broadcastID, contactID int64) error {
	key := fmt.Sprintf("broadcast:%d:sent", broadcastID)
	return r.Client.SAdd(ctx, key, contactID).Err()
}

// WasBroadcastSent checks if a broadcast was already sent to a contact.
func (r *RedisClient) WasBroadcastSent(ctx context.Context, broadcastID, contactID int64) (bool, error) {
	key := fmt.Sprintf("broadcast:%d:sent", broadcastID)
	return r.Client.SIsMember(ctx, key, contactID).Result()
}

// ExpireBroadcastSet sets an expiry on the broadcast dedup set.
func (r *RedisClient) ExpireBroadcastSet(ctx context.Context, broadcastID int64, expiry time.Duration) error {
	key := fmt.Sprintf("broadcast:%d:sent", broadcastID)
	return r.Client.Expire(ctx, key, expiry).Err()
}

// LogEvent logs an event for analytics (e.g., messages sent per hour).
func (r *RedisClient) LogEvent(ctx context.Context, eventType string, userID int64) {
	key := fmt.Sprintf("events:%s:user:%d:%s", eventType, userID, time.Now().Format("2006-01-02"))
	if err := r.Client.Incr(ctx, key).Err(); err != nil {
		log.Printf("[Redis] Failed to log event %s: %v", eventType, err)
	}
	r.Client.Expire(ctx, key, 48*time.Hour)
}
