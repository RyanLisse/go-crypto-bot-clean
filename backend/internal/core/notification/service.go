package notification

import (
	"context"
	"fmt"
	"strings"

	"go-crypto-bot-clean/backend/internal/domain/notification/ports"
	log "github.com/sirupsen/logrus" // Using logrus for better logging
)

// service implements the NotificationService interface.
type service struct {
	notifiers      []ports.Notifier
	preferenceRepo ports.NotificationPreferenceRepository // Added preference repo
}

// NewService creates a new notification service instance.
func NewService(preferenceRepo ports.NotificationPreferenceRepository, notifiers []ports.Notifier) ports.NotificationService {
	if preferenceRepo == nil {
		// A preference repo is essential for this service logic
		log.Fatal("Notification Service: Preference repository cannot be nil")
		// Or return an error, depending on desired startup behavior
		// return nil, errors.New("preference repository cannot be nil")
	}
	return &service{
		notifiers:      notifiers,
		preferenceRepo: preferenceRepo, // Store the repo
	}
}

// SendNotification retrieves user preferences and sends the message via all enabled channels.
// It now uses userID to fetch preferences instead of directly taking channel/recipient.
func (s *service) SendNotification(ctx context.Context, userID string, subject string, message string) error {
	prefs, err := s.preferenceRepo.GetUserPreferences(ctx, userID)
	if err != nil {
		log.WithFields(log.Fields{"userID": userID, "error": err}).Error("Failed to retrieve user notification preferences")
		return fmt.Errorf("failed to get preferences for user %s: %w", userID, err)
	}

	if len(prefs) == 0 {
		log.WithField("userID", userID).Warn("No enabled notification preferences found for user.")
		// Depending on requirements, this might be an error or just informational.
		// Returning nil for now, meaning no error if no preferences are set/enabled.
		return nil
	}

	var sendErrors []string
	successfulSends := 0

	for _, pref := range prefs {
		if !pref.Enabled {
			continue // Skip disabled preferences (though GetUserPreferences should ideally only return enabled)
		}

		foundNotifier := false
		for _, notifier := range s.notifiers {
			if notifier.Supports(pref.Channel) {
				foundNotifier = true
				log.WithFields(log.Fields{
					"userID":    userID,
					"channel":   pref.Channel,
					"recipient": pref.Recipient, // Log recipient for debugging
				}).Info("Attempting to send notification")

				err := notifier.Send(ctx, pref.Recipient, subject, message)
				if err != nil {
					errMsg := fmt.Sprintf("failed to send notification via %s for user %s: %v", pref.Channel, userID, err)
					log.Error(errMsg)
					sendErrors = append(sendErrors, errMsg)
				} else {
					successfulSends++
					log.WithFields(log.Fields{
						"userID":  userID,
						"channel": pref.Channel,
					}).Info("Successfully sent notification")
				}
				break // Move to the next preference once a notifier is found and attempted
			}
		}

		if !foundNotifier {
			errMsg := fmt.Sprintf("no configured notifier supports channel '%s' specified in preference for user %s", pref.Channel, userID)
			log.Warn(errMsg)
			// Decide if this constitutes an error. Adding to sendErrors for now.
			sendErrors = append(sendErrors, errMsg)
		}
	}

	if len(sendErrors) > 0 {
		// If all attempts failed, return a more severe error.
		if successfulSends == 0 {
			log.WithField("userID", userID).Error("All notification attempts failed.")
			return fmt.Errorf("all notification attempts failed for user %s: %s", userID, strings.Join(sendErrors, "; "))
		}
		// If some succeeded, return a less severe error indicating partial failure.
		log.WithField("userID", userID).Warn("Some notification attempts failed.")
		return fmt.Errorf("some notification attempts failed for user %s: %s", userID, strings.Join(sendErrors, "; "))
	}

	log.WithField("userID", userID).Info("All notifications sent successfully.")
	return nil // All preferences processed successfully
}
