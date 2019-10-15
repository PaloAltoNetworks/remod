package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.aporeto.io/remod/internal/remod"
)

var cmdGitSmudge = &cobra.Command{
	Use:    "gitsmudge",
	Short:  "Used by git attributes",
	Hidden: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return viper.BindPFlags(cmd.Flags())
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		idata, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("unable to read input: %s", err)
		}

		if !remod.IsEnabled() {
			_, err = os.Stdout.Write(idata)
			return err
		}

		switch args[0] {

		case "go.mod":

			godev, err := ioutil.ReadFile(remod.GoDev)
			if err != nil {
				return fmt.Errorf("unable to read %s: %s", remod.GoDev, err)
			}

			_, err = os.Stdout.Write(append(idata, append([]byte("\n"), godev...)...))
			return err

		case "go.sum":

			_, err = os.Stdout.Write(idata)
			return err

		default:
			return fmt.Errorf("received non go.mod and non go.sum: %s", args)
		}
	},
}
