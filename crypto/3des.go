package crypto

import (
	"crypto/cipher"
	"crypto/des"
	"errors"
)

// TripleDESType 3DES类型
type TripleDESType int

const (
	Twice  TripleDESType = iota // 2-key 3DES (16字节密钥)
	Triple                       // 3-key 3DES (24字节密钥)
)

// TripleDESCrypto 3DES加密器
type TripleDESCrypto struct {
	key     []byte
	iv      []byte
	mode    EncryptMode
	padding PaddingMode
	desType TripleDESType
}

// NewTripleDES 创建3DES加密器
func NewTripleDES(key, iv string, mode EncryptMode, padding PaddingMode, desType TripleDESType) (*TripleDESCrypto, error) {
	var keyLength int
	if desType == Twice {
		keyLength = 16 // 2-key 3DES
	} else {
		keyLength = 24 // 3-key 3DES
	}
	
	keyBytes := PrepareKey(key, keyLength)
	
	var ivBytes []byte
	var err error
	if mode != ECB {
		ivBytes, err = PrepareIV(iv, 8) // 3DES块大小8字节
		if err != nil {
			return nil, err
		}
	}
	
	return &TripleDESCrypto{
		key:     keyBytes,
		iv:      ivBytes,
		mode:    mode,
		padding: padding,
		desType: desType,
	}, nil
}

// Encrypt 加密
func (t *TripleDESCrypto) Encrypt(plaintext []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(t.key)
	if err != nil {
		return nil, err
	}
	
	// 添加填充
	plaintext, err = Padding(plaintext, block.BlockSize(), t.padding)
	if err != nil {
		return nil, err
	}
	
	ciphertext := make([]byte, len(plaintext))
	
	switch t.mode {
	case ECB:
		return t.encryptECB(block, plaintext)
	case CBC:
		mode := cipher.NewCBCEncrypter(block, t.iv)
		mode.CryptBlocks(ciphertext, plaintext)
		return ciphertext, nil
	case CTR:
		stream := cipher.NewCTR(block, t.iv)
		stream.XORKeyStream(ciphertext, plaintext)
		return ciphertext, nil
	case CFB:
		stream := cipher.NewCFBEncrypter(block, t.iv)
		stream.XORKeyStream(ciphertext, plaintext)
		return ciphertext, nil
	case OFB:
		stream := cipher.NewOFB(block, t.iv)
		stream.XORKeyStream(ciphertext, plaintext)
		return ciphertext, nil
	default:
		return nil, errors.New("不支持的加密模式")
	}
}

// Decrypt 解密
func (t *TripleDESCrypto) Decrypt(ciphertext []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(t.key)
	if err != nil {
		return nil, err
	}
	
	if len(ciphertext) < block.BlockSize() {
		return nil, errors.New("密文长度太短")
	}
	
	plaintext := make([]byte, len(ciphertext))
	
	switch t.mode {
	case ECB:
		plaintext, err = t.decryptECB(block, ciphertext)
		if err != nil {
			return nil, err
		}
	case CBC:
		mode := cipher.NewCBCDecrypter(block, t.iv)
		mode.CryptBlocks(plaintext, ciphertext)
	case CTR:
		stream := cipher.NewCTR(block, t.iv)
		stream.XORKeyStream(plaintext, ciphertext)
	case CFB:
		stream := cipher.NewCFBDecrypter(block, t.iv)
		stream.XORKeyStream(plaintext, ciphertext)
	case OFB:
		stream := cipher.NewOFB(block, t.iv)
		stream.XORKeyStream(plaintext, ciphertext)
	default:
		return nil, errors.New("不支持的解密模式")
	}
	
	// 移除填充
	return UnPadding(plaintext, t.padding)
}

// EncryptToString 加密并输出为字符串
func (t *TripleDESCrypto) EncryptToString(plaintext string, format OutputFormat) (string, error) {
	ciphertext, err := t.Encrypt([]byte(plaintext))
	if err != nil {
		return "", err
	}
	return FormatOutput(ciphertext, format), nil
}

// DecryptFromString 从字符串解密
func (t *TripleDESCrypto) DecryptFromString(ciphertext string, format OutputFormat) (string, error) {
	data, err := ParseInput(ciphertext, format)
	if err != nil {
		return "", err
	}
	
	plaintext, err := t.Decrypt(data)
	if err != nil {
		return "", err
	}
	
	return string(plaintext), nil
}

// ECB模式加密
func (t *TripleDESCrypto) encryptECB(block cipher.Block, src []byte) ([]byte, error) {
	blockSize := block.BlockSize()
	encrypted := make([]byte, len(src))
	
	for i := 0; i < len(src); i += blockSize {
		block.Encrypt(encrypted[i:i+blockSize], src[i:i+blockSize])
	}
	
	return encrypted, nil
}

// ECB模式解密
func (t *TripleDESCrypto) decryptECB(block cipher.Block, src []byte) ([]byte, error) {
	blockSize := block.BlockSize()
	decrypted := make([]byte, len(src))
	
	for i := 0; i < len(src); i += blockSize {
		block.Decrypt(decrypted[i:i+blockSize], src[i:i+blockSize])
	}
	
	return decrypted, nil
}
