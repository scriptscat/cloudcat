package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"io"
)

// AesEncrypt aes-cbc pkcs7padding
type AesEncrypt struct {
	buf  []byte
	mode cipher.BlockMode
	body io.Reader
}

func NewAesEncrypt(key []byte, body io.Reader) (*AesEncrypt, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	mode := cipher.NewCBCEncrypter(block, key[:aes.BlockSize])
	return &AesEncrypt{
		mode: mode,
		body: body,
	}, nil
}

func (a *AesEncrypt) Read(p []byte) (int, error) {
	buf := make([]byte, len(p))
	n, err := a.body.Read(buf)
	if err != nil {
		if err == io.EOF {
			if len(a.buf) == 0 && n == 0 {
				if a.buf != nil {
					// 填充16个16
					buf = bytes.Repeat([]byte{16}, 16)
					if len(p) < len(buf) {
						return 0, io.ErrShortBuffer
					}
					a.mode.CryptBlocks(p, buf)
					a.buf = nil
					return len(buf), nil
				}
				return 0, io.EOF
			}
			// 最后一次加密, 不足blockSize的整数倍的, 补齐
			buf = append(a.buf, buf[:n]...)
			// 需要补齐的长度
			padding := aes.BlockSize - len(buf)%aes.BlockSize
			// 补齐
			buf = append(buf, bytes.Repeat([]byte{byte(padding)}, padding)...)
			if len(p) < len(buf) {
				return 0, io.ErrShortBuffer
			}
			a.mode.CryptBlocks(p, buf)
			a.buf = nil
			return len(buf), nil
		}
		return n, err
	}
	// 只加密blockSize的整数倍
	buf = append(a.buf, buf[:n]...)
	// 有效长度
	n = len(buf) - len(buf)%aes.BlockSize
	a.buf = buf[n:]
	a.mode.CryptBlocks(p[:n], buf[:n])
	return n, nil
}

type AesDecrypt struct {
	buf  []byte
	mode cipher.BlockMode
	body io.Reader
}

func NewAesDecrypt(key []byte, body io.Reader) (*AesDecrypt, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	mode := cipher.NewCBCDecrypter(block, key[:aes.BlockSize])
	return &AesDecrypt{
		mode: mode,
		body: body,
	}, nil
}

func (a *AesDecrypt) Read(p []byte) (int, error) {
	buf := make([]byte, len(p))
	n, err := a.body.Read(buf)
	if err != nil {
		if err == io.EOF {
			if len(a.buf) == 0 && n == 0 {
				return 0, io.EOF
			}
			buf = append(a.buf, buf[:n]...)
			a.mode.CryptBlocks(buf, buf)
			// 去掉填充的数据
			padding := int(buf[len(buf)-1])
			if padding > 0 && padding <= aes.BlockSize {
				n = len(buf) - padding
			} else {
				n = len(buf)
			}
			copy(p, buf[:n])
			a.buf = nil
			return n, nil
		}
		return n, err
	}
	// 只加密blockSize的整数倍
	buf = append(a.buf, buf[:n]...)
	// 有效长度
	n = len(buf) - len(buf)%aes.BlockSize
	// 再留下最后blockSize的长度
	n = n - aes.BlockSize
	if n < 0 {
		a.buf = buf
		return 0, err
	}
	a.buf = buf[n:]
	a.mode.CryptBlocks(p[:n], buf[:n])
	return n, nil
}

type readClose struct {
	io.Reader
	io.Closer
}

func WarpCloser(r io.Reader, c io.Closer) io.ReadCloser {
	return &readClose{
		Reader: r,
		Closer: c,
	}
}
