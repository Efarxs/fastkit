package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
)

// RSACrypto RSA加密器
type RSACrypto struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

// NewRSA 创建RSA加密器（从PEM格式）
func NewRSA(publicKeyPEM, privateKeyPEM string) (*RSACrypto, error) {
	r := &RSACrypto{}
	
	var err error
	if publicKeyPEM != "" {
		r.publicKey, err = parsePublicKey(publicKeyPEM)
		if err != nil {
			return nil, err
		}
	}
	
	if privateKeyPEM != "" {
		r.privateKey, err = parsePrivateKey(privateKeyPEM)
		if err != nil {
			return nil, err
		}
	}
	
	return r, nil
}

// NewRSAFromBase64 从Base64格式创建RSA加密器
func NewRSAFromBase64(publicKeyBase64, privateKeyBase64 string) (*RSACrypto, error) {
	r := &RSACrypto{}
	
	if publicKeyBase64 != "" {
		publicKeyPEM, err := base64.StdEncoding.DecodeString(publicKeyBase64)
		if err != nil {
			return nil, err
		}
		r.publicKey, err = parsePublicKey(string(publicKeyPEM))
		if err != nil {
			return nil, err
		}
	}
	
	if privateKeyBase64 != "" {
		privateKeyPEM, err := base64.StdEncoding.DecodeString(privateKeyBase64)
		if err != nil {
			return nil, err
		}
		r.privateKey, err = parsePrivateKey(string(privateKeyPEM))
		if err != nil {
			return nil, err
		}
	}
	
	return r, nil
}

// EncryptWithPublicKey 使用公钥加密
func (r *RSACrypto) EncryptWithPublicKey(plaintext []byte) ([]byte, error) {
	if r.publicKey == nil {
		return nil, errors.New("公钥未设置")
	}
	
	return rsa.EncryptPKCS1v15(rand.Reader, r.publicKey, plaintext)
}

// DecryptWithPrivateKey 使用私钥解密
func (r *RSACrypto) DecryptWithPrivateKey(ciphertext []byte) ([]byte, error) {
	if r.privateKey == nil {
		return nil, errors.New("私钥未设置")
	}
	
	return rsa.DecryptPKCS1v15(rand.Reader, r.privateKey, ciphertext)
}

// EncryptWithPrivateKey 使用私钥加密（签名）
func (r *RSACrypto) EncryptWithPrivateKey(plaintext []byte) ([]byte, error) {
	if r.privateKey == nil {
		return nil, errors.New("私钥未设置")
	}
	
	// 使用私钥进行加密实际上是签名操作
	// 这里我们使用较低级别的操作来实现
	return rsa.SignPKCS1v15(rand.Reader, r.privateKey, 0, plaintext)
}

// DecryptWithPublicKey 使用公钥解密（验证签名）
func (r *RSACrypto) DecryptWithPublicKey(ciphertext []byte) error {
	if r.publicKey == nil {
		return errors.New("公钥未设置")
	}
	
	// 验证签名
	return rsa.VerifyPKCS1v15(r.publicKey, 0, ciphertext, ciphertext)
}

// EncryptToStringWithPublicKey 使用公钥加密并输出为字符串
func (r *RSACrypto) EncryptToStringWithPublicKey(plaintext string, format OutputFormat) (string, error) {
	ciphertext, err := r.EncryptWithPublicKey([]byte(plaintext))
	if err != nil {
		return "", err
	}
	return FormatOutput(ciphertext, format), nil
}

// DecryptFromStringWithPrivateKey 使用私钥从字符串解密
func (r *RSACrypto) DecryptFromStringWithPrivateKey(ciphertext string, format OutputFormat) (string, error) {
	data, err := ParseInput(ciphertext, format)
	if err != nil {
		return "", err
	}
	
	plaintext, err := r.DecryptWithPrivateKey(data)
	if err != nil {
		return "", err
	}
	
	return string(plaintext), nil
}

// EncryptToStringWithPrivateKey 使用私钥加密并输出为字符串
func (r *RSACrypto) EncryptToStringWithPrivateKey(plaintext string, format OutputFormat) (string, error) {
	ciphertext, err := r.EncryptWithPrivateKey([]byte(plaintext))
	if err != nil {
		return "", err
	}
	return FormatOutput(ciphertext, format), nil
}

// parsePublicKey 解析公钥
func parsePublicKey(publicKeyPEM string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return nil, errors.New("无法解析PEM格式的公钥")
	}
	
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("不是RSA公钥")
	}
	
	return rsaPub, nil
}

// parsePrivateKey 解析私钥
func parsePrivateKey(privateKeyPEM string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, errors.New("无法解析PEM格式的私钥")
	}
	
	// 尝试PKCS1格式
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err == nil {
		return priv, nil
	}
	
	// 尝试PKCS8格式
	privInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	
	rsaPriv, ok := privInterface.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("不是RSA私钥")
	}
	
	return rsaPriv, nil
}

// GenerateRSAKeyPair 生成RSA密钥对
func GenerateRSAKeyPair(bits int) (publicKeyPEM, privateKeyPEM string, err error) {
	// 生成密钥对
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", err
	}
	
	// 编码私钥为PEM格式
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	privateKeyPEM = string(pem.EncodeToMemory(privateKeyBlock))
	
	// 编码公钥为PEM格式
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", err
	}
	publicKeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	publicKeyPEM = string(pem.EncodeToMemory(publicKeyBlock))
	
	return publicKeyPEM, privateKeyPEM, nil
}
