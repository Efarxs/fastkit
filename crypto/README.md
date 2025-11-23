# Crypto - 通用加密库

一个功能完整的Go加密库，支持DES、TripleDES、AES、RSA以及多种哈希算法。

## 特性

- 🔐 **对称加密**
  - DES (支持5种模式和5种填充)
  - TripleDES (支持Twice和Triple)
  - AES (支持128/192/256位)
  
- 🔑 **非对称加密**
  - RSA (公私钥加解密)
  
- 🔨 **哈希算法**
  - MD5、MD4
  - SHA-1、SHA-2系列、SHA-3系列

- ✨ **灵活配置**
  - 5种加密模式: ECB、CBC、CTR、CFB、OFB
  - 5种填充模式: PKCS7、Zero、ISO10126、ANSIX923、NoPadding
  - 支持最多100位字符的KEY和IV
  - 多种输出格式: Base64、Hex、String

## 安装

```bash
go get github.com/efarxs/fastkit/crypto
```

## 快速开始

### DES加密

```go
package main

import (
    "fmt"
    "github.com/efarxs/fastkit/crypto"
)

func main() {
    // 创建DES加密器
    des, err := crypto.NewDES(
        "mykey123",           // 密钥
        "iv123456",           // IV向量
        crypto.CBC,           // 加密模式
        crypto.PKCS7,         // 填充模式
    )
    if err != nil {
        panic(err)
    }
    
    // 加密并输出为Base64
    encrypted, err := des.EncryptToString("Hello World", crypto.OutputBase64)
    if err != nil {
        panic(err)
    }
    fmt.Println("加密:", encrypted)
    
    // 解密
    decrypted, err := des.DecryptFromString(encrypted, crypto.OutputBase64)
    if err != nil {
        panic(err)
    }
    fmt.Println("解密:", decrypted)
}
```

### TripleDES加密

```go
// 创建3DES加密器（Triple模式，24字节密钥）
des3, err := crypto.NewTripleDES(
    "my24bytekey123456789abcd", // 密钥
    "iv123456",                  // IV向量
    crypto.CBC,                  // 加密模式
    crypto.PKCS7,                // 填充模式
    crypto.Triple,               // 3DES类型（Triple或Twice）
)

encrypted, _ := des3.EncryptToString("Hello World", crypto.OutputBase64)
decrypted, _ := des3.DecryptFromString(encrypted, crypto.OutputBase64)
```

### AES加密

```go
// 创建AES加密器（256位）
aes, err := crypto.NewAES(
    "my32bytekey12345678901234567890", // 密钥
    "iv1234567890abcd",                 // IV向量
    crypto.CBC,                         // 加密模式
    crypto.PKCS7,                       // 填充模式
    crypto.AES256,                      // 密钥长度（AES128/AES192/AES256）
)

// 加密
encrypted, _ := aes.EncryptToString("Hello World", crypto.OutputBase64)
// 解密
decrypted, _ := aes.DecryptFromString(encrypted, crypto.OutputBase64)

fmt.Println("AES加密:", encrypted)
fmt.Println("AES解密:", decrypted)
```

### RSA加密

```go
// 生成RSA密钥对
publicKeyPEM, privateKeyPEM, err := crypto.GenerateRSAKeyPair(2048)
if err != nil {
    panic(err)
}

// 创建RSA加密器
rsa, err := crypto.NewRSA(publicKeyPEM, privateKeyPEM)
if err != nil {
    panic(err)
}

// 公钥加密，私钥解密
encrypted, _ := rsa.EncryptToStringWithPublicKey("Hello World", crypto.OutputBase64)
decrypted, _ := rsa.DecryptFromStringWithPrivateKey(encrypted, crypto.OutputBase64)

fmt.Println("RSA加密:", encrypted)
fmt.Println("RSA解密:", decrypted)

// 私钥加密（签名），公钥验证
signature, _ := rsa.EncryptToStringWithPrivateKey("Hello World", crypto.OutputBase64)
fmt.Println("RSA签名:", signature)
```

### 哈希算法

```go
// MD5
md5Hash := crypto.MD5("Hello World")
fmt.Println("MD5:", md5Hash)

// MD4
md4Hash := crypto.MD4("Hello World")
fmt.Println("MD4:", md4Hash)

// SHA系列
sha1Hash := crypto.SHA1("Hello World")
sha256Hash := crypto.SHA256("Hello World")
sha512Hash := crypto.SHA512("Hello World")

fmt.Println("SHA1:", sha1Hash)
fmt.Println("SHA256:", sha256Hash)
fmt.Println("SHA512:", sha512Hash)

// SHA3系列
sha3_256Hash := crypto.SHA3_256("Hello World")
sha3_512Hash := crypto.SHA3_512("Hello World")

fmt.Println("SHA3-256:", sha3_256Hash)
fmt.Println("SHA3-512:", sha3_512Hash)

// 指定输出格式
hexHash := crypto.HashToFormat("Hello World", crypto.SHA256Hash, crypto.OutputHex)
base64Hash := crypto.HashToFormat("Hello World", crypto.SHA256Hash, crypto.OutputBase64)

fmt.Println("Hex格式:", hexHash)
fmt.Println("Base64格式:", base64Hash)
```

## 详细文档

### 加密模式 (EncryptMode)

| 模式 | 说明 | 需要IV |
|-----|------|--------|
| ECB | 电子密码本模式 | ❌ |
| CBC | 密码块链接模式 | ✅ |
| CTR | 计数器模式 | ✅ |
| CFB | 密码反馈模式 | ✅ |
| OFB | 输出反馈模式 | ✅ |

### 填充模式 (PaddingMode)

| 模式 | 说明 |
|-----|------|
| PKCS7 | PKCS#7填充（推荐） |
| Zero | 零字节填充 |
| ISO10126 | ISO10126填充（随机） |
| ANSIX923 | ANSI X9.23填充 |
| NoPadding | 不填充（数据必须是块大小的整数倍） |

### 输出格式 (OutputFormat)

| 格式 | 说明 |
|-----|------|
| OutputBase64 | Base64编码 |
| OutputHex | 十六进制编码 |
| OutputString | 原始字符串 |

### DES API

```go
// 创建DES加密器
des, err := crypto.NewDES(key, iv, mode, padding)

// 加密
ciphertext, err := des.Encrypt(plaintext []byte) ([]byte, error)
encrypted, err := des.EncryptToString(plaintext string, format OutputFormat) (string, error)

// 解密
plaintext, err := des.Decrypt(ciphertext []byte) ([]byte, error)
decrypted, err := des.DecryptFromString(ciphertext string, format OutputFormat) (string, error)
```

### TripleDES API

```go
// 创建3DES加密器
des3, err := crypto.NewTripleDES(key, iv, mode, padding, desType)
// desType: crypto.Twice (16字节) 或 crypto.Triple (24字节)

// 使用方法与DES相同
```

### AES API

```go
// 创建AES加密器
aes, err := crypto.NewAES(key, iv, mode, padding, keySize)
// keySize: crypto.AES128, crypto.AES192, crypto.AES256

// 使用方法与DES相同
```

### RSA API

```go
// 从PEM格式创建
rsa, err := crypto.NewRSA(publicKeyPEM, privateKeyPEM)

// 从Base64格式创建
rsa, err := crypto.NewRSAFromBase64(publicKeyBase64, privateKeyBase64)

// 生成密钥对
publicKeyPEM, privateKeyPEM, err := crypto.GenerateRSAKeyPair(bits)

// 公钥加密
ciphertext, err := rsa.EncryptWithPublicKey(plaintext []byte) ([]byte, error)
encrypted, err := rsa.EncryptToStringWithPublicKey(plaintext string, format) (string, error)

// 私钥解密
plaintext, err := rsa.DecryptWithPrivateKey(ciphertext []byte) ([]byte, error)
decrypted, err := rsa.DecryptFromStringWithPrivateKey(ciphertext string, format) (string, error)

// 私钥加密（签名）
signature, err := rsa.EncryptWithPrivateKey(plaintext []byte) ([]byte, error)
signature, err := rsa.EncryptToStringWithPrivateKey(plaintext string, format) (string, error)
```

### 哈希 API

```go
// 快捷方法
md5Hash := crypto.MD5(data string) string
md4Hash := crypto.MD4(data string) string
sha1Hash := crypto.SHA1(data string) string
sha256Hash := crypto.SHA256(data string) string
sha512Hash := crypto.SHA512(data string) string
sha3_256Hash := crypto.SHA3_256(data string) string
// ... 更多SHA变体

// 通用方法
hash := crypto.Hash(data []byte, hashType HashType) []byte
hashStr := crypto.HashString(data string, hashType HashType) string
hashFormatted := crypto.HashToFormat(data string, hashType HashType, format OutputFormat) string

// HashType支持:
// MD5Hash, MD4Hash, SHA1Hash, SHA224Hash, SHA256Hash, SHA384Hash, SHA512Hash,
// SHA3_224Hash, SHA3_256Hash, SHA3_384Hash, SHA3_512Hash
```

## 示例场景

### 场景1: 数据库密码加密存储

```go
// 使用AES-256-CBC加密
aes, _ := crypto.NewAES("your-secret-key-32bytes123456", "iv1234567890abcd", 
    crypto.CBC, crypto.PKCS7, crypto.AES256)

// 加密密码
encrypted, _ := aes.EncryptToString("user_password123", crypto.OutputBase64)
// 存储到数据库: encrypted

// 从数据库读取后解密
decrypted, _ := aes.DecryptFromString(encrypted, crypto.OutputBase64)
```

### 场景2: API签名验证

```go
// 使用HMAC-SHA256（这里用SHA256示例）
data := "api_key=xxx&timestamp=1234567890"
signature := crypto.SHA256(data + "secret_key")

// 验证
receivedSignature := "..." // 从请求中获取
if signature == receivedSignature {
    // 签名有效
}
```

### 场景3: 文件加密传输

```go
// 读取文件
fileData, _ := os.ReadFile("document.pdf")

// AES加密
aes, _ := crypto.NewAES("key", "iv", crypto.CBC, crypto.PKCS7, crypto.AES256)
encrypted, _ := aes.Encrypt(fileData)

// 保存加密文件
os.WriteFile("document.pdf.enc", encrypted, 0644)

// 解密
decrypted, _ := aes.Decrypt(encrypted)
os.WriteFile("document_decrypted.pdf", decrypted, 0644)
```

### 场景4: 混合加密（RSA+AES）

```go
// 1. 生成随机AES密钥
aesKey := "random-aes-key-32bytes-12345678"

// 2. 使用RSA加密AES密钥
rsa, _ := crypto.NewRSA(publicKeyPEM, "")
encryptedKey, _ := rsa.EncryptToStringWithPublicKey(aesKey, crypto.OutputBase64)

// 3. 使用AES加密大数据
aes, _ := crypto.NewAES(aesKey, "iv1234567890abcd", crypto.CBC, crypto.PKCS7, crypto.AES256)
encryptedData, _ := aes.EncryptToString("large data...", crypto.OutputBase64)

// 4. 传输: encryptedKey + encryptedData

// 解密过程相反
```

## 注意事项

1. **密钥管理**: 请妥善保管密钥，不要硬编码在代码中
2. **IV向量**: CBC、CTR、CFB、OFB模式需要IV，且应该随机生成
3. **ECB模式**: 不推荐使用ECB模式，因为安全性较低
4. **填充**: 使用流模式（CTR、CFB、OFB）时通常不需要填充
5. **RSA**: 加密数据大小不能超过密钥长度-11字节
6. **密钥长度**: 虽然支持最多100字符，但实际会根据算法要求自动调整

## 性能建议

- **小数据**: 直接使用对称加密（AES推荐）
- **大数据**: 使用混合加密（RSA加密密钥 + AES加密数据）
- **哈希**: SHA256是速度和安全性的良好平衡
- **密钥派生**: 使用PBKDF2或bcrypt从密码派生密钥

## License

MIT
