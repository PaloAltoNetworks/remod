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
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.aporeto.io/remod/internal/remod"
)

var preCommitHookBase = []byte(`.git/hooks/pre-commit-remod || exit 1`)

var preCommitHook = []byte(`#!/bin/bash
[ -f go.mod.bak ] && echo "cannot commit while hard remod is on: run 'remod off' first" && exit 1
`)

func init() {
	cmdDevon.Flags().StringSliceP("include", "m", nil, "Set the prefix of the modules to include")
	cmdDevon.Flags().StringSliceP("exclude", "E", nil, "Set the prefix of the modules to exclude")
	cmdDevon.Flags().StringP("prefix", "p", "../", "The prefix to use for the replacements")
	cmdDevon.Flags().String("replace-version", "", "Set the version to use for replacement. It must be set if prefix is not ../ and must not be if different")
}

var cmdDevon = &cobra.Command{
	Use:     "on",
	Aliases: []string{"devon"},
	Short:   "Apply developpment replace directive",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return viper.BindPFlags(cmd.Flags())
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		included := viper.GetStringSlice("include")
		excluded := viper.GetStringSlice("exclude")
		prefix := viper.GetString("prefix")
		version := viper.GetString("replace-version")

		if remod.IsHardMode() {
			return fmt.Errorf("remod hard is already on")
		}

		if !strings.HasPrefix(prefix, ".") && version == "" {
			return fmt.Errorf("you must set --replace-version if --prefix is not local")
		}

		if strings.HasPrefix(prefix, ".") && version != "" {
			return fmt.Errorf("you must not set --replace-version if --prefix is local")
		}

		gomod, err := ioutil.ReadFile("go.mod")
		if err != nil {
			return fmt.Errorf("unable to read go.mod: %s", err)
		}

		gosum, err := ioutil.ReadFile("go.sum")
		if err != nil {
			return fmt.Errorf("unable to read go.sum: %s", err)
		}

		modules, err := remod.Extract(gomod, included, excluded)
		if err != nil {
			return fmt.Errorf("unable to extract modules: %s", err)
		}

		odata, err := remod.MakeDevMod(gomod, modules, prefix, version)
		if err != nil {
			return fmt.Errorf("unable to apply dev replacements: %s", err)
		}
		if odata == nil {
			return nil
		}

		if err := ioutil.WriteFile("go.mod.dev", odata, 0655); err != nil {
			return err
		}

		if err := os.Rename("go.mod", "go.mod.bak"); err != nil {
			return fmt.Errorf("unable to backup go.mod: %s", err)
		}

		if err := ioutil.WriteFile("go.mod", append(gomod, odata...), 0644); err != nil {
			return fmt.Errorf("unable to write build go.mod: %s", err)
		}

		if err := ioutil.WriteFile("go.sum.bak", gosum, 0644); err != nil {
			return fmt.Errorf("unable to write go.sum.bak: %s", err)
		}

		// install the git hook
		if err := installPreCommitHook(); err != nil {
			return fmt.Errorf("unable to manage pre commit hook: %s", err)
		}

		return nil
	},
}

func installPreCommitHook() error {

	// check if there is a pre-commit hook
	_, err := os.Stat(".git/hooks/pre-commit")
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("unable to stat pre-commit: %s", err)
	}

	if err != nil {
		if err := ioutil.WriteFile(".git/hooks/pre-commit", append([]byte("#!/bin/bash\n"), preCommitHookBase...), 0755); err != nil {
			return fmt.Errorf("unable to write pre-commit hook: %s", err)
		}
	} else {
		precommit, err := ioutil.ReadFile(".git/hooks/pre-commit")
		if err != nil {
			return fmt.Errorf("unable to read pre-commit hook: %s", err)
		}

		if !bytes.Contains(precommit, preCommitHookBase) {
			if err := ioutil.WriteFile(".git/hooks/pre-commit", append(precommit, append([]byte("\n"), preCommitHookBase...)...), 0750); err != nil {
				return fmt.Errorf("unable to append pre-commit-remod: %s", err)
			}
		}
	}

	if err := ioutil.WriteFile(".git/hooks/pre-commit-remod", preCommitHook, 0750); err != nil {
		return fmt.Errorf("unable to write pre-commit hook: %s", err)
	}

	return nil
}
