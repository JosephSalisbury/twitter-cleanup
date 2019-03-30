package notifier

import (
	"log"
)

type Notifier interface {
	Notify(string) error
}

type Config struct {
	Logger *log.Logger

	TwilioAccountSid string
	TwilioAuthToken  string
	TwilioNumberTo   string
	TwilioNumberFrom string
}
