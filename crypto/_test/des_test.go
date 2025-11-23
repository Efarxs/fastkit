package crypto

import (
	"testing"
)

func TestDES_CBC_PKCS7(t *testing.T) {
	des, err := NewDES("mykey123", "iv123456", CBC, PKCS7)
	if err != nil {
		t.Fatalf("创建DES加密器失败: %v", err)
	}

	plaintext := "Hello World, this is DES test!"

	// 测试Base64输出
	encrypted, err := des.EncryptToString(plaintext, OutputBase64)
	if err != nil {
		t.Fatalf("加密失败: %v", err)
	}

	decrypted, err := des.DecryptFromString(encrypted, OutputBase64)
	if err != nil {
		t.Fatalf("解密失败: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("解密结果不匹配，期望: %s, 得到: %s", plaintext, decrypted)
	}

	t.Logf("DES CBC PKCS7 测试通过")
}

func TestDES_ECB_PKCS7(t *testing.T) {
	des, err := NewDES("mykey123", "", ECB, PKCS7)
	if err != nil {
		t.Fatalf("创建DES加密器失败: %v", err)
	}

	plaintext := "Hello World!"

	encrypted, err := des.EncryptToString(plaintext, OutputBase64)
	if err != nil {
		t.Fatalf("加密失败: %v", err)
	}

	decrypted, err := des.DecryptFromString(encrypted, OutputBase64)
	if err != nil {
		t.Fatalf("解密失败: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("解密结果不匹配，期望: %s, 得到: %s", plaintext, decrypted)
	}

	t.Logf("DES ECB PKCS7 测试通过")
}

func TestDES_CTR(t *testing.T) {
	// CTR模式不需要填充，但数据可以是任意长度
	des, err := NewDES("mykey123", "iv123456", CTR, PKCS7)
	if err != nil {
		t.Fatalf("创建DES加密器失败: %v", err)
	}

	plaintext := "Hello World!"

	encrypted, err := des.EncryptToString(plaintext, OutputHex)
	if err != nil {
		t.Fatalf("加密失败: %v", err)
	}

	decrypted, err := des.DecryptFromString(encrypted, OutputHex)
	if err != nil {
		t.Fatalf("解密失败: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("解密结果不匹配，期望: %s, 得到: %s", plaintext, decrypted)
	}

	t.Logf("DES CTR 测试通过")
}

func TestDES_AllModes(t *testing.T) {
	modes := []EncryptMode{CBC, CTR, CFB, OFB}
	plaintext := "Test all modes"

	for _, mode := range modes {
		t.Run(string(mode), func(t *testing.T) {
			des, err := NewDES("testkey1", "testiv12", mode, PKCS7)
			if err != nil {
				t.Fatalf("创建DES加密器失败 [%s]: %v", mode, err)
			}

			encrypted, err := des.EncryptToString(plaintext, OutputBase64)
			if err != nil {
				t.Fatalf("加密失败 [%s]: %v", mode, err)
			}

			decrypted, err := des.DecryptFromString(encrypted, OutputBase64)
			if err != nil {
				t.Fatalf("解密失败 [%s]: %v", mode, err)
			}

			if decrypted != plaintext {
				t.Errorf("[%s] 解密结果不匹配，期望: %s, 得到: %s", mode, plaintext, decrypted)
			}
		})
	}
}

func TestDES_AllPaddings(t *testing.T) {
	paddings := []PaddingMode{PKCS7, Zero, ISO10126, ANSIX923}
	plaintext := "Test paddings"

	for _, padding := range paddings {
		t.Run(getPaddingName(padding), func(t *testing.T) {
			des, err := NewDES("testkey1", "testiv12", CBC, padding)
			if err != nil {
				t.Fatalf("创建DES加密器失败: %v", err)
			}

			encrypted, err := des.EncryptToString(plaintext, OutputBase64)
			if err != nil {
				t.Fatalf("加密失败: %v", err)
			}

			decrypted, err := des.DecryptFromString(encrypted, OutputBase64)
			if err != nil {
				t.Fatalf("解密失败: %v", err)
			}

			if decrypted != plaintext {
				t.Errorf("解密结果不匹配，期望: %s, 得到: %s", plaintext, decrypted)
			}
		})
	}
}

func getPaddingName(p PaddingMode) string {
	switch p {
	case PKCS7:
		return "PKCS7"
	case Zero:
		return "Zero"
	case ISO10126:
		return "ISO10126"
	case ANSIX923:
		return "ANSIX923"
	case NoPadding:
		return "NoPadding"
	default:
		return "Unknown"
	}
}

func BenchmarkDES_Encrypt(b *testing.B) {
	des, _ := NewDES("mykey123", "iv123456", CBC, PKCS7)
	plaintext := []byte("Hello World, this is a benchmark test!")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = des.Encrypt(plaintext)
	}
}

func BenchmarkDES_Decrypt(b *testing.B) {
	des, _ := NewDES("mykey123", "iv123456", CBC, PKCS7)
	plaintext := []byte("Hello World, this is a benchmark test!")
	ciphertext, _ := des.Encrypt(plaintext)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = des.Decrypt(ciphertext)
	}
}
