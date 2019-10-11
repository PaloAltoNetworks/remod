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
	"os/exec"
)

// Update will run a go get on the given modules at the given version
func Update(modules []string, version string) error {

	if err := os.Setenv("GO111MODULE", "on"); err != nil {
		return fmt.Errorf("unable set GO111MODULE variable: %s", err)
	}

	for _, mod := range modules {

		mod = fmt.Sprintf("%s@%s", mod, version)

		cmd := exec.Command("go", "get", mod)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("unable to run go get command: %s", err)
		}
	}

	cmd := exec.Command("go", "mod", "tidy")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("unable to run go mod tidy: %s", err)
	}

	return nil
}
