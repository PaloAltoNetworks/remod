package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {

	cobra.OnInitialize(func() {
		viper.SetEnvPrefix("remod")
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	})

	var rootCmd = &cobra.Command{
		Use: "remod",
	}

	var cmdUpdate = &cobra.Command{
		Use:     "update",
		Aliases: []string{"up"},
		Short:   "Update the modules in the given path",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			folder := viper.GetString("folder")
			recursive := viper.GetBool("recursive")

			var paths []string
			if recursive {

				items, err := ioutil.ReadDir(folder)
				if err != nil {
					return fmt.Errorf("unable to list content of dir: %s", err)
				}

				for _, item := range items {
					if !item.IsDir() {
						continue
					}

					p := path.Join(folder, item.Name(), "go.mod")
					_, err := os.Stat(p)
					if err != nil {
						if os.IsNotExist(err) {
							continue
						}
						return fmt.Errorf("unable stat path '%s': %s", p, err)
					}

					paths = append(paths, p)
				}
			} else {
				paths = append(paths, path.Join(folder, "go.mod"))
			}

			modPrefixes := viper.GetStringSlice("module")
			prefixes := make([][]byte, len(modPrefixes))
			for i, m := range modPrefixes {
				prefixes[i] = []byte(m)
			}

			for _, p := range paths {

				basedir := filepath.Dir(p)
				if err := os.Chdir(basedir); err != nil {
					return fmt.Errorf("unable to move to %s: %s", basedir, err)
				}

				fmt.Println("* Entering", basedir)

				file, err := os.Open(p)
				if err != nil {
					return fmt.Errorf("unable to open go.mod: %s", err)
				}
				defer file.Close()

				modules, err := extract(file, prefixes)
				if err != nil {
					return fmt.Errorf("unable to extract modules: %s", err)
				}

				if err := update(modules, viper.GetString("version")); err != nil {
					return fmt.Errorf("unable to extract modules: %s", err)
				}
			}

			return nil
		},
	}
	cmdUpdate.Flags().StringP("folder", "f", "./", "Set the path to the folder file")
	cmdUpdate.Flags().BoolP("recursive", "r", false, "If true, remod will look for mod files in given --folder and all 1 level subfolders")
	cmdUpdate.Flags().StringSliceP("module", "m", nil, "Set the prefix of the modules you are interested in")
	cmdUpdate.Flags().String("version", "latest", "Set to which version you want to update the matching modules")

	var cmdDevon = &cobra.Command{
		Use:     "devon",
		Aliases: []string{"on"},
		Short:   "Apply developpment replace directive",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			modPrefixes := viper.GetStringSlice("module")
			prefixes := make([][]byte, len(modPrefixes))
			for i, m := range modPrefixes {
				prefixes[i] = []byte(m)
			}

			file, err := os.Open("go.mod")
			if err != nil {
				return fmt.Errorf("unable to open go.mod: %s", err)
			}
			defer file.Close()

			modules, err := extract(file, prefixes)
			if err != nil {
				return fmt.Errorf("unable to extract modules: %s", err)
			}

			data, err := enable("go.mod", modules, viper.GetString("local"))
			if err != nil {
				return fmt.Errorf("unable to apply dev replacements: %s", err)
			}
			if data == nil {
				return nil
			}

			return ioutil.WriteFile("go.mod", data, 0655)
		},
	}
	cmdDevon.Flags().StringSliceP("module", "m", nil, "Set the prefix of the modules you are interested in")
	cmdDevon.Flags().StringP("local", "l", "../", "Where the replacements are")
	cmdDevon.Flags().String("version", "latest", "Set to which version you want to update the matching modules")

	var cmdDevoff = &cobra.Command{
		Use:     "devoff",
		Aliases: []string{"off"},
		Short:   "Remove developpment replace directive",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			data, err := disable("go.mod")
			if err != nil {
				return fmt.Errorf("unable to remove dev replacements: %s", err)
			}
			if data == nil {
				return nil
			}

			return ioutil.WriteFile("go.mod", data, 0655)
		},
	}

	rootCmd.AddCommand(
		cmdUpdate,
		cmdDevon,
		cmdDevoff,
	)

	rootCmd.Execute() // nolint: errcheck
}

func extract(file io.Reader, prefixes [][]byte) ([]string, error) {

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

func update(modules []string, version string) error {

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

func enable(path string, modules []string, base string) ([]byte, error) {

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

func disable(path string) ([]byte, error) {

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
