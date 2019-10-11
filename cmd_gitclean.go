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

		idata, err := ioutil.ReadFile("go.mod")
		if err != nil {
			return fmt.Errorf("unable to read go.mod: %s", err)
		}

		cdata, err := remod.Get(idata)
		if err != nil {
			return fmt.Errorf("unable to get dev replacements: %s", err)
		}
		if cdata != nil {
			if err := ioutil.WriteFile(makeRepoSmudge(), cdata, 0644); err != nil {
				return fmt.Errorf("unable to write data for smudging: %s", err)
			}
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
