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
