package config

import "testing"

func TestNew(t *testing.T) {
	c := New()
	if c.Host != defaultHost || c.BaseURL != defaultURL || c.HTTPS || c.DBConfig.DataSourcePath != defaultPath {
		t.Errorf("New() error = %v, %v, %v, %v", c.Host, c.BaseURL, c.HTTPS, c.DBConfig.DataSourcePath)
	}
}

func TestModify(t *testing.T) {
	type fl struct {
		Host  string
		HTTPS bool
	}

	tests := []struct {
		name    string
		file    string
		wantErr bool
		f       fl
		wantF   bool
	}{
		{
			name:    "file not exist",
			file:    "test.txt3",
			wantErr: true,
		},
		{
			name:    "ok",
			file:    "config-test.json",
			wantErr: false,
			wantF:   true,
			f: fl{
				Host:  "localhost:8090",
				HTTPS: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Modify(tt.file); (err != nil) != tt.wantErr {
				t.Errorf("Modify() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantF && (*f.Host != tt.f.Host || *f.HTTPS != tt.f.HTTPS) {
				t.Errorf("Modify() error = %v, %v", *f.Host != tt.f.Host, *f.HTTPS != tt.f.HTTPS)
			}
		})
	}
}
