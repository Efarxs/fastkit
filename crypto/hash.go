package crypto

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"

	"golang.org/x/crypto/md4"
	"golang.org/x/crypto/sha3"
)

// HashType 哈希算法类型
type HashType string

const (
	MD5Hash       HashType = "MD5"
	MD4Hash       HashType = "MD4"
	SHA1Hash      HashType = "SHA1"
	SHA224Hash    HashType = "SHA224"
	SHA256Hash    HashType = "SHA256"
	SHA384Hash    HashType = "SHA384"
	SHA512Hash    HashType = "SHA512"
	SHA3_224Hash  HashType = "SHA3-224"
	SHA3_256Hash  HashType = "SHA3-256"
	SHA3_384Hash  HashType = "SHA3-384"
	SHA3_512Hash  HashType = "SHA3-512"
)

// Hash 通用哈希函数
func Hash(data []byte, hashType HashType) []byte {
	var h hash.Hash
	
	switch hashType {
	case MD5Hash:
		h = md5.New()
	case MD4Hash:
		h = md4.New()
	case SHA1Hash:
		h = sha1.New()
	case SHA224Hash:
		h = sha256.New224()
	case SHA256Hash:
		h = sha256.New()
	case SHA384Hash:
		h = sha512.New384()
	case SHA512Hash:
		h = sha512.New()
	case SHA3_224Hash:
		h = sha3.New224()
	case SHA3_256Hash:
		h = sha3.New256()
	case SHA3_384Hash:
		h = sha3.New384()
	case SHA3_512Hash:
		h = sha3.New512()
	default:
		h = md5.New()
	}
	
	h.Write(data)
	return h.Sum(nil)
}

// HashString 哈希字符串并返回字符串（默认hex格式）
func HashString(data string, hashType HashType) string {
	result := Hash([]byte(data), hashType)
	return hex.EncodeToString(result)
}

// HashToFormat 哈希并按指定格式输出
func HashToFormat(data string, hashType HashType, format OutputFormat) string {
	result := Hash([]byte(data), hashType)
	return FormatOutput(result, format)
}

// MD5 计算MD5哈希
func MD5(data string) string {
	return HashString(data, MD5Hash)
}

// MD4 计算MD4哈希
func MD4(data string) string {
	return HashString(data, MD4Hash)
}

// SHA1 计算SHA1哈希
func SHA1(data string) string {
	return HashString(data, SHA1Hash)
}

// SHA224 计算SHA224哈希
func SHA224(data string) string {
	return HashString(data, SHA224Hash)
}

// SHA256 计算SHA256哈希
func SHA256(data string) string {
	return HashString(data, SHA256Hash)
}

// SHA384 计算SHA384哈希
func SHA384(data string) string {
	return HashString(data, SHA384Hash)
}

// SHA512 计算SHA512哈希
func SHA512(data string) string {
	return HashString(data, SHA512Hash)
}

// SHA3_224 计算SHA3-224哈希
func SHA3_224(data string) string {
	return HashString(data, SHA3_224Hash)
}

// SHA3_256 计算SHA3-256哈希
func SHA3_256(data string) string {
	return HashString(data, SHA3_256Hash)
}

// SHA3_384 计算SHA3-384哈希
func SHA3_384(data string) string {
	return HashString(data, SHA3_384Hash)
}

// SHA3_512 计算SHA3-512哈希
func SHA3_512(data string) string {
	return HashString(data, SHA3_512Hash)
}
