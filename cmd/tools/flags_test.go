package tools_test

import (
	"dall06/go-cleanapi/cmd/tools"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlags(t *testing.T) {
	// define test cases
	testCases := []struct {
		name         string
		inputArgs    []string
		expected     tools.FlagValues
		defaulValues bool
	}{
		{
			name:      "it should get default flags",
			inputArgs: []string{},
			expected: tools.FlagValues{
				Port:    "8080",
				Version: "0.0.0",
			},
			defaulValues: true,
		},
		{
			name:      "it should get custom flags",
			inputArgs: []string{"-p", "9000", "-v", "0.0.1"},
			expected: tools.FlagValues{
				Port:    "9000",
				Version: "0.0.1",
			},
			defaulValues: false,
		},
	}
	i := 0
	// run tests
	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// save original flag values
			originalArgs := os.Args
			args := tc.inputArgs

			os.Args = append([]string{originalArgs[0]}, args...)

			flags := tools.NewFlags().Flags()

			// Check the results
			assert.Equal(t, tc.expected.Port, flags.Port)
			assert.Equal(t, tc.expected.Version, flags.Version)

			// Clean up
			os.Args = originalArgs
		})

		i++
		fmt.Println(i)
	}
}
