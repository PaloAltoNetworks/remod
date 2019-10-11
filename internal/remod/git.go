package remod

import (
	"fmt"
	"io/ioutil"
	"os/exec"
)

// GitConfig installs the needed git config
func GitConfig() error {

	cmd1 := exec.Command("git", "config", "diff.remod.textconv", "remod gitdiff")
	if err := cmd1.Run(); err != nil {
		return fmt.Errorf("unable to update git config for diff.remod: %s", err)
	}

	cmd2 := exec.Command("git", "config", "filter.remod.clean", "remod gitclean")
	if err := cmd2.Run(); err != nil {
		return fmt.Errorf("unable to update git config for filter.remod: %s", err)
	}

	return nil
}

// GitAdd adds the go mod
func GitAdd() error {

	cmd := exec.Command("git", "add", "go.mod")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("unable to git add go.mod: %s", err)
	}

	return nil
}

// GitInit initializes the .gitattribute file.
func GitInit() error {

	if err := ioutil.WriteFile(".gitattributes", []byte("go.mod diff=remod filter=remod\n"), 0644); err != nil {
		return fmt.Errorf("unable to write .gitattributes: %s", err)
	}

	return nil
}
