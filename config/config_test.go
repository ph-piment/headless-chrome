package config

import (
	"testing"
)

func TestNewConfig(t *testing.T) {
	/*
		cases := []struct {
			name     string
			appMode  string
			expected config.Config
		}{
			{
				name:    "app",
				appMode: "app",
				expected: config.Config{
					PATH: config.PathConfig{
						OutputCompareDir: "/go/src/work/outputs/images/compare/",
					},
				},
			},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				res, err := config.NewConfig(c.appMode)

				assert.Equal(t, nil, err)
				assert.Equal(t, c.expected.PATH.OutputCompareDir, res.PATH.OutputCompareDir)
			})
		}
	*/
}
