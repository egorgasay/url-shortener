package rest

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_checkCookies(t *testing.T) {
	type args struct {
		cookie string
		key    []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "bad #1",
			args: args{
				cookie: "test",
				key:    []byte("test"),
			},
			want: false,
		},
		{
			name: "bad #2",
			args: args{
				cookie: "6b8be1b76bbf7e9f1545a76dcb8631a6-0067cec39dbfdfb7d3c3f37b7ddb2a06",
				key:    []byte("key_example"),
			},

			want: false,
		},
		{
			name: "ok",
			args: args{
				cookie: "2daa0f44d32c33a74cfbfd96fd58134649862dd008bd8cba3c331314e81fb551-31363832313834363939393633363234313030",
				key:    []byte("CHANGE ME"),
			},

			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, checkCookies(tt.args.cookie, tt.args.key), "checkCookies(%v, %v)", tt.args.cookie, tt.args.key)
		})
	}
}
