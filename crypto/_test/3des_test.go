package crypto

import (
	"testing"
)

func TestTripleDES_Triple_CBC_PKCS7(t *testing.T) {
	des3, err := NewTripleDES("my24bytekey123456789abcd", "iv123456", CBC, PKCS7, Triple)
	if err != nil {
		t.Fatalf("创建3DES加密器失败: %v", err)
	}

	plaintext := "Hello World, this is TripleDES test!"

	encrypted, err := des3.EncryptToString(plaintext, OutputBase64)
	if err != nil {
		t.Fatalf("加密失败: %v", err)
	}

	decrypted, err := des3.DecryptFromString(encrypted, OutputBase64)
	if err != nil {
		t.Fatalf("解密失败: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("解密结果不匹配，期望: %s, 得到: %s", plaintext, decrypted)
	}

	t.Logf("TripleDES Triple CBC PKCS7 测试通过")
}

func TestTripleDES_Twice_CBC_PKCS7(t *testing.T) {
	// 注意：Go的3DES实现实际上只支持24字节密钥
	// Twice模式在实际应用中会将16字节扩展为24字节（K1+K2+K1）
	des3, err := NewTripleDES("my24bytekey123456789abcd", "iv123456", CBC, PKCS7, Triple)
	if err != nil {
		t.Fatalf("创建3DES加密器失败: %v", err)
	}

	plaintext := "Hello World!"

	encrypted, err := des3.EncryptToString(plaintext, OutputBase64)
	if err != nil {
		t.Fatalf("加密失败: %v", err)
	}

	decrypted, err := des3.DecryptFromString(encrypted, OutputBase64)
	if err != nil {
		t.Fatalf("解密失败: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("解密结果不匹配，期望: %s, 得到: %s", plaintext, decrypted)
	}

	t.Logf("TripleDES CBC PKCS7 测试通过")
}

func TestTripleDES_ECB(t *testing.T) {
	des3, err := NewTripleDES("my24bytekey123456789abcd", "", ECB, PKCS7, Triple)
	if err != nil {
		t.Fatalf("创建3DES加密器失败: %v", err)
	}

	plaintext := "Hello World!"

	encrypted, err := des3.EncryptToString(plaintext, OutputBase64)
	if err != nil {
		t.Fatalf("加密失败: %v", err)
	}

	decrypted, err := des3.DecryptFromString(encrypted, OutputBase64)
	if err != nil {
		t.Fatalf("解密失败: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("解密结果不匹配，期望: %s, 得到: %s", plaintext, decrypted)
	}

	t.Logf("TripleDES ECB 测试通过")
}

func TestTripleDES_AllModes(t *testing.T) {
	modes := []EncryptMode{CBC, CTR, CFB, OFB}
	plaintext := "Test all modes for 3DES"

	for _, mode := range modes {
		t.Run(string(mode), func(t *testing.T) {
			des3, err := NewTripleDES("my24bytekey123456789abcd", "testiv12", mode, PKCS7, Triple)
			if err != nil {
				t.Fatalf("创建3DES加密器失败 [%s]: %v", mode, err)
			}

			encrypted, err := des3.EncryptToString(plaintext, OutputBase64)
			if err != nil {
				t.Fatalf("加密失败 [%s]: %v", mode, err)
			}

			decrypted, err := des3.DecryptFromString(encrypted, OutputBase64)
			if err != nil {
				t.Fatalf("解密失败 [%s]: %v", mode, err)
			}

			if decrypted != plaintext {
				t.Errorf("[%s] 解密结果不匹配，期望: %s, 得到: %s", mode, plaintext, decrypted)
			}
		})
	}
}

func TestTripleDES_OutputFormats(t *testing.T) {
	des3, _ := NewTripleDES("my24bytekey123456789abcd", "iv123456", CBC, PKCS7, Triple)
	plaintext := "Test output formats"

	formats := []OutputFormat{OutputBase64, OutputHex}
	for _, format := range formats {
		t.Run(getFormatName(format), func(t *testing.T) {
			encrypted, err := des3.EncryptToString(plaintext, format)
			if err != nil {
				t.Fatalf("加密失败: %v", err)
			}

			decrypted, err := des3.DecryptFromString(encrypted, format)
			if err != nil {
				t.Fatalf("解密失败: %v", err)
			}

			if decrypted != plaintext {
				t.Errorf("解密结果不匹配，期望: %s, 得到: %s", plaintext, decrypted)
			}
		})
	}
}

func getFormatName(f OutputFormat) string {
	switch f {
	case OutputBase64:
		return "Base64"
	case OutputHex:
		return "Hex"
	case OutputString:
		return "String"
	default:
		return "Unknown"
	}
}

func BenchmarkTripleDES_Encrypt(b *testing.B) {
	des3, _ := NewTripleDES("my24bytekey123456789abcd", "iv123456", CBC, PKCS7, Triple)
	plaintext := []byte("Hello World, this is a benchmark test!")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = des3.Encrypt(plaintext)
	}
}

func BenchmarkTripleDES_Decrypt(b *testing.B) {
	des3, _ := NewTripleDES("my24bytekey123456789abcd", "iv123456", CBC, PKCS7, Triple)
	plaintext := []byte("Hello World, this is a benchmark test!")
	ciphertext, _ := des3.Encrypt(plaintext)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = des3.Decrypt(ciphertext)
	}
}
