package logger

import "testing"

func Test_argsToString(t *testing.T) {
	type args struct {
		args []interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "should pass",
			args: args{
				args: []interface{}{
					"seperate", "strings",
				},
			},
			want: "seperate strings",
		},
		{
			name: "should pass",
			args: args{
				args: []interface{}{
					"seperate", "strings", "plus one",
				},
			},
			want: "seperate strings plus one",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := argsToString(tt.args.args); got != tt.want {
				t.Errorf("argsToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
