package xor

import (
	"reflect"
	"testing"
)

func TestEncryptXor(t *testing.T) {
	type args struct {
		src []byte
		key []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "1",
			args: args{[]byte("hello"), []byte("key")},
			want: DecryptXor([]byte("hello"), []byte("key")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EncryptXor(tt.args.src, tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EncryptXor() = %v, want %v", got, tt.want)
			}
		})
	}
}
