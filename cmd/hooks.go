package cmd

import "github.com/spf13/cobra"

type runHook func(cmd *cobra.Command, args []string) error

var preRunHooks = []runHook{
	ensureOrganization,
	parseListFlags,
}

var postRunHooks = []runHook{
	outputRecords,
}

func preRunE(cmd *cobra.Command, args []string) error {
	return runHooks(preRunHooks, cmd, args)
}

func postRunE(cmd *cobra.Command, args []string) error {
	return runHooks(postRunHooks, cmd, args)
}

// Run all hooks, returning the first error or nil.
func runHooks(hooks []runHook, cmd *cobra.Command, args []string) error {
	for _, hook := range hooks {
		err := hook(cmd, args)
		if err != nil {
			return err
		}
	}

	return nil
}