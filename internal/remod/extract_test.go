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
	"reflect"
	"testing"
)

func TestExtract(t *testing.T) {
	type args struct {
		data     []byte
		prefixes []string
		excluded []string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			"empty",
			args{
				[]byte{},
				nil,
				nil,
			},
			[]string{},
			false,
		},

		{
			"simple single not matching",
			args{
				[]byte(`require a/b v1.0.0`),
				[]string{"z"},
				[]string{},
			},
			[]string{},
			false,
		},
		{
			"simple single full matching",
			args{
				[]byte(`require a/b v1.0.0`),
				[]string{"a/b"},
				[]string{},
			},
			[]string{"a/b"},
			false,
		},
		{
			"simple single partial matching",
			args{
				[]byte(`require a/b v1.0.0`),
				[]string{"a/"},
				[]string{},
			},
			[]string{"a/b"},
			false,
		},

		{
			"multiple single not matching",
			args{
				[]byte(`
require a/b/1 v1.0.0
require a/b/2 v1.0.0
require c/a/1 v1.0.0
`),
				[]string{"z"},
				[]string{},
			},
			[]string{},
			false,
		},
		{
			"multiple single matching",
			args{
				[]byte(`
require a/b/1 v1.0.0
require a/b/2 v1.0.0
require c/a/1 v1.0.0
`),
				[]string{"a/b"},
				[]string{},
			},
			[]string{"a/b/1", "a/b/2"},
			false,
		},

		{
			"multiple multiple not matching",
			args{
				[]byte(`
require (
    a/b/1 v1.0.0
    a/b/2 v1.0.0
    c/a/1 v1.0.0
)

require a/b/3 v1.0.0
`),
				[]string{"z"},
				[]string{},
			},
			[]string{},
			false,
		},
		{
			"multiple multiple matching",
			args{
				[]byte(`
require (
    a/b/1 v1.0.0
    a/b/2 v1.0.0
    c/a/1 v1.0.0
)

require a/b/3 v1.0.0
`),
				[]string{"a/b"},
				[]string{"a/b/3"},
			},
			[]string{"a/b/1", "a/b/2"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Extract(tt.args.data, tt.args.prefixes, tt.args.excluded)
			if (err != nil) != tt.wantErr {
				t.Errorf("Extract() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Extract() = %v, want %v", got, tt.want)
			}
		})
	}
}
