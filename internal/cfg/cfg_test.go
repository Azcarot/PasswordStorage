// Package cfg - обработка флагов и шифрование данных пользователя

package cfg

import (
	"flag"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseFlagsAndENV(t *testing.T) {
	tests := []struct {
		name     string
		want     Flags
		envValue string
		envName  string
	}{
		{name: "NoEnv", want: Flags{FlagAddr: "localhost:8080"}},
		{name: "SomeEnv", want: Flags{FlagAddr: "111"}, envValue: "111", envName: "RUN_ADDRESS"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.envName) != 0 {

				err := os.Setenv(tt.envName, tt.envValue)
				require.NoError(t, err)
			}
			flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
			flagSet.String("myflag", "test_flag_value", "a flag for testing")

			// Parse the flags
			flagSet.Parse([]string{"-myflag=test_flag_value"})
			oldCommandLine := flag.CommandLine
			flag.CommandLine = flagSet
			if got := ParseFlagsAndENV(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseFlagsAndENV() = %v, want %v", got, tt.want)
			}

			flag.CommandLine = oldCommandLine
		})
	}
}
