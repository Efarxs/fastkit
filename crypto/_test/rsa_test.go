package crypto

import (
	"testing"
)

const (
	testPublicKeyPEM = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAw8xnHqYyJ6NhYZOTNRBg
UYLbDJrJKHG9dmvU1GvLFvPnO+L0YmVfPmkRnkH9fYnUYKNlzVFdKqGmVdVVWwXN
R7B3nJM5YOvRnXPvE8KLPb1uPqmVGIq5WxRnqD8HxQkXqVDqKXj4S9GEwGLaNvJe
ZRjWvEFV9YGdOmI2d3dpxVODQqW7K4V5MfnMFxBpn6V1p5nqNGkPxQKlCqPJ7VrS
NhN9V8D8BvG3RHGvRUfPzN5nGvPGcNmvxNmPxGmv5GmPxNmPxNmPxGmvxGmvxGmv
xGmvxGmvxGmvxGmvxGmvxGmvxGmvxGmvxGmvxGmvxGmvxGmvxGmvxGmvxGmvxGmv
xQIDAQAB
-----END PUBLIC KEY-----`

	testPrivateKeyPEM = `-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAw8xnHqYyJ6NhYZOTNRBgUYLbDJrJKHG9dmvU1GvLFvPnO+L0
YmVfPmkRnkH9fYnUYKNlzVFdKqGmVdVVWwXNR7B3nJM5YOvRnXPvE8KLPb1uPqmV
GIq5WxRnqD8HxQkXqVDqKXj4S9GEwGLaNvJeZRjWvEFV9YGdOmI2d3dpxVODQqW7
K4V5MfnMFxBpn6V1p5nqNGkPxQKlCqPJ7VrSNhN9V8D8BvG3RHGvRUfPzN5nGvPG
cNmvxNmPxGmv5GmPxNmPxNmPxGmvxGmvxGmvxGmvxGmvxGmvxGmvxGmvxGmvxGmv
xGmvxGmvxGmvxGmvxGmvxGmvxGmvxGmvxGmvxQIDAQABAoIBAAq2XqLVmQPVIuOG
J7nrjqKJ9Hvnq5KtLHnPqBmFqPxbB9DGmN7HkC2LnGH5qGhqG5qG5qGhqG5qG5qG
hqG5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGh
qG5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhq
G5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhqG
5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhqG5
qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qECgYEA5qG5qG5qG
hqG5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGh
qG5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhq
G5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhqG
5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhqG5
qG5qECgYEA5qG5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qG
hqG5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGh
qG5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhq
G5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhqG
5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhqG5qG5qGhqG5
qG5qECgYBV0ECgYBV0ECgYBV0=
-----END RSA PRIVATE KEY-----`
)

func TestRSA_GenerateKeyPair(t *testing.T) {
	publicKey, privateKey, err := GenerateRSAKeyPair(2048)
	if err != nil {
		t.Fatalf("生成RSA密钥对失败: %v", err)
	}

	if publicKey == "" || privateKey == "" {
		t.Error("生成的密钥为空")
	}

	t.Logf("生成RSA密钥对成功")
	t.Logf("公钥长度: %d", len(publicKey))
	t.Logf("私钥长度: %d", len(privateKey))
}

func TestRSA_PublicEncrypt_PrivateDecrypt(t *testing.T) {
	// 生成密钥对
	publicKey, privateKey, err := GenerateRSAKeyPair(2048)
	if err != nil {
		t.Fatalf("生成RSA密钥对失败: %v", err)
	}

	rsa, err := NewRSA(publicKey, privateKey)
	if err != nil {
		t.Fatalf("创建RSA加密器失败: %v", err)
	}

	plaintext := "Hello RSA!"

	// 公钥加密
	encrypted, err := rsa.EncryptToStringWithPublicKey(plaintext, OutputBase64)
	if err != nil {
		t.Fatalf("公钥加密失败: %v", err)
	}

	// 私钥解密
	decrypted, err := rsa.DecryptFromStringWithPrivateKey(encrypted, OutputBase64)
	if err != nil {
		t.Fatalf("私钥解密失败: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("解密结果不匹配，期望: %s, 得到: %s", plaintext, decrypted)
	}

	t.Logf("RSA 公钥加密私钥解密测试通过")
}

func TestRSA_PrivateEncrypt_PublicDecrypt(t *testing.T) {
	// 生成密钥对
	publicKey, privateKey, err := GenerateRSAKeyPair(2048)
	if err != nil {
		t.Fatalf("生成RSA密钥对失败: %v", err)
	}

	rsa, err := NewRSA(publicKey, privateKey)
	if err != nil {
		t.Fatalf("创建RSA加密器失败: %v", err)
	}

	plaintext := "Hello RSA Signature!"

	// 私钥加密（签名）
	signature, err := rsa.EncryptToStringWithPrivateKey(plaintext, OutputBase64)
	if err != nil {
		t.Fatalf("私钥加密失败: %v", err)
	}

	t.Logf("RSA 签名: %s", signature)
	t.Logf("RSA 私钥加密（签名）测试通过")
}

func TestRSA_OutputFormats(t *testing.T) {
	publicKey, privateKey, _ := GenerateRSAKeyPair(2048)
	rsa, _ := NewRSA(publicKey, privateKey)

	plaintext := "Test output formats"

	formats := []OutputFormat{OutputBase64, OutputHex}
	for _, format := range formats {
		t.Run(getFormatName(format), func(t *testing.T) {
			encrypted, err := rsa.EncryptToStringWithPublicKey(plaintext, format)
			if err != nil {
				t.Fatalf("加密失败: %v", err)
			}

			decrypted, err := rsa.DecryptFromStringWithPrivateKey(encrypted, format)
			if err != nil {
				t.Fatalf("解密失败: %v", err)
			}

			if decrypted != plaintext {
				t.Errorf("解密结果不匹配，期望: %s, 得到: %s", plaintext, decrypted)
			}
		})
	}
}

func TestRSA_OnlyPublicKey(t *testing.T) {
	publicKey, _, _ := GenerateRSAKeyPair(2048)
	rsa, err := NewRSA(publicKey, "")
	if err != nil {
		t.Fatalf("创建RSA加密器失败: %v", err)
	}

	plaintext := "Test with only public key"

	// 应该能加密
	encrypted, err := rsa.EncryptToStringWithPublicKey(plaintext, OutputBase64)
	if err != nil {
		t.Fatalf("公钥加密失败: %v", err)
	}

	// 不应该能解密（没有私钥）
	_, err = rsa.DecryptFromStringWithPrivateKey(encrypted, OutputBase64)
	if err == nil {
		t.Error("没有私钥也能解密，这不应该发生")
	}

	t.Logf("只有公钥测试通过")
}

func TestRSA_OnlyPrivateKey(t *testing.T) {
	_, privateKey, _ := GenerateRSAKeyPair(2048)
	rsa, err := NewRSA("", privateKey)
	if err != nil {
		t.Fatalf("创建RSA加密器失败: %v", err)
	}

	plaintext := "Test with only private key"

	// 不应该能用公钥加密（没有公钥）
	_, err = rsa.EncryptToStringWithPublicKey(plaintext, OutputBase64)
	if err == nil {
		t.Error("没有公钥也能加密，这不应该发生")
	}

	// 应该能用私钥加密（签名）
	_, err = rsa.EncryptToStringWithPrivateKey(plaintext, OutputBase64)
	if err != nil {
		t.Fatalf("私钥加密失败: %v", err)
	}

	t.Logf("只有私钥测试通过")
}

func BenchmarkRSA_Encrypt(b *testing.B) {
	publicKey, privateKey, _ := GenerateRSAKeyPair(2048)
	rsa, _ := NewRSA(publicKey, privateKey)
	plaintext := []byte("Hello World!")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = rsa.EncryptWithPublicKey(plaintext)
	}
}

func BenchmarkRSA_Decrypt(b *testing.B) {
	publicKey, privateKey, _ := GenerateRSAKeyPair(2048)
	rsa, _ := NewRSA(publicKey, privateKey)
	plaintext := []byte("Hello World!")
	ciphertext, _ := rsa.EncryptWithPublicKey(plaintext)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = rsa.DecryptWithPrivateKey(ciphertext)
	}
}

func BenchmarkRSA_KeyGeneration(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = GenerateRSAKeyPair(2048)
	}
}
