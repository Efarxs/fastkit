package crypto

import (
	"testing"
)

func TestAES128_CBC_PKCS7(t *testing.T) {
	aes, err := NewAES("1234567890123456", "iv1234567890abcd", CBC, PKCS7, AES128)
	if err != nil {
		t.Fatalf("创建AES加密器失败: %v", err)
	}

	plaintext := "Hello World, this is AES-128 test!"

	encrypted, err := aes.EncryptToString(plaintext, OutputBase64)
	if err != nil {
		t.Fatalf("加密失败: %v", err)
	}

	decrypted, err := aes.DecryptFromString(encrypted, OutputBase64)
	if err != nil {
		t.Fatalf("解密失败: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("解密结果不匹配，期望: %s, 得到: %s", plaintext, decrypted)
	}

	t.Logf("AES-128 CBC PKCS7 测试通过")
}

func TestAES192_CBC_PKCS7(t *testing.T) {
	aes, err := NewAES("123456789012345678901234", "iv1234567890abcd", CBC, PKCS7, AES192)
	if err != nil {
		t.Fatalf("创建AES加密器失败: %v", err)
	}

	plaintext := "Hello World, this is AES-192 test!"

	encrypted, err := aes.EncryptToString(plaintext, OutputBase64)
	if err != nil {
		t.Fatalf("加密失败: %v", err)
	}

	decrypted, err := aes.DecryptFromString(encrypted, OutputBase64)
	if err != nil {
		t.Fatalf("解密失败: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("解密结果不匹配，期望: %s, 得到: %s", plaintext, decrypted)
	}

	t.Logf("AES-192 CBC PKCS7 测试通过")
}

func TestAES256_CBC_PKCS7(t *testing.T) {
	aes, err := NewAES("12345678901234567890123456789012", "iv1234567890abcd", CBC, PKCS7, AES256)
	if err != nil {
		t.Fatalf("创建AES加密器失败: %v", err)
	}

	plaintext := "Hello World, this is AES-256 test!"

	encrypted, err := aes.EncryptToString(plaintext, OutputBase64)
	if err != nil {
		t.Fatalf("加密失败: %v", err)
	}

	decrypted, err := aes.DecryptFromString(encrypted, OutputBase64)
	if err != nil {
		t.Fatalf("解密失败: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("解密结果不匹配，期望: %s, 得到: %s", plaintext, decrypted)
	}

	t.Logf("AES-256 CBC PKCS7 测试通过")
}

func TestAES_ECB(t *testing.T) {
	aes, err := NewAES("1234567890123456", "", ECB, PKCS7, AES128)
	if err != nil {
		t.Fatalf("创建AES加密器失败: %v", err)
	}

	plaintext := "Hello World!"

	encrypted, err := aes.EncryptToString(plaintext, OutputBase64)
	if err != nil {
		t.Fatalf("加密失败: %v", err)
	}

	decrypted, err := aes.DecryptFromString(encrypted, OutputBase64)
	if err != nil {
		t.Fatalf("解密失败: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("解密结果不匹配，期望: %s, 得到: %s", plaintext, decrypted)
	}

	t.Logf("AES ECB 测试通过")
}

func TestAES_AllModes(t *testing.T) {
	modes := []EncryptMode{CBC, CTR, CFB, OFB}
	plaintext := "Test all AES modes"

	for _, mode := range modes {
		t.Run(string(mode), func(t *testing.T) {
			aes, err := NewAES("1234567890123456", "iv1234567890abcd", mode, PKCS7, AES128)
			if err != nil {
				t.Fatalf("创建AES加密器失败 [%s]: %v", mode, err)
			}

			encrypted, err := aes.EncryptToString(plaintext, OutputBase64)
			if err != nil {
				t.Fatalf("加密失败 [%s]: %v", mode, err)
			}

			decrypted, err := aes.DecryptFromString(encrypted, OutputBase64)
			if err != nil {
				t.Fatalf("解密失败 [%s]: %v", mode, err)
			}

			if decrypted != plaintext {
				t.Errorf("[%s] 解密结果不匹配，期望: %s, 得到: %s", mode, plaintext, decrypted)
			}
		})
	}
}

func TestAES_AllKeySizes(t *testing.T) {
	keySizes := []struct {
		size AESKeySize
		key  string
		name string
	}{
		{AES128, "1234567890123456", "AES-128"},
		{AES192, "123456789012345678901234", "AES-192"},
		{AES256, "12345678901234567890123456789012", "AES-256"},
	}

	plaintext := "Test all key sizes"

	for _, ks := range keySizes {
		t.Run(ks.name, func(t *testing.T) {
			aes, err := NewAES(ks.key, "iv1234567890abcd", CBC, PKCS7, ks.size)
			if err != nil {
				t.Fatalf("创建AES加密器失败 [%s]: %v", ks.name, err)
			}

			encrypted, err := aes.EncryptToString(plaintext, OutputBase64)
			if err != nil {
				t.Fatalf("加密失败 [%s]: %v", ks.name, err)
			}

			decrypted, err := aes.DecryptFromString(encrypted, OutputBase64)
			if err != nil {
				t.Fatalf("解密失败 [%s]: %v", ks.name, err)
			}

			if decrypted != plaintext {
				t.Errorf("[%s] 解密结果不匹配，期望: %s, 得到: %s", ks.name, plaintext, decrypted)
			}
		})
	}
}

func TestAES_AllPaddings(t *testing.T) {
	paddings := []PaddingMode{PKCS7, Zero, ISO10126, ANSIX923}
	plaintext := "Test AES paddings"

	for _, padding := range paddings {
		t.Run(getPaddingName(padding), func(t *testing.T) {
			aes, err := NewAES("1234567890123456", "iv1234567890abcd", CBC, padding, AES128)
			if err != nil {
				t.Fatalf("创建AES加密器失败: %v", err)
			}

			encrypted, err := aes.EncryptToString(plaintext, OutputBase64)
			if err != nil {
				t.Fatalf("加密失败: %v", err)
			}

			decrypted, err := aes.DecryptFromString(encrypted, OutputBase64)
			if err != nil {
				t.Fatalf("解密失败: %v", err)
			}

			if decrypted != plaintext {
				t.Errorf("解密结果不匹配，期望: %s, 得到: %s", plaintext, decrypted)
			}
		})
	}
}

func TestAES_LongKey(t *testing.T) {
	// 测试超长密钥（会自动截取）
	longKey := "this-is-a-very-long-key-that-exceeds-32-bytes-and-should-be-truncated-automatically"
	aes, err := NewAES(longKey, "iv1234567890abcd", CBC, PKCS7, AES256)
	if err != nil {
		t.Fatalf("创建AES加密器失败: %v", err)
	}

	plaintext := "Test with long key"

	encrypted, err := aes.EncryptToString(plaintext, OutputBase64)
	if err != nil {
		t.Fatalf("加密失败: %v", err)
	}

	decrypted, err := aes.DecryptFromString(encrypted, OutputBase64)
	if err != nil {
		t.Fatalf("解密失败: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("解密结果不匹配，期望: %s, 得到: %s", plaintext, decrypted)
	}

	t.Logf("超长密钥测试通过")
}

func BenchmarkAES128_Encrypt(b *testing.B) {
	aes, _ := NewAES("1234567890123456", "iv1234567890abcd", CBC, PKCS7, AES128)
	plaintext := []byte("Hello World, this is a benchmark test!")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = aes.Encrypt(plaintext)
	}
}

func BenchmarkAES256_Encrypt(b *testing.B) {
	aes, _ := NewAES("12345678901234567890123456789012", "iv1234567890abcd", CBC, PKCS7, AES256)
	plaintext := []byte("Hello World, this is a benchmark test!")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = aes.Encrypt(plaintext)
	}
}

func BenchmarkAES128_Decrypt(b *testing.B) {
	aes, _ := NewAES("1234567890123456", "iv1234567890abcd", CBC, PKCS7, AES128)
	plaintext := []byte("Hello World, this is a benchmark test!")
	ciphertext, _ := aes.Encrypt(plaintext)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = aes.Decrypt(ciphertext)
	}
}
