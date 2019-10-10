package remod

import (
	"reflect"
	"testing"
)

func TestEnable(t *testing.T) {

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
				append(
					gomod,
					[]byte(`
replace ( // remod:replacements
    go.aporeto.io/gaia => ../gaia
    go.aporeto.io/manipulate => ../manipulate
    go.aporeto.io/midgard-lib => ../midgard-lib
)
`,
					)...,
				),
				[]string{
					"go.aporeto.io/gaia",
					"go.aporeto.io/manipulate",
					"go.aporeto.io/midgard-lib",
				},
				"../",
				"",
			},
			append(
				gomod,
				[]byte(`
replace ( // remod:replacements
    go.aporeto.io/gaia => ../gaia
    go.aporeto.io/manipulate => ../manipulate
    go.aporeto.io/midgard-lib => ../midgard-lib
)
`,
				)...,
			),
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
			append(
				gomod,
				[]byte(`
replace ( // remod:replacements
	go.aporeto.io/gaia => ../gaia
	go.aporeto.io/manipulate => ../manipulate
	go.aporeto.io/midgard-lib => ../midgard-lib
)
`,
				)...,
			),
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
			append(
				gomod,
				[]byte(`
replace ( // remod:replacements
	go.aporeto.io/gaia => github.com/la/gaia v12.0.1
	go.aporeto.io/manipulate => github.com/la/manipulate v12.0.1
	go.aporeto.io/midgard-lib => github.com/la/midgard-lib v12.0.1
)
`,
				)...,
			),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Enable(tt.args.data, tt.args.modules, tt.args.base, tt.args.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("Enable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Enable() = >>%v<<, want >>%v<<", string(got), string(tt.want))
			}
		})
	}
}

func TestDisable(t *testing.T) {

	gomod := []byte(`module go.aporeto.io/test

go 1.12

require (
    github.com/aporeto-inc/influxdb1-client v0.0.0-20190909164713-fce670a2a4a6
    go.aporeto.io/gaia v1.94.1-0.20191009190518-2222e09dd2f3
    go.aporeto.io/manipulate v1.114.1-0.20191009190511-3ce5141f45cd
    go.aporeto.io/midgard-lib v1.69.1-0.20191009190649-7e0a1cd52585
)

require cloud.google.com/go/storage v1.1.0 // indirect

replace ( // remod:replacements
	go.aporeto.io/gaia => github.com/la/gaia
	go.aporeto.io/manipulate => github.com/la/manipulate
	go.aporeto.io/midgard-lib => github.com/la/midgard-lib
)
`)

	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			"simple",
			args{
				gomod,
			},
			[]byte(`module go.aporeto.io/test

go 1.12

require (
    github.com/aporeto-inc/influxdb1-client v0.0.0-20190909164713-fce670a2a4a6
    go.aporeto.io/gaia v1.94.1-0.20191009190518-2222e09dd2f3
    go.aporeto.io/manipulate v1.114.1-0.20191009190511-3ce5141f45cd
    go.aporeto.io/midgard-lib v1.69.1-0.20191009190649-7e0a1cd52585
)

require cloud.google.com/go/storage v1.1.0 // indirect
`),
			false,
		},

		{
			"not enabled",
			args{
				[]byte(`module go.aporeto.io/test

go 1.12

require (
    github.com/aporeto-inc/influxdb1-client v0.0.0-20190909164713-fce670a2a4a6
    go.aporeto.io/gaia v1.94.1-0.20191009190518-2222e09dd2f3
    go.aporeto.io/manipulate v1.114.1-0.20191009190511-3ce5141f45cd
    go.aporeto.io/midgard-lib v1.69.1-0.20191009190649-7e0a1cd52585
)

require cloud.google.com/go/storage v1.1.0 // indirect
`),
			},
			[]byte(`module go.aporeto.io/test

go 1.12

require (
    github.com/aporeto-inc/influxdb1-client v0.0.0-20190909164713-fce670a2a4a6
    go.aporeto.io/gaia v1.94.1-0.20191009190518-2222e09dd2f3
    go.aporeto.io/manipulate v1.114.1-0.20191009190511-3ce5141f45cd
    go.aporeto.io/midgard-lib v1.69.1-0.20191009190649-7e0a1cd52585
)

require cloud.google.com/go/storage v1.1.0 // indirect
`),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Disable(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Disable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Disable() = %v, want %v", got, tt.want)
			}
		})
	}
}
