package cmd

import (
	"github.com/spf13/cobra"
	"github.com/slushie/ghorg/pkg/output"
	"os"
)

// starsCmd represents the stars command
var starsCmd = &cobra.Command{
	Use:   "stars",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		showTable()
	},
}

func init() {
	rootCmd.AddCommand(starsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// starsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// starsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func showTable() {
	var recs = []output.Record{
		{"potato": "anytime", "eggs": "8am"},
		{"potato": "never", "tacos": "3am"},
	}

	recordWriter.WriteRecords(os.Stdout, recs, nil)
}