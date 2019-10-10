package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.aporeto.io/remod/internal/remod"
)

var cmdGitDiff = &cobra.Command{
	Use:    "gitdiff",
	Short:  "Used by git attributes",
	Hidden: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return viper.BindPFlags(cmd.Flags())
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		idata, err := ioutil.ReadFile(args[0])
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

		_, err = os.Stdout.Write(odata)

		return err
	},
}
