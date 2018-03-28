package cmd

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

var cfgFile string
var organization string
var accessToken string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "ghorg",
	Short: "A Github Organization Stats Tool",
	Long: `
 e88~~\  888   |                       / 
d888     888___|  e88~-_  888-~\ e88~88e 
8888 __  888   | d888   i 888    888 888 
8888   | 888   | 8888   | 888    "88_88" 
Y888   | 888   | Y888   ' 888     /      
 "88__/  888   |  "88_-~  888    Cb      
                                  Y8""8D 

This tool shows basic statistics for your Github organization.

You can set your organization via a command line option, environment variable,
or in the config file. 
`,
	PersistentPreRunE:  preRunE,
	PersistentPostRunE: postRunE,
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

	rootCmd.PersistentFlags().StringVarP(
		&organization,
		"organization",
		"N",
		"",
		"Organization name",
	)

	viper.BindPFlag("organization", rootCmd.PersistentFlags().Lookup("organization"))

	rootCmd.PersistentFlags().StringVarP(
		&accessToken,
		"access-token",
		"T",
		"",
		"Github OAuth2 access token used to authenticate REST calls.",
	)

	viper.BindPFlag("access-token", rootCmd.PersistentFlags().Lookup("access-token"))
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

	// Read in config silently
	viper.ReadInConfig()
}

func parseRootFlags(cmd *cobra.Command, args []string) error {
	organization = viper.GetString("organization")
	accessToken = viper.GetString("access-token")

	if len(args) == 1 {
		organization = args[0]
	}

	if organization == "" {
		return fmt.Errorf("missing required param: organization")
	}

	if accessToken == "" {
		fmt.Println("WARNING: No access token specified, you may hit rate limits from Github!")
	}

	return nil
}
