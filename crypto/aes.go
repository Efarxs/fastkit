package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

// AESKeySize AES密钥大小
type AESKeySize int

const (
	AES128 AESKeySize = 16 // 128位
	AES192 AESKeySize = 24 // 192位
	AES256 AESKeySize = 32 // 256位
)

// AESCrypto AES加密器
type AESCrypto struct {
	key     []byte
	iv      []byte
	mode    EncryptMode
	padding PaddingMode
	keySize AESKeySize
}

// NewAES 创建AES加密器
func NewAES(key, iv string, mode EncryptMode, padding PaddingMode, keySize AESKeySize) (*AESCrypto, error) {
	keyBytes := PrepareKey(key, int(keySize))
	
	var ivBytes []byte
	var err error
	if mode != ECB {
		ivBytes, err = PrepareIV(iv, 16) // AES块大小16字节
		if err != nil {
			return nil, err
		}
	}
	
	return &AESCrypto{
		key:     keyBytes,
		iv:      ivBytes,
		mode:    mode,
		padding: padding,
		keySize: keySize,
	}, nil
}

// Encrypt 加密
func (a *AESCrypto) Encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, err
	}
	
	// 添加填充
	plaintext, err = Padding(plaintext, block.BlockSize(), a.padding)
	if err != nil {
		return nil, err
	}
	
	ciphertext := make([]byte, len(plaintext))
	
	switch a.mode {
	case ECB:
		return a.encryptECB(block, plaintext)
	case CBC:
		mode := cipher.NewCBCEncrypter(block, a.iv)
		mode.CryptBlocks(ciphertext, plaintext)
		return ciphertext, nil
	case CTR:
		stream := cipher.NewCTR(block, a.iv)
		stream.XORKeyStream(ciphertext, plaintext)
		return ciphertext, nil
	case CFB:
		stream := cipher.NewCFBEncrypter(block, a.iv)
		stream.XORKeyStream(ciphertext, plaintext)
		return ciphertext, nil
	case OFB:
		stream := cipher.NewOFB(block, a.iv)
		stream.XORKeyStream(ciphertext, plaintext)
		return ciphertext, nil
	default:
		return nil, errors.New("不支持的加密模式")
	}
}

// Decrypt 解密
func (a *AESCrypto) Decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, err
	}
	
	if len(ciphertext) < block.BlockSize() {
		return nil, errors.New("密文长度太短")
	}
	
	plaintext := make([]byte, len(ciphertext))
	
	switch a.mode {
	case ECB:
		plaintext, err = a.decryptECB(block, ciphertext)
		if err != nil {
			return nil, err
		}
	case CBC:
		mode := cipher.NewCBCDecrypter(block, a.iv)
		mode.CryptBlocks(plaintext, ciphertext)
	case CTR:
		stream := cipher.NewCTR(block, a.iv)
		stream.XORKeyStream(plaintext, ciphertext)
	case CFB:
		stream := cipher.NewCFBDecrypter(block, a.iv)
		stream.XORKeyStream(plaintext, ciphertext)
	case OFB:
		stream := cipher.NewOFB(block, a.iv)
		stream.XORKeyStream(plaintext, ciphertext)
	default:
		return nil, errors.New("不支持的解密模式")
	}
	
	// 移除填充
	return UnPadding(plaintext, a.padding)
}

// EncryptToString 加密并输出为字符串
func (a *AESCrypto) EncryptToString(plaintext string, format OutputFormat) (string, error) {
	ciphertext, err := a.Encrypt([]byte(plaintext))
	if err != nil {
		return "", err
	}
	return FormatOutput(ciphertext, format), nil
}

// DecryptFromString 从字符串解密
func (a *AESCrypto) DecryptFromString(ciphertext string, format OutputFormat) (string, error) {
	data, err := ParseInput(ciphertext, format)
	if err != nil {
		return "", err
	}
	
	plaintext, err := a.Decrypt(data)
	if err != nil {
		return "", err
	}
	
	return string(plaintext), nil
}

// ECB模式加密
func (a *AESCrypto) encryptECB(block cipher.Block, src []byte) ([]byte, error) {
	blockSize := block.BlockSize()
	encrypted := make([]byte, len(src))
	
	for i := 0; i < len(src); i += blockSize {
		block.Encrypt(encrypted[i:i+blockSize], src[i:i+blockSize])
	}
	
	return encrypted, nil
}

// ECB模式解密
func (a *AESCrypto) decryptECB(block cipher.Block, src []byte) ([]byte, error) {
	blockSize := block.BlockSize()
	decrypted := make([]byte, len(src))
	
	for i := 0; i < len(src); i += blockSize {
		block.Decrypt(decrypted[i:i+blockSize], src[i:i+blockSize])
	}
	
	return decrypted, nil
}
