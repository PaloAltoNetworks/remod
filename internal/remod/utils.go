package remod

import (
	"bufio"
	"bytes"
	"fmt"
	"path/filepath"
)

// makeGoModDev builds the dev mod file.
func makeGoModDev(data []byte, modules []string, base string, version string) ([]byte, error) {

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

func must(n int, err error) {
	if err != nil {
		panic(err)
	}
}

func strip(in []byte) ([]byte, error) {

	scanner := bufio.NewScanner(bytes.NewBuffer(in))

	buf := bytes.NewBuffer(nil)

	var skip bool
	var last []byte
	for scanner.Scan() {

		line := scanner.Bytes()

		if bytes.Equal(append(line, '\n'), remodStart) {
			skip = true
			continue
		}

		if skip && bytes.Equal(append(line, '\n'), remodEnd) {
			skip = false
			continue
		}

		if skip {
			continue
		}

		must(buf.Write(line))

		if !bytes.Equal(last, []byte("\n")) {
			must(buf.WriteRune('\n'))
		}

		last = line
	}

	return append(bytes.TrimSpace(buf.Bytes()), '\n'), nil
}

func prepareGoDev(godev []byte) []byte {

	buf := bytes.NewBuffer(nil)
	must(buf.Write(remodStart))
	must(buf.WriteRune('\n'))
	must(buf.Write(godev))
	must(buf.WriteRune('\n'))
	must(buf.Write(remodEnd))

	return buf.Bytes()
}

func assemble(gomod, godev []byte) []byte {

	buf := bytes.NewBuffer(gomod)
	must(buf.WriteRune('\n'))
	must(buf.Write(godev))

	return buf.Bytes()
}
