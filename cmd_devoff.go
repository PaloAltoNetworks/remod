package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmdDevoff = &cobra.Command{
	Use:     "off",
	Aliases: []string{"devoff"},
	Short:   "Remove developpment replace directive",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return viper.BindPFlags(cmd.Flags())
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		return os.RemoveAll("go.mod.dev")
	},
}
