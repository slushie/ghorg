package cmd

import (
	"github.com/google/go-github/github"
	"github.com/slushie/ghorg/pkg/output"
	"strings"
	"fmt"
	"strconv"
)

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
		default:
			return nil, fmt.Errorf("unknown field %v", f)
		}
	}

	return rec, nil
}
