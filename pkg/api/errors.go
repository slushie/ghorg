package api

import (
	"github.com/google/go-github/github"
	"fmt"
	"os"
	"io"
)

// make testing easier by parameterizing os.Exit
var failExit = os.Exit

// another testing shim
var failWriter io.Writer = os.Stdout

func Fail(err error) {
	switch e := err.(type) {
	case *github.RateLimitError, *github.AbuseRateLimitError:
		fmt.Fprintf(failWriter, "ERROR: Rate limit reached! %s\n", e.Error())

	case *github.ErrorResponse:
		fmt.Fprintf(failWriter, "ERROR: Github replied: %d %s\n", e.Response.StatusCode, e.Message)
		fmt.Fprintf(failWriter, "\tfrom: %s\n", e.Response.Request.URL)

		if e.DocumentationURL != "" {
			fmt.Fprintf(failWriter, "\tdocs: %s\n", e.DocumentationURL)
		}

	default:
		fmt.Fprintf(failWriter, "ERROR: [%T] %s\n", err, err.Error())
	}

	failExit(1)
}