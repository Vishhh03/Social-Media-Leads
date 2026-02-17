package models

import "time"

// User represents a SaaS customer (builder, agency, marketer).
type User struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	FullName     string    `json:"full_name"`
	CompanyName  string    `json:"company_name,omitempty"`
	Plan         string    `json:"plan"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Channel represents a connected social media account (WhatsApp, Instagram, Facebook).
type Channel struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	Platform     string    `json:"platform"` // "whatsapp", "instagram", "facebook"
	AccountID    string    `json:"account_id"`
	AccountName  string    `json:"account_name"`
	AccessToken  string    `json:"-"`
	RefreshToken string    `json:"-"`
	TokenExpiry  time.Time `json:"token_expiry,omitempty"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Contact represents a lead/customer who messaged via any channel.
type Contact struct {
	ID               int64     `json:"id"`
	UserID           int64     `json:"user_id"`
	ChannelID        int64     `json:"channel_id"`
	Platform         string    `json:"platform"`
	PlatformUserID   string    `json:"platform_user_id"`
	Name             string    `json:"name"`
	Phone            string    `json:"phone,omitempty"`
	Email            string    `json:"email,omitempty"`
	Budget           string    `json:"budget,omitempty"`
	PreferredLocation string  `json:"preferred_location,omitempty"`
	PurchaseTimeline string    `json:"purchase_timeline,omitempty"`
	Tags             []string  `json:"tags,omitempty"`
	IsHotLead        bool      `json:"is_hot_lead"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// Message represents a single chat message.
type Message struct {
	ID             int64     `json:"id"`
	UserID         int64     `json:"user_id"`
	ChannelID      int64     `json:"channel_id"`
	ContactID      int64     `json:"contact_id"`
	Platform       string    `json:"platform"`
	Direction      string    `json:"direction"` // "inbound" or "outbound"
	Content        string    `json:"content"`
	MessageType    string    `json:"message_type"` // "text", "image", "document", "template"
	PlatformMsgID  string    `json:"platform_msg_id,omitempty"`
	Status         string    `json:"status"` // "sent", "delivered", "read", "failed"
	IsAutomated    bool      `json:"is_automated"`
	CreatedAt      time.Time `json:"created_at"`
}

// Automation represents a keyword-trigger automation rule.
type Automation struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Name        string    `json:"name"`
	TriggerType string    `json:"trigger_type"` // "keyword", "first_message"
	Keywords    []string  `json:"keywords"`
	ReplyText   string    `json:"reply_text"`
	ReplyMedia  string    `json:"reply_media,omitempty"` // URL to file
	DelayMs     int       `json:"delay_ms"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Broadcast represents a bulk message campaign.
type Broadcast struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Name        string    `json:"name"`
	Content     string    `json:"content"`
	MediaURL    string    `json:"media_url,omitempty"`
	Status      string    `json:"status"` // "draft", "scheduled", "sending", "sent"
	TotalSent   int       `json:"total_sent"`
	TotalFailed int       `json:"total_failed"`
	ScheduledAt *time.Time `json:"scheduled_at,omitempty"`
	SentAt      *time.Time `json:"sent_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
