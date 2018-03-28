package cmd

import (
	"github.com/google/go-github/github"
	"github.com/slushie/ghorg/pkg/output"
	"strings"
	"fmt"
	"strconv"
	"github.com/spf13/cobra"
	"context"
	"github.com/slushie/ghorg/pkg/api"
	"golang.org/x/time/rate"
	"sync"
)

const (
	OpenPRState   = "open"
	ClosedPRState = "closed"
	AllPRStates   = "all"
)

var (
	pullState   = AllPRStates
	pullWorkers = 10
	pullLimit   = 5.0
	pullCounts  RepoPullCounts
)

type RepoPullCounts map[int64]uint

func MarshalRepo(repo *github.Repository, fields []string) (output.Record, error) {
	rec := make(output.Record)
	for _, f := range fields {
		switch strings.ToLower(f) {
		case "stars":
			rec[f] = strconv.Itoa(repo.GetStargazersCount())
		case "forks":
			rec[f] = strconv.Itoa(repo.GetForksCount())
		case "name":
			rec[f] = repo.GetName()
		case "url":
			rec[f] = repo.GetHTMLURL()
		case "prs":
			c, _ := pullCounts[repo.GetID()]
			rec[f] = strconv.Itoa(int(c))
		case "contrib-%":
			prs, _ := pullCounts[repo.GetID()]
			forks := repo.GetForksCount()
			if prs == 0 {
				rec[f] = "-"
			} else if forks == 0 {
				rec[f] = "+"
			} else {
				pct := float64(prs) / float64(forks)
				rec[f] = fmt.Sprintf("%d%%", int(100.0*pct))
			}

		default:
			return nil, fmt.Errorf("unknown field %v", f)
		}
	}

	return rec, nil
}

func addPullFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(
		"state",
		"s",
		OpenPRState,
		"Pull request state. One of: open, closed, all",
	)

	cmd.Flags().IntVar(
		&pullWorkers,
		"max-workers",
		pullWorkers,
		"Number of simultaneous workers using Github for gathering stats",
	)

	cmd.Flags().Float64Var(
		&pullLimit,
		"rate-limit",
		pullLimit,
		"Number of requests per second to Github API",
	)
}

func parsePullFlags(cmd *cobra.Command, args []string) error {
	pullState, _ = cmd.Flags().GetString("state")
	pullState = strings.ToLower(pullState)

	switch pullState {
	case OpenPRState:
	case ClosedPRState:
	case AllPRStates:
	default:
		return fmt.Errorf("invalid pr state: %s", pullState)
	}

	pullWorkers, _ = cmd.Flags().GetInt("max-workers")
	if pullWorkers < 1 {
		return fmt.Errorf("invalid max workers: %d", pullWorkers)
	}

	pullLimit, _ = cmd.Flags().GetFloat64("rate-limit")
	if pullLimit < 1 {
		return fmt.Errorf("invalid rate limit: %f", pullLimit)
	}

	return nil
}

// Runs a work queue and waits for the results.
func countRepoPulls(ctx context.Context, c *api.Client, rs []*github.Repository) (RepoPullCounts, error) {
	type result struct {
		Repo  *github.Repository
		Count uint
		Err   error
	}

	jobs := make(chan *github.Repository, len(rs))
	results := make(chan result)
	state := pullState

	// load up the queue of work
	for _, r := range rs {
		jobs <- r
	}
	close(jobs)

	// token bucket with a burst size equal to the worker count
	limit := rate.NewLimiter(rate.Limit(pullLimit), pullWorkers)

	// start workers
	wg := &sync.WaitGroup{}
	for i := 0; i < pullWorkers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for repo := range jobs {
				// rate limit API calls
				if err := limit.Wait(ctx); err != nil {
					break // context cancelled or deadline expired
				}

				opt := github.PullRequestListOptions{State: state}
				count, err := c.CountRepositoryPRs(ctx, repo, &opt)
				results <- result{
					Repo:  repo,
					Count: count,
					Err:   err,
				}
			}
		}(i)
	}

	// wait for workers to finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// collect results
	pullCounts := make(RepoPullCounts)
	for res := range results {
		if res.Err != nil {
			return nil, res.Err
		}

		id := res.Repo.GetID()
		pullCounts[id] = res.Count
	}

	return pullCounts, nil
}
