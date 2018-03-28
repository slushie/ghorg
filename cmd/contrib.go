package cmd

import (
	"github.com/spf13/cobra"
	"github.com/slushie/ghorg/pkg/api"
	"context"
	"github.com/slushie/ghorg/pkg/repos"
	"github.com/google/go-github/github"
)

// contribCmd represents the contrib command
var contribCmd = &cobra.Command{
	Use:     "contrib [organization]",
	Short:   "List repos by contribution percentage",
	Long:    `List all repositories in the organization, sorted by contribution percentage. This
value is calculated by dividing the number of pull requests by the number of forks.`,
	Run:     runContrib,
	PreRunE: parsePullFlags,
}

func init() {
	rootCmd.AddCommand(contribCmd)
	addListFlags(contribCmd)
	addPullFlags(contribCmd)
}

func runContrib(cmd *cobra.Command, args []string) {
	gh := api.NewClient(accessToken)

	// TODO better cancellation and deadline handling via context.Context
	ctx := context.Background()

	rs, err := gh.FetchOrgRepositories(ctx, organization, nil)
	if err != nil {
		api.Fail(err)
	}

	// these are the column names
	fields = []string{"Contrib-%", "PRs", "Forks", "Name", "URL"}

	// set up the output list
	repoList = repos.NewList()
	repoList.Add(rs...)
	repoList.Marshal = MarshalRepo
	repoList.Compare = CompareRepoContrib

	// wait for pr stats
	pullCounts, err = countRepoPulls(ctx, gh, rs)
	if err != nil {
		api.Fail(err)
	}
}

func CompareRepoContrib(a, b *github.Repository) bool {
	aID := a.GetID()
	aPct := float64(pullCounts[aID]) / float64(a.GetForksCount())

	bID := b.GetID()
	bPct := float64(pullCounts[bID]) / float64(b.GetForksCount())

	return aPct > bPct
}
