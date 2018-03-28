package api

import (
	"github.com/google/go-github/github"
	"fmt"
	"os"
)

func Fail(err error) {
	switch err.(type) {
	case *github.RateLimitError, *github.AbuseRateLimitError:
		fmt.Printf("ERROR: Rate limit reached! %s\n", err.Error())

	default:
		fmt.Printf("ERROR: %s\n", err.Error())
	}

	os.Exit(1)
}