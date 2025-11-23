package crypto

import (
	"testing"
)

func TestMD5(t *testing.T) {
	data := "Hello World"
	expected := "b10a8db164e0754105b7a99be72e3fe5" // MD5 of "Hello World"

	result := MD5(data)
	if result != expected {
		t.Errorf("MD5结果不匹配，期望: %s, 得到: %s", expected, result)
	}

	t.Logf("MD5 测试通过: %s", result)
}

func TestMD4(t *testing.T) {
	data := "Hello World"
	result := MD4(data)

	if result == "" {
		t.Error("MD4 结果为空")
	}

	if len(result) != 32 {
		t.Errorf("MD4 结果长度不正确，期望: 32, 得到: %d", len(result))
	}

	t.Logf("MD4 测试通过: %s", result)
}

func TestSHA1(t *testing.T) {
	data := "Hello World"
	expected := "0a4d55a8d778e5022fab701977c5d840bbc486d0" // SHA1 of "Hello World"

	result := SHA1(data)
	if result != expected {
		t.Errorf("SHA1结果不匹配，期望: %s, 得到: %s", expected, result)
	}

	t.Logf("SHA1 测试通过: %s", result)
}

func TestSHA256(t *testing.T) {
	data := "Hello World"
	expected := "a591a6d40bf420404a011733cfb7b190d62c65bf0bcda32b57b277d9ad9f146e" // SHA256 of "Hello World"

	result := SHA256(data)
	if result != expected {
		t.Errorf("SHA256结果不匹配，期望: %s, 得到: %s", expected, result)
	}

	t.Logf("SHA256 测试通过: %s", result)
}

func TestSHA512(t *testing.T) {
	data := "Hello World"
	result := SHA512(data)

	if result == "" {
		t.Error("SHA512 结果为空")
	}

	if len(result) != 128 {
		t.Errorf("SHA512 结果长度不正确，期望: 128, 得到: %d", len(result))
	}

	t.Logf("SHA512 测试通过: %s", result)
}

func TestSHA3_256(t *testing.T) {
	data := "Hello World"
	result := SHA3_256(data)

	if result == "" {
		t.Error("SHA3-256 结果为空")
	}

	if len(result) != 64 {
		t.Errorf("SHA3-256 结果长度不正确，期望: 64, 得到: %d", len(result))
	}

	t.Logf("SHA3-256 测试通过: %s", result)
}

func TestSHA3_512(t *testing.T) {
	data := "Hello World"
	result := SHA3_512(data)

	if result == "" {
		t.Error("SHA3-512 结果为空")
	}

	if len(result) != 128 {
		t.Errorf("SHA3-512 结果长度不正确，期望: 128, 得到: %d", len(result))
	}

	t.Logf("SHA3-512 测试通过: %s", result)
}

func TestAllSHAVariants(t *testing.T) {
	tests := []struct {
		name     string
		hashType HashType
		length   int
	}{
		{"SHA1", SHA1Hash, 40},
		{"SHA224", SHA224Hash, 56},
		{"SHA256", SHA256Hash, 64},
		{"SHA384", SHA384Hash, 96},
		{"SHA512", SHA512Hash, 128},
		{"SHA3-224", SHA3_224Hash, 56},
		{"SHA3-256", SHA3_256Hash, 64},
		{"SHA3-384", SHA3_384Hash, 96},
		{"SHA3-512", SHA3_512Hash, 128},
	}

	data := "Test all SHA variants"

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HashString(data, tt.hashType)

			if result == "" {
				t.Errorf("%s 结果为空", tt.name)
			}

			if len(result) != tt.length {
				t.Errorf("%s 结果长度不正确，期望: %d, 得到: %d", tt.name, tt.length, len(result))
			}

			t.Logf("%s: %s", tt.name, result)
		})
	}
}

func TestHashToFormat_Base64(t *testing.T) {
	data := "Hello World"

	result := HashToFormat(data, SHA256Hash, OutputBase64)
	if result == "" {
		t.Error("Base64格式哈希结果为空")
	}

	// Base64编码的SHA256应该是44个字符（32字节 * 4/3）
	if len(result) < 40 {
		t.Errorf("Base64格式哈希长度异常: %d", len(result))
	}

	t.Logf("SHA256 Base64格式: %s", result)
}

func TestHashToFormat_Hex(t *testing.T) {
	data := "Hello World"

	result := HashToFormat(data, SHA256Hash, OutputHex)
	if result == "" {
		t.Error("Hex格式哈希结果为空")
	}

	// Hex编码的SHA256应该是64个字符
	if len(result) != 64 {
		t.Errorf("Hex格式哈希长度不正确，期望: 64, 得到: %d", len(result))
	}

	t.Logf("SHA256 Hex格式: %s", result)
}

func TestHash_EmptyString(t *testing.T) {
	data := ""

	md5Result := MD5(data)
	if md5Result != "d41d8cd98f00b204e9800998ecf8427e" {
		t.Errorf("空字符串MD5不正确")
	}

	sha256Result := SHA256(data)
	if sha256Result != "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855" {
		t.Errorf("空字符串SHA256不正确")
	}

	t.Logf("空字符串哈希测试通过")
}

func TestHash_LongString(t *testing.T) {
	// 测试长字符串
	data := ""
	for i := 0; i < 1000; i++ {
		data += "Hello World "
	}

	md5Result := MD5(data)
	if md5Result == "" {
		t.Error("长字符串MD5结果为空")
	}

	sha256Result := SHA256(data)
	if sha256Result == "" {
		t.Error("长字符串SHA256结果为空")
	}

	t.Logf("长字符串哈希测试通过")
	t.Logf("数据长度: %d", len(data))
	t.Logf("MD5: %s", md5Result)
	t.Logf("SHA256: %s", sha256Result)
}

func TestHash_ConsistentResults(t *testing.T) {
	data := "Consistency test"

	// 多次哈希应该得到相同结果
	result1 := SHA256(data)
	result2 := SHA256(data)
	result3 := SHA256(data)

	if result1 != result2 || result2 != result3 {
		t.Errorf("哈希结果不一致: %s, %s, %s", result1, result2, result3)
	}

	t.Logf("一致性测试通过")
}

func TestHash_DifferentInputs(t *testing.T) {
	data1 := "Hello"
	data2 := "hello" // 小写

	result1 := SHA256(data1)
	result2 := SHA256(data2)

	if result1 == result2 {
		t.Error("不同输入得到了相同的哈希值")
	}

	t.Logf("不同输入测试通过")
	t.Logf("'Hello': %s", result1)
	t.Logf("'hello': %s", result2)
}

// 性能测试
func BenchmarkMD5(b *testing.B) {
	data := "Hello World, this is a benchmark test!"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MD5(data)
	}
}

func BenchmarkMD4(b *testing.B) {
	data := "Hello World, this is a benchmark test!"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MD4(data)
	}
}

func BenchmarkSHA1(b *testing.B) {
	data := "Hello World, this is a benchmark test!"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = SHA1(data)
	}
}

func BenchmarkSHA256(b *testing.B) {
	data := "Hello World, this is a benchmark test!"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = SHA256(data)
	}
}

func BenchmarkSHA512(b *testing.B) {
	data := "Hello World, this is a benchmark test!"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = SHA512(data)
	}
}

func BenchmarkSHA3_256(b *testing.B) {
	data := "Hello World, this is a benchmark test!"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = SHA3_256(data)
	}
}

func BenchmarkSHA3_512(b *testing.B) {
	data := "Hello World, this is a benchmark test!"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = SHA3_512(data)
	}
}
