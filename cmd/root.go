package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/JosephSalisbury/twitter-cleanup/notifier"
	"github.com/JosephSalisbury/twitter-cleanup/notifier/unified"
	"github.com/JosephSalisbury/twitter-cleanup/twitter"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:              "twitter-cleanup",
		Short:            "twitter-cleanup is a tool for cleaning up Twitter",
		PersistentPreRun: handleDependencies,
	}

	notifierType string

	logger         *log.Logger
	notifierClient notifier.Notifier
	twitterClient  *twitter.Twitter

	twitterAccessToken       string
	twitterAccessTokenSecret string
	twitterConsumerKey       string
	twitterConsumerSecret    string

	twilioAccountSid string
	twilioAuthToken  string
	twilioNumberTo   string
	twilioNumberFrom string
)

func init() {
	rootCmd.PersistentFlags().StringVar(
		&twitterAccessToken, "twitter-access-token", "",
		"Twitter Access Token.",
	)
	rootCmd.PersistentFlags().StringVar(
		&twitterAccessTokenSecret, "twitter-access-token-secret", "",
		"Twitter Access Token Secret.",
	)
	rootCmd.PersistentFlags().StringVar(
		&twitterConsumerKey, "twitter-consumer-key", "",
		"Twitter Consumer Key.",
	)
	rootCmd.PersistentFlags().StringVar(
		&twitterConsumerSecret, "twitter-consumer-secret", "",
		"Twitter Consumer Secret.",
	)

	rootCmd.PersistentFlags().StringVar(
		&notifierType, "notifier", "logger",
		"What notifier to use. Either 'logger' or 'twilio'.",
	)

	rootCmd.PersistentFlags().StringVar(
		&twilioAccountSid, "twilio-account-sid", "",
		"Twilio Account SID.",
	)
	rootCmd.PersistentFlags().StringVar(
		&twilioAuthToken, "twilio-auth-token", "",
		"Twilio Auth Token.",
	)
	rootCmd.PersistentFlags().StringVar(
		&twilioNumberTo, "twilio-number-to", "",
		"Twilio Number To.",
	)
	rootCmd.PersistentFlags().StringVar(
		&twilioNumberFrom, "twilio-number-from", "",
		"Twilio Number From.",
	)
}

func handleDependencies(cmd *cobra.Command, args []string) {
	var err error

	logger = log.New(os.Stdout, "", log.LstdFlags)

	notifierClient, err = unified.GetNotifier(notifierType, notifier.Config{
		Logger: logger,

		TwilioAccountSid: twilioAccountSid,
		TwilioAuthToken:  twilioAuthToken,
		TwilioNumberTo:   twilioNumberTo,
		TwilioNumberFrom: twilioNumberFrom,
	})
	if err != nil {
		panic(err)
	}

	twitterClient, err = twitter.New(twitter.Config{
		Logger:   logger,
		Notifier: notifierClient,

		AccessToken:       twitterAccessToken,
		AccessTokenSecret: twitterAccessTokenSecret,
		ConsumerKey:       twitterConsumerKey,
		ConsumerSecret:    twitterConsumerSecret,
	})
	if err != nil {
		panic(err)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
}
