package notification_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	notification_service "go-crypto-bot-clean/backend/internal/core/notification"
	"go-crypto-bot-clean/backend/internal/domain/notification"
	"go-crypto-bot-clean/backend/internal/domain/notification/ports"
)

// MockNotifier is a mock implementation of the Notifier interface
type MockNotifier struct {
	mock.Mock
	SupportedChannel string
	SendError        error
	SendCalls        int
	LastRecipient    string
}

func (m *MockNotifier) Send(ctx context.Context, recipient string, subject string, message string) error {
	m.SendCalls++
	m.LastRecipient = recipient
	return m.SendError
}

func (m *MockNotifier) Supports(channel string) bool {
	return channel == m.SupportedChannel
}

// MockPreferenceRepository for testing
type MockPreferenceRepository struct {
	prefs map[string][]notification.Preference
}

func (m *MockPreferenceRepository) GetUserPreferences(ctx context.Context, userID string) ([]notification.Preference, error) {
	if prefs, ok := m.prefs[userID]; ok {
		enabledPrefs := make([]notification.Preference, 0, len(prefs))
		for _, p := range prefs {
			if p.Enabled {
				enabledPrefs = append(enabledPrefs, p)
			}
		}
		return enabledPrefs, nil
	}
	return []notification.Preference{}, nil
}

// TestSendNotification tests sending a basic notification
func TestSendNotification(t *testing.T) {
	ctx := context.Background()
	recipient := "user@example.com"
	subject := "Test Subject"
	message := "Test Message"
	channel := "email" // Assume email notifier exists

	mockEmailNotifier := new(MockNotifier)
	// Setup expectations
	mockEmailNotifier.On("Supports", channel).Return(true)
	mockEmailNotifier.On("Send", ctx, recipient, subject, message).Return(nil)

	notifiers := []ports.Notifier{mockEmailNotifier}
	mockRepo := &MockPreferenceRepository{ // Create mock repo
		prefs: map[string][]notification.Preference{
			"user@example.com": {{UserID: "user@example.com", Channel: channel, Recipient: recipient, Enabled: true}},
		},
	}
	service := notification_service.NewService(mockRepo, notifiers)

	err := service.SendNotification(ctx, recipient, subject, message) // Removed channel argument

	assert.NoError(t, err)
	mockEmailNotifier.AssertExpectations(t)
}

// TestSendNotificationUnsupportedChannel tests sending when the *specific* channel isn't supported by any available notifier.
func TestSendNotificationUnsupportedChannel(t *testing.T) {
	ctx := context.Background()
	recipient := "user@example.com"
	subject := "Test Subject"
	message := "Test Message"
	channel := "unsupported_channel"

	mockEmailNotifier := new(MockNotifier)
	// Setup expectations - Supports returns false for the target channel
	mockEmailNotifier.On("Supports", channel).Return(false)
	// Send should NOT be called

	notifiers := []ports.Notifier{mockEmailNotifier}
	mockRepo := &MockPreferenceRepository{ // Create mock repo
		prefs: map[string][]notification.Preference{
			"user@example.com": {{UserID: "user@example.com", Channel: channel, Recipient: recipient, Enabled: true}},
		},
	}
	service := notification_service.NewService(mockRepo, notifiers)

	err := service.SendNotification(ctx, recipient, subject, message) // Removed channel argument

	assert.Error(t, err) // Expect an error
	// Adjust assertion to match the actual error message when no supporting notifier is found for the user's PREFERRED channel
	assert.Contains(t, err.Error(), "no enabled preferences or supporting notifiers found for user")
	mockEmailNotifier.AssertNotCalled(t, "Send", ctx, recipient, subject, message) // Ensure Send wasn't called
}

// TestSendNotificationNoSupportingNotifier tests sending when no notifier supports the channel
func TestSendNotificationNoSupportingNotifier(t *testing.T) {
	ctx := context.Background()
	recipient := "user@example.com"
	subject := "Test Subject"
	message := "Test Message"
	channel := "email" // User preference

	mockOtherNotifier := new(MockNotifier)
	// Setup expectations - Supports returns false for the target channel
	mockOtherNotifier.On("Supports", channel).Return(false)

	notifiers := []ports.Notifier{mockOtherNotifier}
	mockRepo := &MockPreferenceRepository{ // Create mock repo
		prefs: map[string][]notification.Preference{
			"user@example.com": {{UserID: "user@example.com", Channel: channel, Recipient: recipient, Enabled: true}},
		},
	}
	service := notification_service.NewService(mockRepo, notifiers)

	err := service.SendNotification(ctx, recipient, subject, message) // Removed channel argument

	assert.Error(t, err)                                                                    // Expect an error
	assert.Contains(t, err.Error(), "no enabled preferences or supporting notifiers found") // Check specific error message
	mockOtherNotifier.AssertNotCalled(t, "Send", ctx, recipient, subject, message)
}

func TestNotificationService_SendNotification(t *testing.T) {
	// Use local mocks defined in this file
	mockNotifier := &MockNotifier{}
	mockRepo := &MockPreferenceRepository{
		prefs: make(map[string][]notification.Preference), // Initialize map
	}

	t.Run("successful notification", func(t *testing.T) {
		// Reset mock calls and expectations for this subtest
		mockNotifier = &MockNotifier{}
		mockRepo = &MockPreferenceRepository{
			prefs: map[string][]notification.Preference{
				"test_user": {{UserID: "test_user", Channel: "test_channel", Recipient: "test_user", Enabled: true}},
			},
		}

		// Setup expectations on local mock
		mockNotifier.On("Supports", "test_channel").Return(true)
		mockNotifier.On("Send", mock.Anything, "test_user", "test_subject", "test_message").Return(nil)

		service := notification_service.NewService(mockRepo, []ports.Notifier{mockNotifier})
		err := service.SendNotification(context.Background(), "test_user", "test_subject", "test_message")
		require.NoError(t, err)
		mockNotifier.AssertExpectations(t)
	})

	t.Run("notification disabled", func(t *testing.T) {
		// Reset mock calls and expectations
		mockNotifier = &MockNotifier{}
		mockRepo = &MockPreferenceRepository{
			prefs: map[string][]notification.Preference{
				"test_user": {{UserID: "test_user", Channel: "test_channel", Recipient: "test_user", Enabled: false}},
			},
		}

		service := notification_service.NewService(mockRepo, []ports.Notifier{mockNotifier})
		err := service.SendNotification(context.Background(), "test_user", "test_subject", "test_message")
		require.NoError(t, err)
		mockNotifier.AssertNotCalled(t, "Send", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("notification error", func(t *testing.T) {
		// Reset mock calls and expectations
		mockNotifier = &MockNotifier{}
		mockRepo = &MockPreferenceRepository{
			prefs: map[string][]notification.Preference{
				"test_user": {{UserID: "test_user", Channel: "test_channel", Recipient: "test_user", Enabled: true}},
			},
		}

		// Setup expectations on local mock
		mockNotifier.On("Supports", "test_channel").Return(true)
		mockNotifier.On("Send", mock.Anything, "test_user", "test_subject", "test_message").Return(errors.New("test error"))

		service := notification_service.NewService(mockRepo, []ports.Notifier{mockNotifier})
		err := service.SendNotification(context.Background(), "test_user", "test_subject", "test_message")
		require.Error(t, err)
		mockNotifier.AssertExpectations(t)
	})
}
