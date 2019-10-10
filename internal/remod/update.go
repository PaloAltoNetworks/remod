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

	return nil
}
