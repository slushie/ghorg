package api

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"errors"
	"github.com/google/go-github/github"
	"bytes"
	"net/http"
)

func TestFail(t *testing.T) {
	Convey(".Fail(err error)", t, func() {
		originalExit := failExit
		exitCode := 0
		exitCalled := false
		failExit = func(code int) {
			exitCalled = true
			exitCode = code
		}

		originalWriter := failWriter

		Convey("Calls os.Exit(1)", func() {
			Fail(errors.New("test-error"))
			So(exitCalled, ShouldBeTrue)
			So(exitCode, ShouldEqual, 1)
		})

		Convey("Handles RateLimitError", func() {
			err := &github.RateLimitError{
				Message: "test-rate-limit",
				Response: &http.Response{
					Request: &http.Request{},
				},
			}

			var buf bytes.Buffer
			failWriter = &buf
			Fail(err)

			output := buf.String()
			So(output, ShouldContainSubstring, "Rate limit reached")
			So(output, ShouldContainSubstring, "test-rate-limit")
		})

		Convey("Handles ErrorResponse", func() {
			err := &github.ErrorResponse{
				Message: "test-error-response",
				DocumentationURL: "test-documentation-url",
				Response: &http.Response{
					Request: &http.Request{},
				},
			}

			var buf bytes.Buffer
			failWriter = &buf
			Fail(err)

			output := buf.String()
			So(output, ShouldContainSubstring, "Github replied")
			So(output, ShouldContainSubstring, "test-error-response")

			Convey("Shows Documentation URL", func() {
				So(output, ShouldContainSubstring, "test-documentation-url")
			})
		})

		Convey("Handles other errors", func() {
			err := errors.New("test-generic-error")

			var buf bytes.Buffer
			failWriter = &buf
			Fail(err)

			output := buf.String()
			So(output, ShouldContainSubstring, "test-generic-error")
		})

		failExit = originalExit
		failWriter = originalWriter
	})

}