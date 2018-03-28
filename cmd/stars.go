package cmd

import (
	"github.com/spf13/cobra"
	"github.com/slushie/ghorg/pkg/api"
	"context"
	"github.com/slushie/ghorg/pkg/repos"
	"github.com/google/go-github/github"
)

// starsCmd represents the stars command
var starsCmd = &cobra.Command{
	Use:   "stars [organization]",
	Short: "List repos by stargazers",
	Long:  `List all repositories in the organization, sorted by number of stars.`,
	Run:   runStars,
}

func init() {
	rootCmd.AddCommand(starsCmd)
	addListFlags(starsCmd)
}

func runStars(cmd *cobra.Command, args []string) {
	gh := api.NewClient(accessToken)
	ctx := context.Background()

	rs, err := gh.FetchOrgRepositories(ctx, organization, nil)
	if err != nil {
		api.Fail(err)
	}

	// these are the column names
	fields = []string{"Stars", "Name", "URL"}

	// set up the output list
	repoList = repos.NewList()
	repoList.Add(rs...)
	repoList.Marshal = MarshalRepo
	repoList.Compare = CompareRepoStars
}

func CompareRepoStars(a, b *github.Repository) bool {
	return a.GetStargazersCount() > b.GetStargazersCount()
}
