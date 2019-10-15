package remod

import "fmt"

const (
	goDev     = "remod.dev"
	modBackup = ".remod/mod"
	sumBackup = ".remod/sum"
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

	return fmt.Sprintf("%s.%s", modBackup, br)
}

func goSumBackup() string {

	br, err := branchName()
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%s.%s", sumBackup, br)
}
