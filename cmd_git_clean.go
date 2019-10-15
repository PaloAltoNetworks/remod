package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.aporeto.io/remod/internal/remod"
)

func makeRepoSmudge() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return path.Join(os.TempDir(), strings.ReplaceAll(strings.ReplaceAll(dir, "/", "_"), "\\", "_"))
}

var cmdGitClean = &cobra.Command{
	Use:    "gitclean",
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

			gomod, err := ioutil.ReadFile(remod.GoModBackup)
			if err != nil {
				return fmt.Errorf("unable to update previous bak: %s", err)
			}

			_, err = os.Stdout.Write(gomod)
			return err

		case "go.sum":

			gosum, err := ioutil.ReadFile(remod.GoSumBackup)
			if err != nil {
				return fmt.Errorf("unable to update previous bak: %s", err)
			}

			_, err = os.Stdout.Write(gosum)
			return err
		}

		return nil
	},
}
