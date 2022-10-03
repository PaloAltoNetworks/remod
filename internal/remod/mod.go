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

package remod

import (
	"fmt"
	"os"
	"strings"
)

// Enabled checks if remod is enabled.
func Enabled() bool {

	_, err := os.Stat(goModBackup())
	if err == nil {
		return true
	}

	if !os.IsNotExist(err) {
		panic(err)
	}

	return false
}

// Initialized checks if remod is initialized.
func Initialized() bool {

	_, err := os.Stat("remod.dev")
	if err == nil {
		return true
	}

	if !os.IsNotExist(err) {
		panic(err)
	}

	return false
}

// Install installs go mod in the current repo.
func Install(prefix string, version string, included []string, excluded []string) error {

	if Enabled() {
		return nil
	}

	if !strings.HasPrefix(prefix, ".") && version == "" {
		return fmt.Errorf("you must set --replace-version if --prefix is not local")
	}

	if strings.HasPrefix(prefix, ".") && version != "" {
		return fmt.Errorf("you must not set --replace-version if --prefix is local")
	}

	gomod, err := os.ReadFile("go.mod")
	if err != nil {
		return fmt.Errorf("unable to read go.mod: %s", err)
	}

	modules, err := Extract(gomod, included, excluded)
	if err != nil {
		return fmt.Errorf("unable to extract modules: %s", err)
	}

	odata := makeGoModDev(modules, prefix, version)
	if odata == nil {
		return nil
	}

	if err := os.WriteFile(goDev, odata, 0655); err != nil {
		return err
	}

	return nil
}

// On enables remod
func On() error {

	// we read the current state
	gomod, err := os.ReadFile("go.mod")
	if err != nil {
		return fmt.Errorf("unable to read go.mod: %s", err)
	}

	gosum, err := os.ReadFile("go.sum")
	if err != nil {
		return fmt.Errorf("unable to read go.sum: %s", err)
	}

	godev, err := os.ReadFile(goDev)
	if err != nil {
		return fmt.Errorf("unable to read %s: %s", goDev, err)
	}

	// we strip any previous remod replacements
	gomod = strip(gomod)

	mbak := goModBackup()
	sback := goSumBackup()

	if err := os.MkdirAll(".remod", 0700); err != nil {
		return fmt.Errorf("unable to create remod directory: %s", err)
	}

	if err := os.WriteFile(mbak, gomod, 0644); err != nil {
		return fmt.Errorf("unable to write %s: %s", mbak, err)
	}

	if err := os.WriteFile(sback, gosum, 0644); err != nil {
		return fmt.Errorf("unable to write %s: %s", sback, err)
	}

	if err := os.WriteFile("go.mod", assemble(gomod, prepareGoDev(godev)), 0644); err != nil {
		return fmt.Errorf("unable to write go.mod: %s", err)
	}

	if err := os.WriteFile("go.sum", gosum, 0644); err != nil {
		return fmt.Errorf("unable to write go.sum: %s", err)
	}

	return nil
}

// Off disables remod
func Off() error {

	if !Enabled() {
		return nil
	}

	mbak := goModBackup()
	sback := goSumBackup()

	gomod, err := os.ReadFile(mbak)
	if err != nil {
		return fmt.Errorf("unable to read go.mod: %s", err)
	}

	gosum, err := os.ReadFile(sback)
	if err != nil {
		return fmt.Errorf("unable to read go.sum: %s", err)
	}

	if err := os.RemoveAll(mbak); err != nil {
		return fmt.Errorf("unable to remove %s: %s", mbak, err)
	}

	if err := os.RemoveAll(sback); err != nil {
		return fmt.Errorf("unable to remove %s: %s", sback, err)
	}

	if err := os.WriteFile("go.mod", gomod, 0644); err != nil {
		return fmt.Errorf("unable to write %s: %s", mbak, err)
	}

	if err := os.WriteFile("go.sum", gosum, 0644); err != nil {
		return fmt.Errorf("unable to write %s: %s", sback, err)
	}

	return nil
}

// Uninstall uninstalls go mod from the current repo.
func Uninstall() error {

	if Enabled() {
		return fmt.Errorf("run remod off first")
	}

	if err := GitRemoveConfig(); err != nil {
		return fmt.Errorf("unable to restore git config: %s", err)
	}

	if err := os.RemoveAll(goDev); err != nil {
		return fmt.Errorf("unable to restore go.mod: %s", err)
	}

	return nil
}
