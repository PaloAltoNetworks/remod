package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.aporeto.io/remod/internal/remod"
)

var cmdGitInit = &cobra.Command{
	Use:   "init",
	Short: "Initializes the git attribute file",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return viper.BindPFlags(cmd.Flags())
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		if err := remod.GitInit(); err != nil {
			return err
		}

		if err := remod.GitConfig(); err != nil {
			return err
		}

		return nil
	},
}
