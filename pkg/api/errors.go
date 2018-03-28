package api

import (
	"github.com/google/go-github/github"
	"fmt"
	"os"
)

func Fail(err error) {
	switch e := err.(type) {
	case *github.RateLimitError, *github.AbuseRateLimitError:
		fmt.Printf("ERROR: Rate limit reached! %s\n", err.Error())

	case *github.ErrorResponse:
		fmt.Printf("ERROR: Github replied: %d %s\n", e.Response.StatusCode, e.Message)
		if e.DocumentationURL != "" {
			fmt.Printf("\tsee: %s\n", e.DocumentationURL)
		}

	default:
		fmt.Printf("ERROR: [%T] %s\n", err, err.Error())
	}

	os.Exit(1)
}