package remod

import (
	"bufio"
	"bytes"
	"fmt"
	"path/filepath"
)

// makeGoModDev builds the dev mod file.
func makeGoModDev(modules []string, base string, version string) []byte {

	if version != "" {
		version = " " + version
	}

	buf := bytes.NewBuffer(nil)

	switch len(modules) {
	case 0:
		must(buf.WriteString("// insert your development replacements in the remod.dev file"))
	case 1:
		must(buf.WriteString(fmt.Sprintf("replace %s => %s%s%s\n", modules[0], base, filepath.Base(modules[0]), version)))
	default:
		must(buf.WriteString("replace (\n"))
		for _, m := range modules {
			must(buf.WriteString(fmt.Sprintf("\t%s => %s%s%s\n", m, base, filepath.Base(m), version)))
		}
		must(buf.WriteString(")\n"))
	}

	return append(bytes.TrimSpace(buf.Bytes()), '\n')
}

func must(_ int, err error) {
	if err != nil {
		panic(err)
	}
}

func strip(in []byte) []byte {

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

	return append(bytes.TrimSpace(buf.Bytes()), '\n')
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
