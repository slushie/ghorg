package cmd

import (
	"github.com/spf13/cobra"
	"github.com/slushie/ghorg/pkg/api"
	"context"
	"github.com/slushie/ghorg/pkg/repos"
	"github.com/google/go-github/github"
)

// forksCmd represents the forks command
var forksCmd = &cobra.Command{
	Use:   "forks [organization]",
	Short: "List repos by forks",
	Long:  `List all repositories in the organization, sorted by number of forks.`,
	Run:   runForks,
}

func init() {
	rootCmd.AddCommand(forksCmd)
	addListFlags(forksCmd)
}

func runForks(cmd *cobra.Command, args []string) {
	gh := api.NewClient(accessToken)
	ctx := context.Background()

	rs, err := gh.FetchOrgRepositories(ctx, organization, nil)
	if err != nil {
		api.Fail(err)
	}

	// these are the column names
	fields = []string{"Forks", "Name", "URL"}

	// set up the output list
	repoList = repos.NewList()
	repoList.Add(rs...)
	repoList.Marshal = MarshalRepo
	repoList.Compare = CompareRepoForks
}

func CompareRepoForks(a, b *github.Repository) bool {
	return a.GetStargazersCount() > b.GetStargazersCount()
}
