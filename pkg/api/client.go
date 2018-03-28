package api

import (
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"net/http"
	"context"
)

// Create a type alias for local extension methods.
type (
	Client github.Client
	castClient = *Client
)

type RepositoryFetcher func(
	opt *github.RepositoryListByOrgOptions,
) ([]*github.Repository, *github.Response, error)

// Create a new Github client
func NewClient(accessToken string) *Client {
	// TODO set friendlier defaults for eg, read timeout
	var c = http.DefaultClient

	if accessToken != "" {
		c = createOAuthClient(accessToken)
	}

	// TODO support persistent http caching via "httpcache"

	return castClient(github.NewClient(c))
}

// TODO use the oauth 3-legged flow to get an access token
func createOAuthClient(token string) *http.Client {
	ctx := context.Background()
	src := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: token,
	})

	return oauth2.NewClient(ctx, src)
}

// Fetch all repos for an org, with pagination support.
func (c *Client) FetchOrgRepositories(
	ctx context.Context,
	org string,
	opt *github.RepositoryListByOrgOptions) ([]*github.Repository, error) {
	if opt == nil {
		opt = &github.RepositoryListByOrgOptions{}
	}

	// use the max page size of 100 repos
	opt.ListOptions.PerPage = 100

	// get all pages of results
	var allRepos []*github.Repository
	for {
		repos, resp, err := c.Repositories.ListByOrg(ctx, org, opt)
		if err != nil {
			return nil, err
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return allRepos, nil
}

// Count all pull requests for a given repo, with pagination support.
func (c *Client) CountRepositoryPRs(
	ctx context.Context,
	repo *github.Repository,
	opt *github.PullRequestListOptions) (uint, error) {
	if opt == nil {
		opt = &github.PullRequestListOptions{}
	}

	opt.ListOptions.PerPage = 100

	var count uint = 0
	for {
		pulls, resp, err := c.PullRequests.List(
			ctx,
			repo.GetOwner().GetName(),
			repo.GetName(),
			opt,
		)
		if err != nil {
			return 0, err
		}

		count += uint(len(pulls))
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return count, nil
}