package pkg

import (
	"context"
	"net/http"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func GithubClient() *github.Client {
	var httpclient *http.Client
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {

		httpclient = oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))
	}
	githubClient := github.NewClient(httpclient)
	return githubClient
}
