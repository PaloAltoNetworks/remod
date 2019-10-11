package main

import (
	"fmt"
	"os"

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

		if _, err := os.Stat(".gitattributes"); err != nil {
			if !os.IsNotExist(err) {
				return fmt.Errorf("unable to check for .gitattributes: %s", err)
			}
		} else {
			fmt.Println("warning: .gitattributes file already exists. You can add the relevant remod part by running:")
			fmt.Println("")
			fmt.Println("   echo 'go.mod diff=remod filter=remod' >> .gitattributes")
			fmt.Println("")
			return nil
		}

		if err := remod.GitInit(); err != nil {
			return err
		}

		if err := remod.GitConfig(); err != nil {
			return err
		}

		return nil
	},
}
