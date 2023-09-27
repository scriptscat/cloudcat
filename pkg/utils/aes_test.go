package utils

import (
	"bytes"
	"github.com/codfrm/cago/pkg/utils"
	"io"
	"testing"
)

func TestAes(t *testing.T) {
	type args struct {
		key  []byte
		data string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"case1", args{[]byte("1234567890123456"), "你好"}, false},
		{"case2", args{[]byte("1234567890123456"), "0000000000000000"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.NewBufferString(tt.args.data)
			got, err := NewAesEncrypt(tt.args.key, buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAesEncrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			encryptData := bytes.NewBuffer(nil)
			_, err = io.Copy(encryptData, got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Copy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			decrypt, err := NewAesDecrypt(tt.args.key, encryptData)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAesDecrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			_, err = io.Copy(buf, decrypt)
			if (err != nil) != tt.wantErr {
				t.Errorf("Copy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if buf.String() != tt.args.data {
				t.Errorf("Decrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func FuzzAes(f *testing.F) {
	key := []byte(utils.RandString(16, utils.Mix))
	testcases := [][]byte{
		[]byte("Hello, world"),
		[]byte(" "),
		[]byte("!12345"),
		[]byte("0000000000000000"),
		[]byte(""),
	}
	for _, tc := range testcases {
		f.Add(tc)
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		buf := bytes.NewBuffer(data)
		got, err := NewAesEncrypt(key, buf)
		if err != nil {
			t.Errorf("NewAesEncrypt() error = %v", err)
			return
		}
		encryptData := bytes.NewBuffer(nil)
		_, err = io.Copy(encryptData, got)
		if err != nil {
			t.Errorf("Copy() error = %v", err)
			return
		}
		decrypt, err := NewAesDecrypt(key, encryptData)
		if err != nil {
			t.Errorf("NewAesDecrypt() error = %v", err)
			return
		}
		_, err = io.Copy(buf, decrypt)
		if err != nil {
			t.Errorf("Copy() error = %v", err)
			return
		}
		if !bytes.Equal(buf.Bytes(), data) {
			t.Errorf("Decrypt() error = %v", err)
			return
		}
	})
}
