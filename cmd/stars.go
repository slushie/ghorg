package cmd

import (
	"github.com/spf13/cobra"
	"github.com/slushie/ghorg/pkg/api"
	"context"
	"github.com/slushie/ghorg/pkg/repos"
	"github.com/google/go-github/github"
	"github.com/slushie/ghorg/pkg/output"
	"strings"
	"strconv"
	"fmt"
)

// starsCmd represents the stars command
var starsCmd = &cobra.Command{
	Use:   "stars",
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
		panic(err.Error())
	}

	// these are the relevant fields
	fields = []string{"Stars", "Name", "URL"}

	// set up the output list
	repoList = repos.NewList()
	repoList.Add(rs...)
	repoList.Marshal = MarshalRepoStars
	repoList.Compare = CompareRepoStars
}

func CompareRepoStars(a, b *github.Repository) bool {
	return a.GetStargazersCount() < b.GetStargazersCount()
}

func MarshalRepoStars(repo *github.Repository, fields []string) (output.Record, error) {
	rec := make(output.Record)
	for _, f := range fields {
		switch strings.ToLower(f) {
		case "stars":
			rec[f] = strconv.Itoa(repo.GetStargazersCount())
		case "name":
			rec[f] = repo.GetName()
		case "url":
			rec[f] = repo.GetURL()
		default:
			return nil, fmt.Errorf("unknown field %v", f)
		}
	}

	return rec, nil
}