package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

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

		target := args[0]

		idata, err := ioutil.ReadFile(target)
		if err != nil {
			return fmt.Errorf("unable to read go.mod: %s", err)
		}

		if !remod.IsEnabled() {
			_, err = os.Stdout.Write(idata)
			return err
		}

		if strings.HasSuffix(target, "go.mod") {
			data, err := ioutil.ReadFile(remod.GoModBackup)
			if err != nil {
				return fmt.Errorf("unable to read %s: %s", remod.GoModBackup, err)
			}
			_, err = os.Stdout.Write(data)
			return err
		}

		if strings.HasSuffix(target, "go.sum") {
			data, err := ioutil.ReadFile(remod.GoSumBackup)
			if err != nil {
				return fmt.Errorf("unable to read %s: %s", remod.GoSumBackup, err)
			}
			_, err = os.Stdout.Write(data)
			return err
		}

		return fmt.Errorf("remod gitdiff ran on unsuported file: %s", target)

	},
}
