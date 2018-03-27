package cmd

import (
	"github.com/spf13/cobra"
	"github.com/slushie/ghorg/pkg/output"
	"os"
	"strings"
	"fmt"
	"github.com/slushie/ghorg/pkg/repos"
)

var recordWriter output.RecordWriter
var fields []string
var repoList *repos.List

// List flags are parsed and stored in this struct
var listOptions = struct {
	Count   uint
	Reverse bool
	Format  string
}{
	Count:   0,
	Reverse: false,
	Format:  "table",
}

// All list commands should support common flags.
func addListFlags(cmd *cobra.Command) {
	cmd.Flags().UintVarP(
		&listOptions.Count,
		"count",
		"c",
		5,
		"Number of repos to list",
	)

	cmd.Flags().BoolVarP(
		&listOptions.Reverse,
		"reverse",
		"R",
		false,
		"Sort in reverse order",
	)

	cmd.Flags().StringVarP(
		&listOptions.Format,
		"output-format",
		"F",
		"table",
		"Output format. One of: table, json",
	)

	// TODO expose arbitrary repo data in the output
	// TODO repo list options (eg, private repos)
}

// A pre-run hook that validates input and sets listOptions from flags.
func parseListFlags(cmd *cobra.Command, args []string) error {
	listOptions.Format = strings.ToLower(listOptions.Format)

	switch listOptions.Format {
	case "json":
		recordWriter = output.NewJson()
	case "table":
		recordWriter = output.NewTable()
	default:
		return fmt.Errorf("invalid output format: %s", listOptions.Format)
	}

	return nil
}

// A post-run hook that sorts and writes records after processing.
func outputRecords(cmd *cobra.Command, args []string) error {
	if repoList != nil {
		if listOptions.Reverse {
			repoList.SortReverse()
		} else {
			repoList.Sort()
		}

		records, err := repoList.Records(fields)
		if err != nil {
			return err
		}

		if listOptions.Count != 0 {
			records = records[0:listOptions.Count]
		}

		return recordWriter.WriteRecords(os.Stdout, records, fields)
	} else {
		return nil
	}
}
