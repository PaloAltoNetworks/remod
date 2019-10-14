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
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmdDevoff = &cobra.Command{
	Use:     "off",
	Aliases: []string{"devoff"},
	Short:   "Remove developpment replace directive",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return viper.BindPFlags(cmd.Flags())
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		if err := os.Rename("go.mod.bak", "go.mod"); err != nil {
			return fmt.Errorf("unable to restore go.mod: %s", err)
		}

		if err := os.Rename("go.sum.bak", "go.sum"); err != nil {
			return fmt.Errorf("unable to restore go.mod: %s", err)
		}

		return nil
	},
}
