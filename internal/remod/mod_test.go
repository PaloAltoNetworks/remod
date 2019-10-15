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

func TestMakeGoModDev(t *testing.T) {

	gomod := []byte(`module go.aporeto.io/test

go 1.12

require (
    github.com/aporeto-inc/influxdb1-client v0.0.0-20190909164713-fce670a2a4a6
    go.aporeto.io/gaia v1.94.1-0.20191009190518-2222e09dd2f3
    go.aporeto.io/manipulate v1.114.1-0.20191009190511-3ce5141f45cd
    go.aporeto.io/midgard-lib v1.69.1-0.20191009190649-7e0a1cd52585
)

require cloud.google.com/go/storage v1.1.0 // indirect
`)
	type args struct {
		data    []byte
		modules []string
		base    string
		version string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"no match",
			args{
				gomod,
				[]string{},
				"../",
				"",
			},
			gomod,
			false,
		},

		{
			"already enabled",
			args{
				[]byte(`replace (
	go.aporeto.io/gaia => ../gaia
	go.aporeto.io/manipulate => ../manipulate
	go.aporeto.io/midgard-lib => ../midgard-lib
)
`),
				[]string{
					"go.aporeto.io/gaia",
					"go.aporeto.io/manipulate",
					"go.aporeto.io/midgard-lib",
				},
				"../",
				"",
			},
			[]byte(`replace (
	go.aporeto.io/gaia => ../gaia
	go.aporeto.io/manipulate => ../manipulate
	go.aporeto.io/midgard-lib => ../midgard-lib
)
`),
			false,
		},

		{
			"simple matching go.aporeto.io on ../",
			args{
				gomod,
				[]string{
					"go.aporeto.io/gaia",
					"go.aporeto.io/manipulate",
					"go.aporeto.io/midgard-lib",
				},
				"../",
				"",
			},
			[]byte(`replace (
	go.aporeto.io/gaia => ../gaia
	go.aporeto.io/manipulate => ../manipulate
	go.aporeto.io/midgard-lib => ../midgard-lib
)
`),
			false,
		},

		{
			"simple matching go.aporeto.io on github.com/la",
			args{
				gomod,
				[]string{
					"go.aporeto.io/gaia",
					"go.aporeto.io/manipulate",
					"go.aporeto.io/midgard-lib",
				},
				"github.com/la/",
				"v12.0.1",
			},
			[]byte(`replace (
	go.aporeto.io/gaia => github.com/la/gaia v12.0.1
	go.aporeto.io/manipulate => github.com/la/manipulate v12.0.1
	go.aporeto.io/midgard-lib => github.com/la/midgard-lib v12.0.1
)
`),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := makeGoModDev(tt.args.data, tt.args.modules, tt.args.base, tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("makeGoModDev() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("makeGoModDev() = >>%v<<, want >>%v<<", string(got), string(tt.want))
			}
		})
	}
}
