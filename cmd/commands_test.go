package cmd

import (
	"testing"

	"github.com/jessevdk/go-flags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppendDBTypes(t *testing.T) {
	tests := []struct {
		name         string
		data         interface{}
		expectOption bool
	}{
		{
			name: "no matching options",
			data: &struct {
				Command1 *struct {
					Foo string `long:"foo"`
					Bar string `long:"bar"`
				} `command:"test"`
			}{},
		},
		{
			name: "db-type option",
			data: &struct {
				Command1 *struct {
					Foo string `long:"db-type"`
					Bar string `long:"bar"`
				} `command:"test"`
			}{},
			expectOption: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := flags.NewParser(tt.data, flags.Default)
			AppendDBTypes(p)
			require.Len(t, p.Commands(), 1)

			opt := p.Commands()[0].FindOptionByLongName("db-type")
			if tt.expectOption {
				require.NotNil(t, opt)
				assert.ElementsMatch(t, []string{"mysql", "sqlite3"}, opt.Choices)
			} else {
				assert.Nil(t, opt)
			}
		})
	}
}
