package crypto

import (
	"bytes"
	"crypto/rand"
	"errors"
)

// PaddingMode 填充模式
type PaddingMode int

const (
	PKCS7     PaddingMode = iota // PKCS7填充
	Zero                          // Zero填充
	ISO10126                      // ISO10126填充
	ANSIX923                      // ANSI X9.23填充
	NoPadding                     // 不填充
)

// Padding 添加填充
func Padding(src []byte, blockSize int, mode PaddingMode) ([]byte, error) {
	switch mode {
	case PKCS7:
		return pkcs7Padding(src, blockSize), nil
	case Zero:
		return zeroPadding(src, blockSize), nil
	case ISO10126:
		return iso10126Padding(src, blockSize)
	case ANSIX923:
		return ansix923Padding(src, blockSize), nil
	case NoPadding:
		if len(src)%blockSize != 0 {
			return nil, errors.New("数据长度必须是块大小的整数倍")
		}
		return src, nil
	default:
		return nil, errors.New("不支持的填充模式")
	}
}

// UnPadding 移除填充
func UnPadding(src []byte, mode PaddingMode) ([]byte, error) {
	switch mode {
	case PKCS7:
		return pkcs7UnPadding(src)
	case Zero:
		return zeroUnPadding(src), nil
	case ISO10126:
		return iso10126UnPadding(src)
	case ANSIX923:
		return ansix923UnPadding(src)
	case NoPadding:
		return src, nil
	default:
		return nil, errors.New("不支持的填充模式")
	}
}

// PKCS7填充
func pkcs7Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

// PKCS7去填充
func pkcs7UnPadding(src []byte) ([]byte, error) {
	length := len(src)
	if length == 0 {
		return nil, errors.New("数据为空")
	}
	unpadding := int(src[length-1])
	if unpadding > length {
		return nil, errors.New("无效的填充")
	}
	return src[:(length - unpadding)], nil
}

// Zero填充
func zeroPadding(src []byte, blockSize int) []byte {
	paddingCount := blockSize - len(src)%blockSize
	if paddingCount == blockSize {
		return src
	}
	return append(src, bytes.Repeat([]byte{0}, paddingCount)...)
}

// Zero去填充
func zeroUnPadding(src []byte) []byte {
	return bytes.TrimRight(src, string([]byte{0}))
}

// ISO10126填充（随机填充+长度）
func iso10126Padding(src []byte, blockSize int) ([]byte, error) {
	padding := blockSize - len(src)%blockSize
	if padding == 0 {
		padding = blockSize
	}
	padtext := make([]byte, padding)
	if padding > 1 {
		_, err := rand.Read(padtext[:padding-1])
		if err != nil {
			return nil, err
		}
	}
	padtext[padding-1] = byte(padding)
	return append(src, padtext...), nil
}

// ISO10126去填充
func iso10126UnPadding(src []byte) ([]byte, error) {
	length := len(src)
	if length == 0 {
		return nil, errors.New("数据为空")
	}
	unpadding := int(src[length-1])
	if unpadding > length || unpadding == 0 {
		return nil, errors.New("无效的填充")
	}
	return src[:(length - unpadding)], nil
}

// ANSI X9.23填充（零填充+长度）
func ansix923Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	if padding == 0 {
		padding = blockSize
	}
	padtext := make([]byte, padding)
	padtext[padding-1] = byte(padding)
	return append(src, padtext...)
}

// ANSI X9.23去填充
func ansix923UnPadding(src []byte) ([]byte, error) {
	length := len(src)
	if length == 0 {
		return nil, errors.New("数据为空")
	}
	unpadding := int(src[length-1])
	if unpadding > length || unpadding == 0 {
		return nil, errors.New("无效的填充")
	}
	return src[:(length - unpadding)], nil
}
