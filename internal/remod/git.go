package remod

import (
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
)

// GitConfig installs the needed git config
func GitConfig() error {

	// cmd1 := exec.Command("git", "config", "diff.textconv", "remod gitdiff")
	// if err := cmd1.Run(); err != nil {
	// 	return fmt.Errorf("unable to update git config for diff.remod: %s", err)
	// }

	cmd2 := exec.Command("git", "config", "filter.clean", "remod gitclean %f")
	if err := cmd2.Run(); err != nil {
		return fmt.Errorf("unable to update git config for filter.clean: %s", err)
	}

	cmd3 := exec.Command("git", "config", "filter.smudge", "remod gitsmudge %f")
	if err := cmd3.Run(); err != nil {
		return fmt.Errorf("unable to update git config for filter.smudge: %s", err)
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

	if err := ioutil.WriteFile(".gitattributes", []byte("go.mod filter=remod\ngo.sum filter=remod\n"), 0644); err != nil {
		return fmt.Errorf("unable to write .gitattributes: %s", err)
	}

	return nil
}

// GitFilterClean is used by git filter.
func GitFilterClean(filename string, input io.Reader, output io.Writer) error {

	idata, err := ioutil.ReadAll(input)
	if err != nil {
		return fmt.Errorf("unable to read input: %s", err)
	}

	if !Enabled() {
		_, err = output.Write(idata)
		return err
	}

	switch filename {

	case "go.mod":

		gomod, err := ioutil.ReadFile(goModBackup)
		if err != nil {
			return fmt.Errorf("unable to update previous bak: %s", err)
		}

		_, err = output.Write(gomod)
		return err

	case "go.sum":

		gosum, err := ioutil.ReadFile(goSumBackup)
		if err != nil {
			return fmt.Errorf("unable to update previous bak: %s", err)
		}

		_, err = output.Write(gosum)
		return err
	}

	return nil
}

// GitFilterSmudge is used by git filter.
func GitFilterSmudge(filename string, input io.Reader, output io.Writer) error {

	idata, err := ioutil.ReadAll(input)
	if err != nil {
		return fmt.Errorf("unable to read input: %s", err)
	}

	if !Enabled() {
		_, err = output.Write(idata)
		return err
	}

	switch filename {

	case "go.mod":

		godev, err := ioutil.ReadFile(goDev)
		if err != nil {
			return fmt.Errorf("unable to read %s: %s", goDev, err)
		}

		_, err = output.Write(append(idata, append([]byte("\n"), godev...)...))
		return err

	case "go.sum":

		_, err = output.Write(idata)
		return err

	default:
		return fmt.Errorf("received non go.mod and non go.sum: %s", filename)
	}
}
