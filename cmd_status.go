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
	"go.aporeto.io/remod/internal/remod"
)

var cmdStatus = &cobra.Command{
	Use:     "status",
	Aliases: []string{"stat", "st"},
	Short:   "Tells if remod is currently active on the branch",
	Long: `This command will print in stdout if remod is currently
on for the current branch.

It will return 0 if on, 1 otherwise.
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return viper.BindPFlags(cmd.Flags())
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		if remod.Enabled() {
			fmt.Println("remod is on")
			os.Exit(0)
		} else {
			fmt.Println("remod is off")
			os.Exit(1)
		}

		return nil
	},
}
