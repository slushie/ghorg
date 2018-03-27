package cmd

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
	"github.com/slushie/ghorg/pkg/output"
)

var cfgFile string

var recordWriter output.RecordWriter

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ghorg",
	Short: "A Github Organization Stats Tool",
	Long: `This tool shows basic statistics for your Github organization.

You can set your organization via a command line option, environment variable,
or in the config file. 
`,
	PersistentPreRunE: preRun,
	PersistentPostRun: postRun,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(
		&cfgFile,
		"config",
		"",
		"Path to ghorg config file",
	)

	rootCmd.PersistentFlags().StringP(
		"organization",
		"N",
		"",
		"Organization name",
	)

	rootCmd.PersistentFlags().StringP(
		"output-format",
		"F",
		"table",
		"Output format. One of: table, json",
	)

	rootCmd.MarkFlagRequired("organization")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".ghorg" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".ghorg")
	}

	// Sane environment variable naming
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func preRun(cmd *cobra.Command, args []string) error {
	f, _ := cmd.Flags().GetString("output-format")

	switch strings.ToLower(f) {
	case "json":
		recordWriter = output.NewJson()
	case "table":
		recordWriter = output.NewTable()
	default:
		return fmt.Errorf("invalid output format: %s", f)
	}

	return nil
}

func postRun(cmd *cobra.Command, args []string) {

}