package shortenalgorithm

import "testing"

func TestGetShortName(t *testing.T) {
	tests := []struct {
		name   string
		lastID int
		want   string
	}{
		{
			name:   "test1",
			lastID: 1,
			want:   "zE",
		},
		{
			name:   "test1",
			lastID: 2,
			want:   "Xz",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := GetShortName(tt.lastID)
			if got != tt.want {
				t.Errorf("GetShortName() got = %v, want %v", got, tt.want)
			}
		})
	}
}
