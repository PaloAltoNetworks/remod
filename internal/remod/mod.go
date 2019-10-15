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
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// IsEnabled checks if remod is enabled.
func IsEnabled() bool {

	_, err := os.Stat(GoModBackup)
	if err == nil {
		return true
	}

	if !os.IsNotExist(err) {
		panic(err)
	}

	return false
}

// MakeDevMod builds the dev mod file.
func MakeDevMod(data []byte, modules []string, base string, version string) ([]byte, error) {

	if len(modules) == 0 {
		return data, nil
	}

	if version != "" {
		version = " " + version
	}

	buf := bytes.NewBuffer(nil)

	must(buf.WriteString("replace (\n"))
	for _, m := range modules {
		must(buf.WriteString(fmt.Sprintf("\t%s => %s%s%s\n", m, base, filepath.Base(m), version)))
	}
	must(buf.WriteString(")\n"))

	return append(bytes.TrimSpace(buf.Bytes()), '\n'), nil
}

// On enables remod
func On() error {

	if IsEnabled() {
		return fmt.Errorf("remod is already on")
	}

	if err := os.MkdirAll(".remod", 0700); err != nil {
		return fmt.Errorf("unable to create remod directory: %s", err)
	}

	gomod, err := ioutil.ReadFile("go.mod")
	if err != nil {
		return fmt.Errorf("unable to read go.mod: %s", err)
	}

	gosum, err := ioutil.ReadFile("go.sum")
	if err != nil {
		return fmt.Errorf("unable to read go.sum: %s", err)
	}

	godev, err := ioutil.ReadFile(GoDev)
	if err != nil {
		return fmt.Errorf("unable to read %s: %s", GoDev, err)
	}

	if err := ioutil.WriteFile(GoModBackup, gomod, 0644); err != nil {
		return fmt.Errorf("unable to write %s: %s", GoModBackup, err)
	}

	if err := ioutil.WriteFile(GoSumBackup, gosum, 0644); err != nil {
		return fmt.Errorf("unable to write %s: %s", GoSumBackup, err)
	}

	if err := ioutil.WriteFile("go.mod", append(gomod, append([]byte("\n"), godev...)...), 0644); err != nil {
		return fmt.Errorf("unable to write go.mod: %s", err)
	}

	if err := ioutil.WriteFile("go.sum", gosum, 0644); err != nil {
		return fmt.Errorf("unable to write go.sum: %s", err)
	}

	return nil
}

// Off disables remod
func Off() error {

	if !IsEnabled() {
		return fmt.Errorf("remod is not enabled")
	}

	gomod, err := ioutil.ReadFile(GoModBackup)
	if err != nil {
		return fmt.Errorf("unable to read go.mod: %s", err)
	}

	gosum, err := ioutil.ReadFile(GoSumBackup)
	if err != nil {
		return fmt.Errorf("unable to read go.sum: %s", err)
	}

	if err := os.Remove(GoModBackup); err != nil {
		return fmt.Errorf("unable to remove %s: %s", GoModBackup, err)
	}

	if err := os.Remove(GoSumBackup); err != nil {
		return fmt.Errorf("unable to remove %s: %s", GoSumBackup, err)
	}

	if err := ioutil.WriteFile("go.mod", gomod, 0644); err != nil {
		return fmt.Errorf("unable to write %s: %s", GoModBackup, err)
	}

	if err := ioutil.WriteFile("go.sum", gosum, 0644); err != nil {
		return fmt.Errorf("unable to write %s: %s", GoModBackup, err)
	}

	return nil
}

func must(n int, err error) {
	if err != nil {
		panic(err)
	}
}
