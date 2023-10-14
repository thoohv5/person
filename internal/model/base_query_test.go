package model

import "testing"

func Test_escapeLikePattern(t *testing.T) {
	type args struct {
		pattern string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"%", args{"%"}, `\%`},
		{"%1", args{"%12"}, `\%12`},
		{"%_", args{"%12_"}, `\%12\_`},
		{"%_\\", args{"%12_\\"}, `\%12\_\\`},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := escapeLikePattern(tt.args.pattern); got != tt.want {
				t.Errorf("escapeLikePattern() = %v, want %v", got, tt.want)
			}
		})
	}
}
