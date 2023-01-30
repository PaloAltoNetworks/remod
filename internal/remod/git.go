package remod

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// GitConfig installs the needed git config
func GitConfig() error {

	cmd1 := exec.Command("git", "config", "filter.remod.clean", "remod gitclean")
	if err := cmd1.Run(); err != nil {
		return fmt.Errorf("unable to update git config for filter.clean: %s", err)
	}

	cm2 := exec.Command("git", "config", "filter.remod.smudge", "remod gitsmudge")
	if err := cm2.Run(); err != nil {
		return fmt.Errorf("unable to update git config for filter.smudge: %s", err)
	}

	return nil
}

// GitRemoveConfig removes the remod git config
func GitRemoveConfig() error {

	cmd1 := exec.Command("git", "config", "--unset", "filter.remod.clean", "remod gitclean")
	if err := cmd1.Run(); err != nil {
		return fmt.Errorf("unable to update git config for filter.clean: %s", err)
	}

	cm2 := exec.Command("git", "config", "--unset", "filter.remod.smudge", "remod gitsmudge")
	if err := cm2.Run(); err != nil {
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

	idata, err := os.ReadFile(".gitattributes")
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if !bytes.Contains(idata, []byte("go.mod filter=remod")) {
		idata = append(idata, []byte("go.mod filter=remod\n")...)
	}

	if !bytes.Contains(idata, []byte("go.sum filter=remod")) {
		idata = append(idata, []byte("go.sum filter=remod\n")...)
	}

	if err := os.WriteFile(".gitattributes", idata, 0644); err != nil {
		return fmt.Errorf("unable to write .gitattributes: %s", err)
	}

	idata, err = os.ReadFile(".gitignore")
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if !bytes.Contains(idata, []byte("remod.dev")) {
		idata = append(idata, []byte("remod.dev\n")...)
	}

	if !bytes.Contains(idata, []byte(".remod")) {
		idata = append(idata, []byte(".remod\n")...)
	}

	if err := os.WriteFile(".gitignore", idata, 0644); err != nil {
		return fmt.Errorf("unable to write .gitignore: %s", err)
	}

	return nil
}

// GitFilterClean is used by git filter.
func GitFilterClean(input io.Reader, output io.Writer) error {

	idata, err := io.ReadAll(input)
	if err != nil {
		return fmt.Errorf("unable to read input: %s", err)
	}

	if !Enabled() {
		must(output.Write(strip(idata)))
		return nil
	}

	if bytes.Contains(idata, []byte("module ")) {

		mbak := goModBackup()
		gomod, err := os.ReadFile(mbak)
		if err != nil {
			return fmt.Errorf("unable to read %s: %s", mbak, err)
		}

		must(output.Write(gomod))

		return nil

	} else if bytes.Contains(idata, []byte(" h1:")) {

		sbak := goSumBackup()
		gosum, err := os.ReadFile(sbak)
		if err != nil {
			return fmt.Errorf("unable to read %s: %s", sbak, err)
		}

		must(output.Write(gosum))

		return nil

	}

	panic(fmt.Errorf("received non go.mod and non go.sum: %s", string(idata)))
}

// GitFilterSmudge is used by git filter.
func GitFilterSmudge(input io.Reader, output io.Writer) error {

	idata, err := io.ReadAll(input)
	if err != nil {
		return fmt.Errorf("unable to read input: %s", err)
	}

	if !Enabled() {
		must(output.Write(idata))
		return nil
	}

	if bytes.Contains(idata, []byte("module ")) {

		godev, err := os.ReadFile(goDev)

		if err != nil {
			return fmt.Errorf("unable to read %s: %s", goDev, err)
		}

		must(output.Write(assemble(idata, prepareGoDev(godev))))

		return nil

	} else if bytes.Contains(idata, []byte(" h1:")) {

		must(output.Write(idata))

		return nil
	}

	panic(fmt.Errorf("received non go.mod and non go.sum: %s", string(idata)))
}

func branchName() (string, error) {

	out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return "", err
	}

	return strings.Trim(strings.ReplaceAll(strings.TrimPrefix(string(out), "heads/"), "/", "_"), "\n\t"), nil
}
