package twitter

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/JosephSalisbury/twitter-cleanup/notifier"
)

type Config struct {
	Logger   *log.Logger
	Notifier notifier.Notifier

	AccessToken       string
	AccessTokenSecret string
	ConsumerKey       string
	ConsumerSecret    string
}

type Twitter struct {
	api      *anaconda.TwitterApi
	logger   *log.Logger
	notifier notifier.Notifier
}

func New(config Config) (*Twitter, error) {
	if config.Logger == nil {
		return nil, errors.New("Twitter Logger cannot be empty.")
	}
	if config.Notifier == nil {
		return nil, errors.New("Twitter Notifier cannot be empty.")
	}

	if config.AccessToken == "" {
		return nil, errors.New("Twitter Access Token cannot be empty.")
	}
	if config.AccessTokenSecret == "" {
		return nil, errors.New("Twitter Access Token Secret cannot be empty.")
	}
	if config.ConsumerKey == "" {
		return nil, errors.New("Twitter Consumer Key cannot be empty.")
	}
	if config.ConsumerSecret == "" {
		return nil, errors.New("Twitter Consumer Secret cannot be empty.")
	}

	api := anaconda.NewTwitterApiWithCredentials(
		config.AccessToken,
		config.AccessTokenSecret,
		config.ConsumerKey,
		config.ConsumerSecret,
	)

	t := &Twitter{
		api:      api,
		logger:   config.Logger,
		notifier: config.Notifier,
	}

	return t, nil
}

func (t *Twitter) UnfollowStaleAccounts(timespan time.Duration, maxUnfollow int) error {
	t.logger.Printf("Finding stale accounts\n")

	unfollow := []anaconda.User{}
	for page := range t.api.GetFriendsListAll(nil) {
		t.logger.Printf("Getting next page of friends\n")

		for _, user := range page.Friends {
			t.logger.Printf("Checking staleness of '%v' (%v)\n", user.Name, user.ScreenName)

			v := url.Values{}
			v.Set("user_id", user.IdStr)
			v.Set("count", "1")
			timeline, err := t.api.GetUserTimeline(v)
			if err != nil {
				return err
			}

			// If the user has never tweeted, move on.
			if len(timeline) == 0 {
				continue
			}

			latestTweet := timeline[0]
			createdTime, err := latestTweet.CreatedAtTime()
			if err != nil {
				return err
			}

			// If the time of the last tweet, plus the timespan,
			// is less than now, add to the unfollow list.
			if createdTime.Add(timespan).Before(time.Now()) {
				t.logger.Printf("User is stale\n")
				unfollow = append(unfollow, user)
			}

			// If we have enough stale accounts, stop.
			if len(unfollow) >= maxUnfollow {
				t.logger.Printf("We have %v stale accounts, moving on", len(unfollow))
				goto Unfollow
			}
		}
	}

Unfollow:
	t.logger.Printf("Unfollowing stale accounts\n")

	for _, user := range unfollow {
		t.logger.Printf("Unfollowing '%v' (%v)\n", user.Name, user.ScreenName)

		if _, err := t.api.UnfollowUserId(user.Id); err != nil {
			return err
		}

		if err := t.notifier.Notify(fmt.Sprintf("You have unfollowed '%v' (%v) due to inactivity.", user.Name, user.ScreenName)); err != nil {
			return err
		}
	}

	t.logger.Printf("Unfollowed stale accounts\n")

	return nil
}
