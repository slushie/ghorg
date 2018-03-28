package cmd

import (
	"github.com/spf13/cobra"
	"github.com/slushie/ghorg/pkg/api"
	"context"
	"github.com/slushie/ghorg/pkg/repos"
	"github.com/google/go-github/github"
)

// pullsCmd represents the pulls command
var pullsCmd = &cobra.Command{
	Use:   "pulls [organization]",
	Aliases: []string{"pr", "prs"},
	Short: "List repos by PRs",
	Long:  `List all repositories in the organization, sorted by number of pull requests.`,
	Run:   runPulls,
}

var pullCounts RepoPullCounts

func init() {
	rootCmd.AddCommand(pullsCmd)
	addListFlags(pullsCmd)
	addPullFlags(pullsCmd)
}

func runPulls(cmd *cobra.Command, args []string) {
	gh := api.NewClient(accessToken)
	ctx := context.Background()

	rs, err := gh.FetchOrgRepositories(ctx, organization, nil)
	if err != nil {
		api.Fail(err)
	}

	// these are the column names
	fields = []string{"PRs", "Name", "URL"}

	// set up the output list
	repoList = repos.NewList()
	repoList.Add(rs...)
	repoList.Marshal = MarshalRepo
	repoList.Compare = CompareRepoPulls

	// wait for pr stats
	pullCounts, err = countRepoPulls(ctx, gh, rs)
	if err != nil {
		api.Fail(err)
	}
}

func CompareRepoPulls(a, b *github.Repository) bool {
	var (
		aID = a.GetID()
		bID = b.GetID()
	)

	return pullCounts[aID].Count > pullCounts[bID].Count
}

