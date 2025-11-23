package crypto

import (
	"crypto/cipher"
	"crypto/des"
	"errors"
)

// DESCrypto DES加密器
type DESCrypto struct {
	key     []byte
	iv      []byte
	mode    EncryptMode
	padding PaddingMode
}

// NewDES 创建DES加密器
func NewDES(key, iv string, mode EncryptMode, padding PaddingMode) (*DESCrypto, error) {
	keyBytes := PrepareKey(key, 8) // DES密钥固定8字节
	
	var ivBytes []byte
	var err error
	if mode != ECB {
		ivBytes, err = PrepareIV(iv, 8) // DES块大小8字节
		if err != nil {
			return nil, err
		}
	}
	
	return &DESCrypto{
		key:     keyBytes,
		iv:      ivBytes,
		mode:    mode,
		padding: padding,
	}, nil
}

// Encrypt 加密
func (d *DESCrypto) Encrypt(plaintext []byte) ([]byte, error) {
	block, err := des.NewCipher(d.key)
	if err != nil {
		return nil, err
	}
	
	// 添加填充
	plaintext, err = Padding(plaintext, block.BlockSize(), d.padding)
	if err != nil {
		return nil, err
	}
	
	ciphertext := make([]byte, len(plaintext))
	
	switch d.mode {
	case ECB:
		return d.encryptECB(block, plaintext)
	case CBC:
		mode := cipher.NewCBCEncrypter(block, d.iv)
		mode.CryptBlocks(ciphertext, plaintext)
		return ciphertext, nil
	case CTR:
		stream := cipher.NewCTR(block, d.iv)
		stream.XORKeyStream(ciphertext, plaintext)
		return ciphertext, nil
	case CFB:
		stream := cipher.NewCFBEncrypter(block, d.iv)
		stream.XORKeyStream(ciphertext, plaintext)
		return ciphertext, nil
	case OFB:
		stream := cipher.NewOFB(block, d.iv)
		stream.XORKeyStream(ciphertext, plaintext)
		return ciphertext, nil
	default:
		return nil, errors.New("不支持的加密模式")
	}
}

// Decrypt 解密
func (d *DESCrypto) Decrypt(ciphertext []byte) ([]byte, error) {
	block, err := des.NewCipher(d.key)
	if err != nil {
		return nil, err
	}
	
	if len(ciphertext) < block.BlockSize() {
		return nil, errors.New("密文长度太短")
	}
	
	plaintext := make([]byte, len(ciphertext))
	
	switch d.mode {
	case ECB:
		plaintext, err = d.decryptECB(block, ciphertext)
		if err != nil {
			return nil, err
		}
	case CBC:
		mode := cipher.NewCBCDecrypter(block, d.iv)
		mode.CryptBlocks(plaintext, ciphertext)
	case CTR:
		stream := cipher.NewCTR(block, d.iv)
		stream.XORKeyStream(plaintext, ciphertext)
	case CFB:
		stream := cipher.NewCFBDecrypter(block, d.iv)
		stream.XORKeyStream(plaintext, ciphertext)
	case OFB:
		stream := cipher.NewOFB(block, d.iv)
		stream.XORKeyStream(plaintext, ciphertext)
	default:
		return nil, errors.New("不支持的解密模式")
	}
	
	// 移除填充
	return UnPadding(plaintext, d.padding)
}

// EncryptToString 加密并输出为字符串
func (d *DESCrypto) EncryptToString(plaintext string, format OutputFormat) (string, error) {
	ciphertext, err := d.Encrypt([]byte(plaintext))
	if err != nil {
		return "", err
	}
	return FormatOutput(ciphertext, format), nil
}

// DecryptFromString 从字符串解密
func (d *DESCrypto) DecryptFromString(ciphertext string, format OutputFormat) (string, error) {
	data, err := ParseInput(ciphertext, format)
	if err != nil {
		return "", err
	}
	
	plaintext, err := d.Decrypt(data)
	if err != nil {
		return "", err
	}
	
	return string(plaintext), nil
}

// ECB模式加密
func (d *DESCrypto) encryptECB(block cipher.Block, src []byte) ([]byte, error) {
	blockSize := block.BlockSize()
	encrypted := make([]byte, len(src))
	
	for i := 0; i < len(src); i += blockSize {
		block.Encrypt(encrypted[i:i+blockSize], src[i:i+blockSize])
	}
	
	return encrypted, nil
}

// ECB模式解密
func (d *DESCrypto) decryptECB(block cipher.Block, src []byte) ([]byte, error) {
	blockSize := block.BlockSize()
	decrypted := make([]byte, len(src))
	
	for i := 0; i < len(src); i += blockSize {
		block.Decrypt(decrypted[i:i+blockSize], src[i:i+blockSize])
	}
	
	return decrypted, nil
}
