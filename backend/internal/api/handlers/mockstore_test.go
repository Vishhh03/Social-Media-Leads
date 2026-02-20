package handlers_test

import (
	"context"
	"errors"
	"time"

	"github.com/social-media-lead/backend/internal/models"
)

// MockStore is a mock implementation of the store.Store interface for testing.
type MockStore struct {
	Users          map[int64]*models.User
	UsersByEmail   map[string]*models.User
	CreateUserFunc func(ctx context.Context, user *models.User) error
}

func NewMockStore() *MockStore {
	return &MockStore{
		Users:        make(map[int64]*models.User),
		UsersByEmail: make(map[string]*models.User),
	}
}

// Implement required interfaces for testing Auth
func (m *MockStore) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	if user, exists := m.UsersByEmail[email]; exists {
		return user, nil
	}
	return nil, errors.New("user not found")
}

func (m *MockStore) CreateUser(ctx context.Context, user *models.User) error {
	if m.CreateUserFunc != nil {
		return m.CreateUserFunc(ctx, user)
	}
	user.ID = int64(len(m.Users) + 1)
	m.Users[user.ID] = user
	m.UsersByEmail[user.Email] = user
	return nil
}

func (m *MockStore) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	if user, exists := m.Users[id]; exists {
		return user, nil
	}
	return nil, errors.New("user not found")
}

// Dummy Implementations for the rest of the interface to satisfy the compiler
func (m *MockStore) GetOrCreateOAuthUser(ctx context.Context, oauthUser *models.User) (*models.User, error) { return nil, nil }
func (m *MockStore) UpdateUserProfile(ctx context.Context, userID int64, fullName, email, companyName string) (*models.User, error) { return nil, nil }
func (m *MockStore) UpdateUserPassword(ctx context.Context, userID int64, passwordHash string) error { return nil }
func (m *MockStore) UpdateChannelToken(ctx context.Context, channelID int64, accessToken string, expiry time.Time) error { return nil }
func (m *MockStore) Close() {}
func (m *MockStore) RunMigrations() error { return nil }
func (m *MockStore) CreateMessage(ctx context.Context, msg *models.Message) error { return nil }
func (m *MockStore) GetMessagesByContact(ctx context.Context, contactID int64, limit, offset int) ([]models.Message, error) { return nil, nil }
func (m *MockStore) GetConversations(ctx context.Context, userID int64, limit, offset int) ([]models.Message, error) { return nil, nil }
func (m *MockStore) CreateContact(ctx context.Context, c *models.Contact) error { return nil }
func (m *MockStore) GetOrCreateContact(ctx context.Context, c *models.Contact) error { return nil }
func (m *MockStore) GetContactsByUser(ctx context.Context, userID int64, limit, offset int) ([]models.Contact, error) { return nil, nil }
func (m *MockStore) UpdateContactLead(ctx context.Context, contactID int64, budget, location, timeline, phone string, isHot bool) error { return nil }
func (m *MockStore) GetContactByID(ctx context.Context, contactID int64) (*models.Contact, error) { return nil, nil }
func (m *MockStore) CreateChannel(ctx context.Context, ch *models.Channel) error { return nil }
func (m *MockStore) GetChannelsByUser(ctx context.Context, userID int64) ([]models.Channel, error) { return nil, nil }
func (m *MockStore) GetChannelByAccountID(ctx context.Context, platform, accountID string) (*models.Channel, error) { return nil, nil }
func (m *MockStore) GetChannelByID(ctx context.Context, channelID int64) (*models.Channel, error) { return nil, nil }
func (m *MockStore) DeleteChannel(ctx context.Context, channelID, userID int64) error { return nil }
func (m *MockStore) CreateBroadcast(ctx context.Context, b *models.Broadcast) error { return nil }
func (m *MockStore) GetBroadcastsByUser(ctx context.Context, userID int64, limit, offset int) ([]models.Broadcast, error) { return nil, nil }
func (m *MockStore) GetBroadcastByID(ctx context.Context, broadcastID int64) (*models.Broadcast, error) { return nil, nil }
func (m *MockStore) UpdateBroadcastStatus(ctx context.Context, broadcastID int64, status string, totalSent, totalFailed int) error { return nil }
func (m *MockStore) CreateAutomation(ctx context.Context, a *models.Automation) error { return nil }
func (m *MockStore) GetAutomationsByUser(ctx context.Context, userID int64) ([]models.Automation, error) { return nil, nil }
func (m *MockStore) UpdateAutomation(ctx context.Context, a *models.Automation) error { return nil }
func (m *MockStore) DeleteAutomation(ctx context.Context, automationID, userID int64) error { return nil }
