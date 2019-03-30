package cmd

import (
	"time"

	"github.com/spf13/cobra"
)

var (
	followingCmd = &cobra.Command{
		Use:   "following",
		Short: "Unfollow stale followings",
		Long:  "Unfollows followings that have not tweeted in the configured timespan.",
		Run:   followers,
	}

	timespan    time.Duration
	maxUnfollow int
)

func init() {
	rootCmd.AddCommand(followingCmd)

	followingCmd.Flags().DurationVar(
		&timespan, "timespan", time.Hour*24*30*6,
		"How long since last tweet to determine if a following is stale",
	)
	followingCmd.Flags().IntVar(
		&maxUnfollow, "max-unfollow", 1,
		"The maximum number of accounts to unfollow at once",
	)
}

func followers(cmd *cobra.Command, args []string) {
	if err := twitterClient.UnfollowStaleAccounts(timespan, maxUnfollow); err != nil {
		panic(err)
	}
}
