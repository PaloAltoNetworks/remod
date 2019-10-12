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
	"os/exec"
	"path/filepath"
)

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

// WrapGoCommand wraps a go command into a workable remod version.
func WrapGoCommand(args ...string) error {

	var exists bool

	_, err := os.Stat("go.mod.dev")
	if err == nil {
		exists = true
	}

	cmd := exec.Command("go", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if !exists {
		return cmd.Run()
	}

	// read go.mod
	gomod, err := ioutil.ReadFile("go.mod")
	if err != nil {
		return fmt.Errorf("unable to read go.mod: %s", err)
	}

	// read go.sum
	gosum, err := ioutil.ReadFile("go.sum")
	if err != nil {
		return fmt.Errorf("unable to read go.sum: %s", err)
	}

	// read go.sum.dev
	gomoddev, err := ioutil.ReadFile("go.mod.dev")
	if err != nil {
		return fmt.Errorf("unable to read go.mod.dev: %s", err)
	}

	// remove go.mod and defer its recreation.
	if err := os.RemoveAll("go.mod"); err != nil {
		return fmt.Errorf("unable to remove go.mod: %s", err)
	}
	defer ioutil.WriteFile("go.mod", gomod, 0644) // nolint

	// remove go.sum and defer its recreation.
	if err := os.RemoveAll("go.sum"); err != nil {
		return fmt.Errorf("unable to remove go.sum: %s", err)
	}
	defer ioutil.WriteFile("go.sum", gosum, 0644) // nolint

	// write new combined go.mod
	if err := ioutil.WriteFile("go.mod", append(gomod, gomoddev...), 0644); err != nil {
		return fmt.Errorf("unable to write build go.mod: %s", err)
	}

	// run the command
	return cmd.Run()
}

// IsHardMode checks if go.mod.bak is present.
func IsHardMode() bool {

	_, err := os.Stat("go.mod.bak")
	if err == nil {
		return true
	}

	if !os.IsNotExist(err) {
		panic(err)
	}

	return false
}

func must(n int, err error) {
	if err != nil {
		panic(err)
	}
}
