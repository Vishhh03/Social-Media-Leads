package models

import "time"

// User represents a SaaS customer (builder, agency, marketer).
type User struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	FullName     string    `json:"full_name"`
	CompanyName  string    `json:"company_name,omitempty"`
	GoogleID     string    `json:"google_id,omitempty"`
	AvatarURL    string    `json:"avatar_url,omitempty"`
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

// ============================================
// AI & Orchestrator Models
// ============================================

// KnowledgeBase represents a document chunk used for RAG
type KnowledgeBase struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	// Note: We don't expose the 'embedding' float32 array in standard JSON responses
	// to save bandwidth, unless specifically requested.
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Workflow represents the Blueprint (DAG) of an automation
type Workflow struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Name        string    `json:"name"`
	TriggerType string    `json:"trigger_type"` // e.g., 'meta_dm_received'
	Status      string    `json:"status"`       // "draft", "published"
	Prompt      string    `json:"prompt,omitempty"` // Original AI prompt if generated
	// Nodes and Edges are stored as JSONB in DB, we use generic map/interfaces or raw JSON here
	Nodes       []byte    `json:"nodes"` 
	Edges       []byte    `json:"edges"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// WorkflowExecution represents the runtime state of a specific Lead passing through a Workflow
type WorkflowExecution struct {
	ID            int64     `json:"id"`
	WorkflowID    int64     `json:"workflow_id"`
	ContactID     int64     `json:"contact_id"`
	CurrentNodeID string    `json:"current_node_id"`
	Status        string    `json:"status"` // "running", "waiting", "completed", "failed"
	StateData     []byte    `json:"state_data"` // Context payload (JSONB)
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
