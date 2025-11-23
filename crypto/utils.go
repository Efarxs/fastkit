package crypto

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
)

// OutputFormat 输出格式
type OutputFormat int

const (
	OutputBase64 OutputFormat = iota // Base64输出
	OutputString                     // 字符串输出
	OutputHex                        // 十六进制输出
)

// EncryptMode 加密模式
type EncryptMode string

const (
	ECB EncryptMode = "ECB" // 电子密码本模式
	CBC EncryptMode = "CBC" // 密码块链接模式
	CTR EncryptMode = "CTR" // 计数器模式
	CFB EncryptMode = "CFB" // 密码反馈模式
	OFB EncryptMode = "OFB" // 输出反馈模式
)

// FormatOutput 格式化输出
func FormatOutput(data []byte, format OutputFormat) string {
	switch format {
	case OutputBase64:
		return base64.StdEncoding.EncodeToString(data)
	case OutputHex:
		return hex.EncodeToString(data)
	case OutputString:
		return string(data)
	default:
		return string(data)
	}
}

// ParseInput 解析输入（支持base64、hex、原始字符串）
func ParseInput(data string, format OutputFormat) ([]byte, error) {
	switch format {
	case OutputBase64:
		return base64.StdEncoding.DecodeString(data)
	case OutputHex:
		return hex.DecodeString(data)
	case OutputString:
		return []byte(data), nil
	default:
		return []byte(data), nil
	}
}

// PrepareKey 准备密钥（支持最大100字符，根据算法要求调整）
func PrepareKey(key string, requiredLength int) []byte {
	keyBytes := []byte(key)
	
	// 如果密钥长度正好，直接返回
	if len(keyBytes) == requiredLength {
		return keyBytes
	}
	
	// 如果密钥过长，截取
	if len(keyBytes) > requiredLength {
		return keyBytes[:requiredLength]
	}
	
	// 如果密钥过短，使用MD5扩展或重复填充
	if requiredLength <= 16 {
		// 对于较短的密钥长度（如DES的8字节），使用MD5后截取
		hash := md5.Sum(keyBytes)
		return hash[:requiredLength]
	}
	
	// 对于较长的密钥，使用重复填充
	result := make([]byte, requiredLength)
	for i := 0; i < requiredLength; i++ {
		result[i] = keyBytes[i%len(keyBytes)]
	}
	return result
}

// PrepareIV 准备初始化向量
func PrepareIV(iv string, requiredLength int) ([]byte, error) {
	if iv == "" {
		return nil, errors.New("IV不能为空")
	}
	
	ivBytes := []byte(iv)
	
	// 如果IV长度正好，直接返回
	if len(ivBytes) == requiredLength {
		return ivBytes, nil
	}
	
	// 如果IV过长，截取
	if len(ivBytes) > requiredLength {
		return ivBytes[:requiredLength], nil
	}
	
	// 如果IV过短，使用MD5扩展或重复填充
	if requiredLength <= 16 {
		hash := md5.Sum(ivBytes)
		return hash[:requiredLength], nil
	}
	
	// 对于较长的IV，使用重复填充
	result := make([]byte, requiredLength)
	for i := 0; i < requiredLength; i++ {
		result[i] = ivBytes[i%len(ivBytes)]
	}
	return result, nil
}
