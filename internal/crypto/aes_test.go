package crypto

import (
	"reflect"
	"testing"
)

func TestGenerateRandomKey(t *testing.T) {
	type tCase struct {
		name    string
		want    []byte
		wantErr bool
	}

	tests := []tCase{
		{name: "", want: []byte(""), wantErr: false},
		{name: "", want: []byte("这里写期望的结果"), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateRandomKey()
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateRandomKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateRandomKey() got = %v, want %v", got, tt.want)
			}
		})
	}
}
