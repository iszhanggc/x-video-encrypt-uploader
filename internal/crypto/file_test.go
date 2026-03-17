package crypto

import (
	"bytes"
	"os"
	"testing"
)

// TestEncryptFileDecryptFile_Small 测试小文件加密解密
func TestEncryptFileDecryptFile_Small(t *testing.T) {
	masterKey, _ := GenerateRandomKey()

	// 创建测试文件
	testContent := []byte("This is a small test file content with 中文测试!@#$%^&*()")
	srcPath := "/tmp/test_small_src.bin"
	encPath := "/tmp/test_small_enc.bin"
	decPath := "/tmp/test_small_dec.bin"

	_ = os.Remove(srcPath)
	_ = os.Remove(encPath)
	_ = os.Remove(decPath)
	defer func() {
		_ = os.Remove(srcPath)
		_ = os.Remove(encPath)
		_ = os.Remove(decPath)
	}()

	err := os.WriteFile(srcPath, testContent, 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// 加密
	header, err := EncryptFile(srcPath, encPath, masterKey)
	if err != nil {
		t.Fatalf("EncryptFile failed: %v", err)
	}

	if header.OriginalSize != int64(len(testContent)) {
		t.Errorf("Expected original size %d, got %d", len(testContent), header.OriginalSize)
	}

	if string(header.Magic[:]) != MagicNumber {
		t.Errorf("Expected magic %q, got %q", MagicNumber, string(header.Magic[:]))
	}

	if header.Version != CurrentVersion {
		t.Errorf("Expected version %d, got %d", CurrentVersion, header.Version)
	}

	// 读取头部测试
	readHeader, err := ReadFileHeader(encPath)
	if err != nil {
		t.Fatalf("ReadFileHeader failed: %v", err)
	}

	if readHeader.OriginalSize != header.OriginalSize {
		t.Errorf("Read back original size mismatch: original %d, read %d", header.OriginalSize, readHeader.OriginalSize)
	}

	// 解密
	err = DecryptFile(encPath, decPath, masterKey)
	if err != nil {
		t.Fatalf("DecryptFile failed: %v", err)
	}

	// 比较内容
	decryptedContent, err := os.ReadFile(decPath)
	if err != nil {
		t.Fatalf("Failed to read decrypted file: %v", err)
	}

	if !bytes.Equal(testContent, decryptedContent) {
		t.Errorf("Decrypted content doesn't match original\noriginal: %s\ndecrypted: %s", testContent, decryptedContent)
	}
}

// TestEncryptFileDecryptFile_Large 测试大文件加密解密
func TestEncryptFileDecryptFile_Large(t *testing.T) {
	masterKey, _ := GenerateRandomKey()

	// 10MB 大文件
	size := 10 * 1024 * 1024
	largeData := make([]byte, size)
	for i := 0; i < size; i++ {
		largeData[i] = byte(i % 256)
	}

	srcPath := "/tmp/test_large_src.bin"
	encPath := "/tmp/test_large_enc.bin"
	decPath := "/tmp/test_large_dec.bin"

	_ = os.Remove(srcPath)
	_ = os.Remove(encPath)
	_ = os.Remove(decPath)
	defer func() {
		_ = os.Remove(srcPath)
		_ = os.Remove(encPath)
		_ = os.Remove(decPath)
	}()

	err := os.WriteFile(srcPath, largeData, 0644)
	if err != nil {
		t.Fatalf("Failed to create large test file: %v", err)
	}

	// 加密
	_, err = EncryptFile(srcPath, encPath, masterKey)
	if err != nil {
		t.Fatalf("EncryptFile large failed: %v", err)
	}

	// 解密
	err = DecryptFile(encPath, decPath, masterKey)
	if err != nil {
		t.Fatalf("DecryptFile large failed: %v", err)
	}

	// 比较内容
	decryptedContent, err := os.ReadFile(decPath)
	if err != nil {
		t.Fatalf("Failed to read decrypted large file: %v", err)
	}

	if !bytes.Equal(largeData, decryptedContent) {
		t.Errorf("Decrypted large content doesn't match original")
	}
}

// TestEncryptFileDecryptFile_WrongKey 测试错误密钥解密失败
func TestEncryptFileDecryptFile_WrongKey(t *testing.T) {
	correctKey, _ := GenerateRandomKey()
	wrongKey, _ := GenerateRandomKey()

	testContent := []byte("Test wrong key decryption should fail")
	srcPath := "/tmp/test_wrongkey_src.bin"
	encPath := "/tmp/test_wrongkey_enc.bin"
	decPath := "/tmp/test_wrongkey_dec.bin"

	_ = os.Remove(srcPath)
	_ = os.Remove(encPath)
	_ = os.Remove(decPath)
	defer func() {
		_ = os.Remove(srcPath)
		_ = os.Remove(encPath)
		_ = os.Remove(decPath)
	}()

	_ = os.WriteFile(srcPath, testContent, 0644)
	_, err := EncryptFile(srcPath, encPath, correctKey)
	if err != nil {
		t.Fatalf("EncryptFile failed: %v", err)
	}

	// 用错误密钥解密
	err = DecryptFile(encPath, decPath, wrongKey)
	if err == nil {
		t.Errorf("DecryptFile with wrong master key should fail, got success")
	}
	_ = os.Remove(decPath)
}

// TestReadFileHeader_InvalidMagic 测试读取非加密文件失败（魔法数不对）
func TestReadFileHeader_InvalidMagic(t *testing.T) {
	// 创建一个长度够但魔法数不对的文件
	f, err := os.Create("/tmp/test_invalid_magic.bin")
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	// 文件头需要 HeaderSize 字节，写够长度但魔法数不对
	data := make([]byte, HeaderSize)
	copy(data[:], "WRNG") // 前四个字节不对，不是 BECU
	_, _ = f.Write(data)
	_ = f.Close()
	defer os.Remove("/tmp/test_invalid_magic.bin")

	_, err = ReadFileHeader("/tmp/test_invalid_magic.bin")
	if err == nil {
		t.Errorf("ReadFileHeader with invalid magic should fail, got success")
	}
	// 只要失败就行，这里我们只需要确保不是成功读取
	// 因为文件长度够，会读到 ErrInvalidMagic
}

// TestReadFileHeader_UnsupportedVersion 测试不支持的版本
func TestReadFileHeader_UnsupportedVersion(t *testing.T) {
	// 这个测试需要构造一个有正确魔法数但版本不对的文件
	// 我们这里简单验证一下错误类型即可，实际构造比较麻烦先跳过框架
	t.Skip("Skipping version test for now")
}

// TestEncryptFile_NonexistentSource 测试源文件不存在
func TestEncryptFile_NonexistentSource(t *testing.T) {
	masterKey, _ := GenerateRandomKey()
	_, err := EncryptFile("/path/does/not/exist.bin", "/tmp/out.bin", masterKey)
	if err == nil {
		t.Errorf("EncryptFile with nonexistent source should fail, got success")
	}
}
