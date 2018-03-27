package org

import (
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"net/http"
	"context"
)

func NewClient(accessToken string) *github.Client {
	var c = http.DefaultClient
	if accessToken != "" {
		c = createOAuthClient(accessToken)
	}

	// TODO support persistent http caching via "httpcache"
	return github.NewClient(c)
}

// TODO use the oauth 3-legged flow to get an access token
func createOAuthClient(token string) *http.Client {
	ctx := context.Background()
	src := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: token,
	})

	return oauth2.NewClient(ctx, src)
}