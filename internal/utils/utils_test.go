// Package utils - обработка флагов и шифрование данных пользователя

package utils

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

func TestShaData(t *testing.T) {
	type args struct {
		result string
		key    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "BasicTest", args: args{result: "1", key: "1"}, want: "b_bO7bHBACZoSdfk4TRJyY0cmOVpmU-N2_57UueA3q0="},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ShaData(tt.args.result, tt.args.key); got != tt.want {
				t.Errorf("ShaData() = %v, want %v", got, tt.want)
			}
		})
	}
}
