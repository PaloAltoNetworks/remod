package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
			return fmt.Errorf("unable to read go.mod: %s", err)
		}

		smudgefile := makeRepoSmudge()
		defer os.RemoveAll(smudgefile)

		odata, err := ioutil.ReadFile(smudgefile)
		if err != nil {
			if !os.IsNotExist(err) {
				return fmt.Errorf("unable to extract remod data: %s", err)
			}
		}

		out := idata
		if len(odata) > 0 {
			out = append(out, append([]byte("\n"), odata...)...)
		}

		_, err = os.Stdout.Write(out)

		return err
	},
}
