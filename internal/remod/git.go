package remod

import (
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"strings"
)

// GitConfig installs the needed git config
func GitConfig() error {

	// cmd1 := exec.Command("git", "config", "diff.textconv", "remod gitdiff")
	// if err := cmd1.Run(); err != nil {
	// 	return fmt.Errorf("unable to update git config for diff.remod: %s", err)
	// }

	cmd2 := exec.Command("git", "config", "filter.remod.clean", "remod gitclean %f")
	if err := cmd2.Run(); err != nil {
		return fmt.Errorf("unable to update git config for filter.clean: %s", err)
	}

	cmd3 := exec.Command("git", "config", "filter.remod.smudge", "remod gitsmudge %f")
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
		must(output.Write(idata))
		return nil
	}

	switch filename {

	case "go.mod":

		mbak := goModBackup()
		gomod, err := ioutil.ReadFile(mbak)
		if err != nil {
			return fmt.Errorf("unable to read %s: %s", mbak, err)
		}

		must(output.Write(gomod))

		return nil

	case "go.sum":

		sbak := goSumBackup()
		gosum, err := ioutil.ReadFile(sbak)
		if err != nil {
			return fmt.Errorf("unable to read %s: %s", sbak, err)
		}

		must(output.Write(gosum))

		return nil

	default:

		panic(fmt.Errorf("received non go.mod and non go.sum: %s", filename))
	}
}

// GitFilterSmudge is used by git filter.
func GitFilterSmudge(filename string, input io.Reader, output io.Writer) error {

	idata, err := ioutil.ReadAll(input)
	if err != nil {
		return fmt.Errorf("unable to read input: %s", err)
	}

	if !Enabled() {
		must(output.Write(idata))
		return nil
	}

	switch filename {

	case "go.mod":

		godev, err := ioutil.ReadFile(goDev)
		if err != nil {
			return fmt.Errorf("unable to read %s: %s", goDev, err)
		}

		must(output.Write(assemble(idata, prepareGoDev(godev))))

		return nil

	case "go.sum":

		must(output.Write(idata))

		return nil

	default:

		panic(fmt.Errorf("received non go.mod and non go.sum: %s", filename))
	}
}

func branchName() (string, error) {

	out, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return "", err
	}

	return strings.Trim(strings.ReplaceAll(strings.TrimPrefix(string(out), "heads/"), "/", "_"), "\n\t"), nil
}
