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
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.aporeto.io/remod/internal/remod"
)

var cmdGet = &cobra.Command{
	Use:                "get",
	Aliases:            []string{"g"},
	Short:              "Run a wrapper go get command",
	DisableFlagParsing: true,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return viper.BindPFlags(cmd.Flags())
	},
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(args) == 0 {
			return fmt.Errorf("you must at least pass one argument")
		}

		if err := remod.Off(); err != nil {
			return fmt.Errorf("unable to set remod to off: %s", err)
		}

		defer func() {
			if err := remod.On(); err != nil {
				panic(err)
			}
		}()

		c := exec.Command("go", append([]string{"get"}, args...)...)
		c.Stdin = os.Stdin
		c.Stderr = os.Stderr
		c.Stdout = os.Stdout

		if err := c.Start(); err != nil {
			return err
		}

		return c.Wait()
	},
}
