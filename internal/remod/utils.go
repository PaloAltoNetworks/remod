package remod

import (
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
