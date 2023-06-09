package getfreeport

import (
	"net"
	"testing"
)

func TestGetFreePort(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "ok",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFreePort()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFreePort() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if _, err = net.Listen("tcp", "localhost:"+string(got)); err != nil {
				t.Errorf("GetFreePort() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
