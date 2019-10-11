// Copyright 2019 Aporeto Inc.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package remod

import (
	"bufio"
	"bytes"
	"sort"
)

// Extract extracts the required modules with the given prefixes
func Extract(data []byte, prefixes []string, excluded []string) ([]string, error) {

	incpfx := make([][]byte, len(prefixes))
	for i, m := range prefixes {
		incpfx[i] = []byte(m)
	}

	expfx := make([][]byte, len(excluded))
	for i, m := range excluded {
		expfx[i] = []byte(m)
	}

	scanner := bufio.NewScanner(bytes.NewBuffer(data))
	singleRequireStartPrefix := []byte("require ")
	multiRequireStartPrefix := []byte("require (")
	multiRequireEndPrefix := []byte(")")

	modules := map[string]struct{}{}

	match := func(mod []byte) bool {

		for _, prefix := range expfx {
			if bytes.HasPrefix(mod, prefix) {
				return false
			}
		}

		for _, prefix := range incpfx {
			if bytes.HasPrefix(mod, prefix) {
				return true
			}
		}

		return false
	}

	var multiStart bool
	for scanner.Scan() {

		line := scanner.Bytes()

		if bytes.HasPrefix(line, multiRequireStartPrefix) {
			multiStart = true
			continue
		} else if bytes.HasPrefix(line, singleRequireStartPrefix) {
			mod := bytes.TrimSpace(line)
			mod = bytes.Replace(mod, singleRequireStartPrefix, nil, -1)

			if match(mod) {
				modules[string(bytes.SplitN(mod, []byte(" "), 2)[0])] = struct{}{}
			}

			continue
		}

		if multiStart && bytes.HasPrefix(line, multiRequireEndPrefix) {
			multiStart = false
			continue
		}

		if multiStart {
			mod := bytes.TrimSpace(line)

			if match(mod) {
				modules[string(bytes.SplitN(mod, []byte(" "), 2)[0])] = struct{}{}
			}
		}
	}

	out := make([]string, len(modules))
	var i int
	for mod := range modules {
		out[i] = mod
		i++
	}

	sort.Strings(out)

	return out, nil
}
