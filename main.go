package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {

	cobra.OnInitialize(func() {
		viper.SetEnvPrefix("remod")
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	})

	var rootCmd = &cobra.Command{
		Use:           "remod",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	rootCmd.AddCommand(
		cmdUpdate,
		cmdDevon,
		cmdDevoff,
		cmdGitClean,
		cmdGitDiff,
		cmdGitInit,
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
}
