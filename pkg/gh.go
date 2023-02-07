package pkg

import (
	"context"
	"net/http"
	"os"

	"github.com/google/go-github/github"
	"github.com/projectdiscovery/gologger"
	"golang.org/x/oauth2"
)

func GithubClient() *github.Client {
	var httpclient *http.Client
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		httpclient = oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))
	}
	githubClient := github.NewClient(httpclient)
	rateLimit, _, err := githubClient.RateLimits(context.Background())
	if err != nil {
		gologger.Error().Msgf("error while getting rate limits %s ", err.Error())

	}
	if rateLimit.Core.Remaining <= 0 {
		gologger.Error().Msgf("error for remaining request per hour: %s", err.Error())
		if arlErr, ok := err.(*github.AbuseRateLimitError); ok {
			// Provide user with more info regarding the rate limit
			gologger.Error().Msgf("error for remaining request per hour: %s, RetryAfter: %s", err.Error(), arlErr.RetryAfter)
		}
	}
	return githubClient
}
