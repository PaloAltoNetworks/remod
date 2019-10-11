package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.aporeto.io/remod/internal/remod"
)

var cmdGo = &cobra.Command{
	Use:                "go",
	Aliases:            []string{"g"},
	Short:              "Run a go command wrapped so it used go.mod.dev",
	DisableFlagParsing: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return viper.BindPFlags(cmd.Flags())
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		return remod.WrapGoCommand(args...)
	},
}
