package util

import "testing"

func TestParseBearerToken(t *testing.T) {
	type args struct {
		token string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "normal",
			args: args{
				token: "bearer token",
			},
			want: "token",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseBearerToken(tt.args.token); got != tt.want {
				t.Errorf("ParseBearerToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
