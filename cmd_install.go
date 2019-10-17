// Copyright 2019 Aporeto Inc.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.aporeto.io/remod/internal/remod"
)

func init() {
	cmdInstall.Flags().StringSliceP("include", "m", nil, "Set the prefix of the modules to include")
	cmdInstall.Flags().StringSliceP("exclude", "E", nil, "Set the prefix of the modules to exclude")
	cmdInstall.Flags().StringP("prefix", "p", "../", "The prefix to use for the replacements")
	cmdInstall.Flags().String("replace-version", "", "Set the version to use for replacement. It must be set if prefix is not ../ and must not be if different")
}

var cmdInstall = &cobra.Command{
	Use:     "init",
	Aliases: []string{"i", "install"},
	Short:   "Initializes the repository for remod",
	Long: `This command prepares and align your repository to work with remod.

It will:
- Prepare the '.gitattribute' file if needed
- Prepare the '.gitignore' file if needed
- Prepare the '.git/config' filter if needed

It will also create the 'remod.dev' file if needed with the eventual replacements
provided through the '--include' and '--exclude' flags. If no replacements are provided,
the remod.dev file will be empty. You can then add your replacements manually.

The flag '--replace-version' allows to pass a version to use for replacement.
The flag '--prefix' allows to set the replacement prefix. If the replacement prefix
if local (starts with '.') the '--replace-version' cannot be set. If it is remote,
'--replace-version' must be set.

Once the repository is initialized or a change has been made to the 'remod.dev' file
you need to run 'remod on'.
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return viper.BindPFlags(cmd.Flags())
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		if err := remod.Install(
			viper.GetString("prefix"),
			viper.GetString("replace-version"),
			viper.GetStringSlice("include"),
			viper.GetStringSlice("exclude"),
		); err != nil {
			return err
		}

		if err := remod.GitConfig(); err != nil {
			return err
		}

		if err := remod.GitInit(); err != nil {
			return err
		}

		return remod.GitAdd()
	},
}
