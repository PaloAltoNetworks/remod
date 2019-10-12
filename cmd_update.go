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
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.aporeto.io/remod/internal/remod"
)

func init() {
	cmdUpdate.Flags().StringSliceP("include", "m", nil, "Set the prefix of the modules to include")
	cmdUpdate.Flags().StringSliceP("exclude", "E", nil, "Set the prefix of the modules to exclude")
	cmdUpdate.Flags().StringP("folder", "f", "./", "Set the path to the folder file")
	cmdUpdate.Flags().BoolP("recursive", "r", false, "If true, remod will look for mod files in given --folder and all 1 level subfolders")
	cmdUpdate.Flags().String("version", "latest", "Set to which version you want to update the matching modules")
}

var cmdUpdate = &cobra.Command{
	Use:     "update",
	Aliases: []string{"up"},
	Short:   "Update the modules in the given path",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return viper.BindPFlags(cmd.Flags())
	},
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		folder := "."
		if len(args) > 0 {
			folder = args[0]
		}

		if remod.IsHardMode() {
			return fmt.Errorf("unable to up when remod is on in hard mode. run remod off first")
		}

		recursive := viper.GetBool("recursive")
		included := viper.GetStringSlice("include")
		excluded := viper.GetStringSlice("exclude")
		version := viper.GetString("version")

		var paths []string
		if recursive && folder != "." {

			items, err := ioutil.ReadDir(folder)
			if err != nil {
				return fmt.Errorf("unable to list content of dir: %s", err)
			}

			for _, item := range items {
				if !item.IsDir() {
					continue
				}

				p := path.Join(folder, item.Name(), "go.mod")
				_, err := os.Stat(p)
				if err != nil {
					if os.IsNotExist(err) {
						continue
					}
					return fmt.Errorf("unable stat path '%s': %s", p, err)
				}

				paths = append(paths, p)
			}
		} else {
			paths = append(paths, path.Join(folder, "go.mod"))
		}

		for _, p := range paths {

			basedir := filepath.Dir(p)
			if err := os.Chdir(basedir); err != nil {
				return fmt.Errorf("unable to move to %s: %s", basedir, err)
			}

			if len(paths) > 1 {
				fmt.Println("* Entering", basedir)
			}

			idata, err := ioutil.ReadFile(p)
			if err != nil {
				return fmt.Errorf("unable to read go.mod: %s", err)
			}

			modules, err := remod.Extract(idata, included, excluded)
			if err != nil {
				return fmt.Errorf("unable to extract modules: %s", err)
			}

			if err := remod.Update(modules, version); err != nil {
				return fmt.Errorf("unable to extract modules: %s", err)
			}
		}

		return nil
	},
}
