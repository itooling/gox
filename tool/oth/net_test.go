package oth

import "testing"

func TestGetIp(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ip, err := GetIP()
			if err == nil {
				t.Log(ip)
			}
		})
	}
}
