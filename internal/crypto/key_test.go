package crypto

import (
	"bytes"
	"testing"
)

// TestGenerateRandomKey 测试生成随机密钥
func TestGenerateRandomKey(t *testing.T) {
	// 生成第一个密钥
	key1, err := GenerateRandomKey()
	if err != nil {
		t.Fatalf("GenerateRandomKey failed: %v", err)
	}

	// 检查长度必须是32字节（AES-256）
	if len(key1) != 32 {
		t.Errorf("Expected key length 32, got %d", len(key1))
	}

	// 生成第二个密钥
	key2, err := GenerateRandomKey()
	if err != nil {
		t.Fatalf("Second GenerateRandomKey failed: %v", err)
	}

	// 两次生成结果必须不同（随机性测试）
	if bytes.Equal(key1, key2) {
		t.Errorf("Two generated keys are identical, expected different keys for randomness")
	}
}

// TestGenerateRandomIV 测试生成随机IV
func TestGenerateRandomIV(t *testing.T) {
	// 生成第一个IV
	iv1, err := GenerateRandomIV()
	if err != nil {
		t.Fatalf("GenerateRandomIV failed: %v", err)
	}

	// 检查长度必须是12字节（GCM推荐）
	if len(iv1) != 12 {
		t.Errorf("Expected IV length 12, got %d", len(iv1))
	}

	// 生成第二个IV
	iv2, err := GenerateRandomIV()
	if err != nil {
		t.Fatalf("Second GenerateRandomIV failed: %v", err)
	}

	// 两次生成结果必须不同
	if bytes.Equal(iv1, iv2) {
		t.Errorf("Two generated IVs are identical, expected different for randomness")
	}
}

// TestPasswordToKey 测试密码转换为密钥
func TestPasswordToKey(t *testing.T) {
	password1 := "my-secret-password-123"
	password2 := "another-different-password"

	// 相同密码应该生成相同密钥
	key1a := PasswordToKey(password1)
	key1b := PasswordToKey(password1)
	if !bytes.Equal(key1a, key1b) {
		t.Errorf("Same password should generate same key, got different results")
	}

	// 检查长度
	if len(key1a) != 32 {
		t.Errorf("Expected key length 32, got %d", len(key1a))
	}

	// 不同密码应该生成不同密钥
	key2 := PasswordToKey(password2)
	if bytes.Equal(key1a, key2) {
		t.Errorf("Different passwords should generate different keys, got same result")
	}
}

// TestEncryptDecryptKey 测试密钥加密解密
func TestEncryptDecryptKey(t *testing.T) {
	// 生成主密钥和文件密钥
	masterKey, _ := GenerateRandomKey()
	fileKey, _ := GenerateRandomKey()

	// 加密文件密钥
	encryptedKey, err := EncryptKey(fileKey, masterKey)
	if err != nil {
		t.Fatalf("EncryptKey failed: %v", err)
	}

	// 检查最小长度
	if len(encryptedKey) < 12+16 { // IV + tag
		t.Errorf("EncryptedKey too short: %d bytes", len(encryptedKey))
	}

	// 正确密钥解密
	decryptedKey, err := DecryptKey(encryptedKey, masterKey)
	if err != nil {
		t.Fatalf("DecryptKey with correct key failed: %v", err)
	}

	// 解密后必须和原文一致
	if !bytes.Equal(fileKey, decryptedKey) {
		t.Errorf("Decrypted key doesn't match original, got %x, expected %x", decryptedKey, fileKey)
	}

	// 错误密钥解密应该失败
	wrongMasterKey, _ := GenerateRandomKey()
	_, err = DecryptKey(encryptedKey, wrongMasterKey)
	if err == nil {
		t.Errorf("DecryptKey with wrong master key should fail, got success")
	}

	// 太短的密文应该失败
	shortKey := make([]byte, 10)
	_, err = DecryptKey(shortKey, masterKey)
	if err == nil {
		t.Errorf("DecryptKey with too short encryptedKey should fail, got success")
	}
}

// TestEncodeDecodeKeyBase64 测试密钥Base64编码解码
func TestEncodeDecodeKeyBase64(t *testing.T) {
	originalKey, _ := GenerateRandomKey()

	// 编码
	encoded := EncodeKeyToBase64(originalKey)

	// 解码
	decoded, err := DecodeKeyFromBase64(encoded)
	if err != nil {
		t.Fatalf("DecodeKeyFromBase64 failed: %v", err)
	}

	// 必须一致
	if !bytes.Equal(originalKey, decoded) {
		t.Errorf("Decoded key doesn't match original after base64 encode/decode")
	}
}
