package crypto

import (
	"bytes"
	"testing"
)

// TestNewAESGCM 测试创建AES实例
func TestNewAESGCM(t *testing.T) {
	// 32字节密钥应该成功
	validKey := make([]byte, 32)
	for i := 0; i < 32; i++ {
		validKey[i] = byte(i)
	}
	aes, err := NewAESGCM(validKey)
	if err != nil {
		t.Fatalf("NewAESGCM with 32-byte key failed: %v", err)
	}
	if aes == nil {
		t.Fatalf("NewAESGCM returned nil AES")
	}

	// 错误长度密钥应该失败
	shortKey := make([]byte, 16)
	_, err = NewAESGCM(shortKey)
	if err != ErrInvalidKeyLength {
		t.Errorf("Expected ErrInvalidKeyLength for 16-byte key, got %v", err)
	}

	longKey := make([]byte, 64)
	_, err = NewAESGCM(longKey)
	if err != ErrInvalidKeyLength {
		t.Errorf("Expected ErrInvalidKeyLength for 64-byte key, got %v", err)
	}
}

// TestEncryptDecrypt 测试基础加密解密
func TestEncryptDecrypt(t *testing.T) {
	key, _ := GenerateRandomKey()
	iv, _ := GenerateRandomIV()
	aes, _ := NewAESGCM(key)

	// 测试普通明文
	plaintext := []byte("Hello, World! 这是中文测试 123!@#")
	ciphertext, err := aes.Encrypt(plaintext, iv)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	// 密文长度应该是明文 + 16字节tag
	if len(ciphertext) != len(plaintext)+16 {
		t.Errorf("Expected ciphertext length %d, got %d", len(plaintext)+16, len(ciphertext))
	}

	// 正确解密
	decrypted, err := aes.Decrypt(ciphertext, iv)
	if err != nil {
		t.Fatalf("Decrypt with correct key/iv failed: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Errorf("Decrypted text doesn't match original\noriginal: %x\ndecrypted: %x", plaintext, decrypted)
	}

	// 测试空明文
	emptyPlaintext := []byte{}
	ciphertextEmpty, err := aes.Encrypt(emptyPlaintext, iv)
	if err != nil {
		t.Fatalf("Encrypt empty plaintext failed: %v", err)
	}
	decryptedEmpty, err := aes.Decrypt(ciphertextEmpty, iv)
	if err != nil {
		t.Fatalf("Decrypt empty plaintext failed: %v", err)
	}
	if !bytes.Equal(emptyPlaintext, decryptedEmpty) {
		t.Errorf("Empty plaintext decryption failed")
	}
}

// TestEncryptDecrypt_Tampered 测试密文被篡改后解密失败
func TestEncryptDecrypt_Tampered(t *testing.T) {
	key, _ := GenerateRandomKey()
	iv, _ := GenerateRandomIV()
	aes, _ := NewAESGCM(key)

	plaintext := []byte("Original content that shouldn't be tampered")
	ciphertext, _ := aes.Encrypt(plaintext, iv)

	// 篡改一个字节
	tampered := make([]byte, len(ciphertext))
	copy(tampered, ciphertext)
	tampered[len(tampered)/2] ^= 0xFF // 翻转一位

	_, err := aes.Decrypt(tampered, iv)
	if err == nil {
		t.Errorf("Decrypt of tampered ciphertext should fail, got success")
	}
	if err != ErrDecryptFailed {
		t.Errorf("Expected ErrDecryptFailed for tampered ciphertext, got %v", err)
	}
}

// TestEncryptDecrypt_WrongIV 测试错误IV解密失败
func TestEncryptDecrypt_WrongIV(t *testing.T) {
	key, _ := GenerateRandomKey()
	iv, _ := GenerateRandomIV()
	wrongIV, _ := GenerateRandomIV()
	aes, _ := NewAESGCM(key)

	plaintext := []byte("Test with wrong IV")
	ciphertext, _ := aes.Encrypt(plaintext, iv)

	_, err := aes.Decrypt(ciphertext, wrongIV)
	if err == nil {
		t.Errorf("Decrypt with wrong IV should fail, got success")
	}
}

// TestEncryptDecrypt_WrongKey 测试错误密钥解密失败
func TestEncryptDecrypt_WrongKey(t *testing.T) {
	rightKey, _ := GenerateRandomKey()
	wrongKey, _ := GenerateRandomKey()
	iv, _ := GenerateRandomIV()
	rightAES, _ := NewAESGCM(rightKey)
	wrongAES, _ := NewAESGCM(wrongKey)

	plaintext := []byte("Test with wrong key")
	ciphertext, _ := rightAES.Encrypt(plaintext, iv)

	_, err := wrongAES.Decrypt(ciphertext, iv)
	if err == nil {
		t.Errorf("Decrypt with wrong key should fail, got success")
	}
	if err != ErrDecryptFailed {
		t.Errorf("Expected ErrDecryptFailed for wrong key, got %v", err)
	}
}

// TestEncryptStreamDecryptStream 测试流式加密解密
func TestEncryptStreamDecryptStream(t *testing.T) {
	key, _ := GenerateRandomKey()
	iv, _ := GenerateRandomIV()
	aes, _ := NewAESGCM(key)

	// 测试小块数据
	smallData := []byte("This is a small block of data for streaming test")
	var encrypted bytes.Buffer
	srcReader := bytes.NewReader(smallData)

	_, err := aes.EncryptStream(srcReader, &encrypted, iv)
	if err != nil {
		t.Fatalf("EncryptStream small data failed: %v", err)
	}

	var decrypted bytes.Buffer
	encryptedReader := bytes.NewReader(encrypted.Bytes())
	_, err = aes.DecryptStream(encryptedReader, &decrypted, iv)
	if err != nil {
		t.Fatalf("DecryptStream small data failed: %v", err)
	}

	if !bytes.Equal(smallData, decrypted.Bytes()) {
		t.Errorf("Streaming decryption doesn't match original for small data")
	}

	// 测试大块数据（1MB）
	largeData := make([]byte, 1024*1024)
	for i := 0; i < len(largeData); i++ {
		largeData[i] = byte(i % 256)
	}

	var encryptedLarge bytes.Buffer
	srcLargeReader := bytes.NewReader(largeData)
	_, err = aes.EncryptStream(srcLargeReader, &encryptedLarge, iv)
	if err != nil {
		t.Fatalf("EncryptStream 1MB data failed: %v", err)
	}

	var decryptedLarge bytes.Buffer
	encryptedLargeReader := bytes.NewReader(encryptedLarge.Bytes())
	_, err = aes.DecryptStream(encryptedLargeReader, &decryptedLarge, iv)
	if err != nil {
		t.Fatalf("DecryptStream 1MB data failed: %v", err)
	}

	if !bytes.Equal(largeData, decryptedLarge.Bytes()) {
		t.Errorf("Streaming decryption doesn't match original for 1MB data")
	}
}
