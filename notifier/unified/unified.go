package unified

import (
	"errors"
	"log"

	"github.com/JosephSalisbury/twitter-cleanup/notifier"
	"github.com/JosephSalisbury/twitter-cleanup/notifier/logger"
	"github.com/JosephSalisbury/twitter-cleanup/notifier/twilio"
)

type Config struct {
	Logger *log.Logger

	TwilioAccountSid string
	TwilioAuthToken  string
	TwilioNumberTo   string
	TwilioNumberFrom string
}

func GetNotifier(notifierType string, config notifier.Config) (notifier.Notifier, error) {
	switch notifierType {
	case "logger":
		return logger.New(config)
	case "twilio":
		return twilio.New(config)
	default:
		return nil, errors.New("Unknown notifier type")
	}
}
