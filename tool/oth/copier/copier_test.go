package copier

import "testing"

func TestCopy(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "one",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Copy()
		})
	}
}
