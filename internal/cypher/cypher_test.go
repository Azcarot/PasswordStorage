// package cypher - функции (де)/шифрования

package cypher

import "testing"

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
