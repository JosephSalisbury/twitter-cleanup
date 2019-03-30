package twilio

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/JosephSalisbury/twitter-cleanup/notifier"
)

const (
	apiUrlFormat = "https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json"
)

type Twilio struct {
	logger *log.Logger

	accountSid string
	authToken  string
	numberTo   string
	numberFrom string
}

func New(config notifier.Config) (*Twilio, error) {
	if config.Logger == nil {
		return nil, errors.New("Twilio logger cannot be empty.")
	}

	if config.TwilioAccountSid == "" {
		return nil, errors.New("Twilio Account Sid cannot be empty.")
	}
	if config.TwilioAuthToken == "" {
		return nil, errors.New("Twilio Auth Token cannot be empty.")
	}
	if config.TwilioNumberTo == "" {
		return nil, errors.New("Twilio Number To cannot be empty.")
	}
	if config.TwilioNumberFrom == "" {
		return nil, errors.New("Twilio Number From cannot be empty.")
	}

	t := &Twilio{
		logger: config.Logger,

		accountSid: config.TwilioAccountSid,
		authToken:  config.TwilioAuthToken,
		numberTo:   config.TwilioNumberTo,
		numberFrom: config.TwilioNumberFrom,
	}

	return t, nil
}

func (t *Twilio) Notify(message string) error {
	t.logger.Printf("Sending text message\n")

	apiUrl := fmt.Sprintf(apiUrlFormat, t.accountSid)

	data := url.Values{}
	data.Set("To", t.numberTo)
	data.Set("From", t.numberFrom)
	data.Set("Body", message)

	reader := strings.NewReader(data.Encode())

	req, err := http.NewRequest(http.MethodPost, apiUrl, reader)
	if err != nil {
		return err
	}

	req.SetBasicAuth(t.accountSid, t.authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		t.logger.Printf("Failed to send text message: %s\n", string(body))
		return errors.New("Could not send text message")
	}

	return nil
}
