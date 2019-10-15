package remod

import "fmt"

const (
	goDev = "remod.dev"
)

var (
	remodStart = []byte("// remod:start\n")
	remodEnd   = []byte("// remod:end\n")
)

func goModBackup() string {

	br, err := branchName()
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf(".remod/%s.mod", br)
}

func goSumBackup() string {

	br, err := branchName()
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf(".remod/%s.sum", br)
}
