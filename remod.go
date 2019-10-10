package remod

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

// Extract extracts the required modules with the given prefixes
func Extract(file io.Reader, prefixes [][]byte) ([]string, error) {

	scanner := bufio.NewScanner(file)
	singleRequireStartPrefix := []byte("require ")
	multiRequireStartPrefix := []byte("require (")
	multiRequireEndPrefix := []byte(")")

	modules := map[string]struct{}{}

	var multiStart bool
	for scanner.Scan() {

		line := scanner.Bytes()

		if bytes.HasPrefix(line, multiRequireStartPrefix) {
			multiStart = true
			continue
		} else if bytes.HasPrefix(line, singleRequireStartPrefix) {
			mod := bytes.TrimSpace(line)
			mod = bytes.Replace(line, singleRequireStartPrefix, nil, -1)
			modules[string(bytes.SplitN(mod, []byte(" "), 2)[0])] = struct{}{}
			continue
		}

		if multiStart && bytes.HasPrefix(line, multiRequireEndPrefix) {
			multiStart = false
			continue
		}

		if multiStart {
			mod := bytes.TrimSpace(line)

			var found bool
			for _, prefix := range prefixes {
				if bytes.HasPrefix(mod, prefix) {
					found = true
				}
			}

			if !found {
				continue
			}

			modules[string(bytes.SplitN(mod, []byte(" "), 2)[0])] = struct{}{}
		}
	}

	out := make([]string, len(modules))
	var i int
	for mod := range modules {
		out[i] = mod
		i++
	}

	return out, nil
}

// Update will run a go get on the given modules at the given version
func Update(modules []string, version string) error {

	if err := os.Setenv("GO111MODULE", "on"); err != nil {
		return fmt.Errorf("unable set GO111MODULE variable: %s", err)
	}

	for _, mod := range modules {

		mod = fmt.Sprintf("%s@%s", mod, version)

		fmt.Printf("- Getting %s\n", mod)

		cmd := exec.Command("go", "get", mod)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("unable to run go get command: %s", err)
		}
	}

	return nil
}

// Enable will enable remod dev replacements for the given modules.
func Enable(path string, modules []string, base string) ([]byte, error) {

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to read go.mod: %s", err)
	}

	if bytes.Contains(data, []byte("replace ( // remod:replacements")) {
		return nil, nil
	}

	buf := bytes.NewBuffer(data)

	buf.WriteString("\nreplace ( // remod:replacements\n")
	for _, m := range modules {
		buf.WriteString(fmt.Sprintf("\t%s => %s%s\n", m, base, filepath.Base(m)))
	}
	buf.WriteString(")\n")

	return append(bytes.TrimSpace(buf.Bytes()), '\n'), nil
}

// Disable will disable remod dev replacements.
func Disable(path string) ([]byte, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to read go.mod: %s", err)
	}

	scanner := bufio.NewScanner(file)

	buf := bytes.NewBuffer(nil)

	start := []byte("replace ( // remod:replacements")
	end := []byte(")")

	var strip bool
	for scanner.Scan() {

		line := scanner.Bytes()

		if bytes.Equal(line, start) {
			strip = true
			continue
		}

		if strip && bytes.HasPrefix(line, end) {
			strip = false
			continue
		}

		if strip {
			continue
		}

		buf.Write(line)
		buf.WriteByte('\n')
	}

	return append(bytes.TrimSpace(buf.Bytes()), '\n'), nil
}
