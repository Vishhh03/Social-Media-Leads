package store

import (
	"context"
	"time"

	"github.com/social-media-lead/backend/internal/models"
)

// Store defines the interface for our database layer.
// This allows us to use either a real PostgreSQL connection or a Mock store for testing.
type Store interface {
	// Users
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id int64) (*models.User, error)
	GetOrCreateOAuthUser(ctx context.Context, oauthUser *models.User) (*models.User, error)
	UpdateUserProfile(ctx context.Context, userID int64, fullName, email, companyName string) (*models.User, error)
	UpdateUserPassword(ctx context.Context, userID int64, passwordHash string) error
	UpdateChannelToken(ctx context.Context, channelID int64, accessToken string, expiry time.Time) error

	// Storage Management
	Close()
	RunMigrations() error

	// Messages
	CreateMessage(ctx context.Context, m *models.Message) error
	GetMessagesByContact(ctx context.Context, contactID int64, limit, offset int) ([]models.Message, error)
	GetConversations(ctx context.Context, userID int64, limit, offset int) ([]models.Message, error)

	// Contacts
	CreateContact(ctx context.Context, c *models.Contact) error
	GetOrCreateContact(ctx context.Context, c *models.Contact) error
	GetContactsByUser(ctx context.Context, userID int64, limit, offset int) ([]models.Contact, error)
	UpdateContactLead(ctx context.Context, contactID int64, budget, location, timeline, phone string, isHot bool) error
	GetContactByID(ctx context.Context, contactID int64) (*models.Contact, error)

	// Channels
	CreateChannel(ctx context.Context, ch *models.Channel) error
	GetChannelsByUser(ctx context.Context, userID int64) ([]models.Channel, error)
	GetChannelByAccountID(ctx context.Context, platform, accountID string) (*models.Channel, error)
	GetChannelByID(ctx context.Context, channelID int64) (*models.Channel, error)
	DeleteChannel(ctx context.Context, channelID, userID int64) error

	// Broadcasts
	CreateBroadcast(ctx context.Context, b *models.Broadcast) error
	GetBroadcastsByUser(ctx context.Context, userID int64, limit, offset int) ([]models.Broadcast, error)
	GetBroadcastByID(ctx context.Context, broadcastID int64) (*models.Broadcast, error)
	UpdateBroadcastStatus(ctx context.Context, broadcastID int64, status string, totalSent, totalFailed int) error

	// Automations
	CreateAutomation(ctx context.Context, a *models.Automation) error
	GetAutomationsByUser(ctx context.Context, userID int64) ([]models.Automation, error)
	UpdateAutomation(ctx context.Context, a *models.Automation) error
	DeleteAutomation(ctx context.Context, automationID, userID int64) error

	// Knowledge Base (RAG)
	CreateKnowledgeBaseEntry(ctx context.Context, entry *models.KnowledgeBase, embedding []float32) error
	GetKnowledgeBaseEntriesByUser(ctx context.Context, userID int64) ([]models.KnowledgeBase, error)
	SearchKnowledgeBase(ctx context.Context, userID int64, queryEmbedding []float32, limit int) ([]models.KnowledgeBase, error)
	DeleteKnowledgeBaseEntry(ctx context.Context, entryID, userID int64) error

	// Workflows
	CreateWorkflow(ctx context.Context, w *models.Workflow) error
	GetWorkflowByID(ctx context.Context, workflowID int64) (*models.Workflow, error)
	GetWorkflowsByUser(ctx context.Context, userID int64) ([]models.Workflow, error)
	GetActiveWorkflowsByTrigger(ctx context.Context, userID int64, triggerType string) ([]models.Workflow, error)
	UpdateWorkflow(ctx context.Context, w *models.Workflow) error
	DeleteWorkflow(ctx context.Context, workflowID, userID int64) error

	// Workflow Executions
	CreateWorkflowExecution(ctx context.Context, exec *models.WorkflowExecution) error
	GetWorkflowExecutionByID(ctx context.Context, executionID int64) (*models.WorkflowExecution, error)
	UpdateWorkflowExecution(ctx context.Context, exec *models.WorkflowExecution) error
}

// Ensure Storage implements Store at compile time.
var _ Store = (*Storage)(nil)
