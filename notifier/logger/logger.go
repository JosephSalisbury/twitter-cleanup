package logger

import (
	"errors"
	"log"

	"github.com/JosephSalisbury/twitter-cleanup/notifier"
)

type Logger struct {
	logger *log.Logger
}

func New(config notifier.Config) (*Logger, error) {
	if config.Logger == nil {
		return nil, errors.New("Logger logger cannot be empty.")
	}

	l := &Logger{
		logger: config.Logger,
	}

	return l, nil
}

func (l *Logger) Notify(message string) error {
	l.logger.Printf("%s\n", message)

	return nil
}
