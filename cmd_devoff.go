package main

import (
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.aporeto.io/remod/internal/remod"
)

var cmdDevoff = &cobra.Command{
	Use:     "devoff",
	Aliases: []string{"off"},
	Short:   "Remove developpment replace directive",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return viper.BindPFlags(cmd.Flags())
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		idata, err := ioutil.ReadFile("go.mod")
		if err != nil {
			return fmt.Errorf("unable to read go.mod: %s", err)
		}

		odata, err := remod.Disable(idata)
		if err != nil {
			return fmt.Errorf("unable to remove dev replacements: %s", err)
		}
		if odata == nil {
			return nil
		}

		if err := ioutil.WriteFile("go.mod", odata, 0655); err != nil {
			return err
		}

		return remod.GitAdd()
	},
}
