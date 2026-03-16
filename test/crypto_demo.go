package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"os"

	"github.com/pangge/baiduyun-encrypt-uploader/internal/crypto"
)

func main() {
	fmt.Println("🔐 开始测试加密模块...")
	fmt.Println("=====================================")

	// 1. 测试基础加密解密
	fmt.Println("\n1. 测试基础AES加密解密:")
	key, err := crypto.GenerateRandomKey()
	if err != nil {
		fmt.Printf("❌ 生成密钥失败: %v\n", err)
		return
	}
	fmt.Printf("✅ 生成32字节AES-256密钥: %x...\n", key[:8])

	iv, err := crypto.GenerateRandomIV()
	if err != nil {
		fmt.Printf("❌ 生成IV失败: %v\n", err)
		return
	}
	fmt.Printf("✅ 生成12字节IV: %x...\n", iv[:8])

	aes, err := crypto.NewAESGCM(key)
	if err != nil {
		fmt.Printf("❌ 创建AES实例失败: %v\n", err)
		return
	}

	plaintext := []byte("这是测试加密解密的文本内容123456!@#$%^&*()")
	ciphertext, err := aes.Encrypt(plaintext, iv)
	if err != nil {
		fmt.Printf("❌ 加密失败: %v\n", err)
		return
	}
	fmt.Printf("✅ 加密成功，密文长度: %d字节\n", len(ciphertext))

	decrypted, err := aes.Decrypt(ciphertext, iv)
	if err != nil {
		fmt.Printf("❌ 解密失败: %v\n", err)
		return
	}

	if !bytes.Equal(plaintext, decrypted) {
		fmt.Printf("❌ 解密结果不匹配: 原文=%s, 解密后=%s\n", plaintext, decrypted)
		return
	}
	fmt.Println("✅ 解密成功，内容完全一致!")

	// 2. 测试错误密钥解密
	fmt.Println("\n2. 测试错误密钥解密:")
	wrongKey, _ := crypto.GenerateRandomKey()
	wrongAes, _ := crypto.NewAESGCM(wrongKey)
	_, err = wrongAes.Decrypt(ciphertext, iv)
	if err == nil {
		fmt.Println("❌ 错误密钥应该解密失败，但实际成功了!")
		return
	}
	fmt.Printf("✅ 错误密钥正确返回失败: %v\n", err)

	// 3. 测试文件加密解密
	fmt.Println("\n3. 测试文件加密解密:")
	testContent := []byte("这是测试文件加密的内容，包含各种字符：!@#$%^&*()_+，测试中文也没问题")
	srcPath := "/tmp/test_src.txt"
	encryptedPath := "/tmp/test_encrypted.bin"
	decryptedPath := "/tmp/test_decrypted.txt"

	err = os.WriteFile(srcPath, testContent, 0644)
	if err != nil {
		fmt.Printf("❌ 写入测试文件失败: %v\n", err)
		return
	}
	defer os.Remove(srcPath)
	defer os.Remove(encryptedPath)
	defer os.Remove(decryptedPath)

	header, err := crypto.EncryptFile(srcPath, encryptedPath, key)
	if err != nil {
		fmt.Printf("❌ 加密文件失败: %v\n", err)
		return
	}
	fmt.Printf("✅ 文件加密成功，原大小=%d, 加密后大小=%d\n", header.OriginalSize, header.EncryptedSize)

	// 读取文件头部
	readHeader, err := crypto.ReadFileHeader(encryptedPath)
	if err != nil {
		fmt.Printf("❌ 读取文件头部失败: %v\n", err)
		return
	}
	if string(readHeader.Magic[:]) != crypto.MagicNumber {
		fmt.Printf("❌ 文件魔法数不匹配: 期望=%s, 实际=%s\n", crypto.MagicNumber, readHeader.Magic)
		return
	}
	if readHeader.Version != crypto.CurrentVersion {
		fmt.Printf("❌ 文件版本不匹配: 期望=%d, 实际=%d\n", crypto.CurrentVersion, readHeader.Version)
		return
	}
	fmt.Println("✅ 文件头部校验通过")

	// 解密文件
	err = crypto.DecryptFile(encryptedPath, decryptedPath, key)
	if err != nil {
		fmt.Printf("❌ 解密文件失败: %v\n", err)
		return
	}

	decryptedContent, err := os.ReadFile(decryptedPath)
	if err != nil {
		fmt.Printf("❌ 读取解密文件失败: %v\n", err)
		return
	}

	if !bytes.Equal(testContent, decryptedContent) {
		fmt.Println("❌ 解密文件内容不匹配!")
		return
	}
	fmt.Println("✅ 文件解密成功，内容完全一致!")

	// 4. 测试大文件加密解密
	fmt.Println("\n4. 测试10MB大文件加密解密:")
	largeData := make([]byte, 10*1024*1024) // 10MB
	_, err = rand.Read(largeData)
	if err != nil {
		fmt.Printf("❌ 生成随机数据失败: %v\n", err)
		return
	}

	largeSrcPath := "/tmp/test_large_src.bin"
	largeEncryptedPath := "/tmp/test_large_encrypted.bin"
	largeDecryptedPath := "/tmp/test_large_decrypted.bin"

	err = os.WriteFile(largeSrcPath, largeData, 0644)
	if err != nil {
		fmt.Printf("❌ 写入大文件失败: %v\n", err)
		return
	}
	defer os.Remove(largeSrcPath)
	defer os.Remove(largeEncryptedPath)
	defer os.Remove(largeDecryptedPath)

	_, err = crypto.EncryptFile(largeSrcPath, largeEncryptedPath, key)
	if err != nil {
		fmt.Printf("❌ 加密大文件失败: %v\n", err)
		return
	}
	fmt.Println("✅ 10MB大文件加密成功")

	err = crypto.DecryptFile(largeEncryptedPath, largeDecryptedPath, key)
	if err != nil {
		fmt.Printf("❌ 解密大文件失败: %v\n", err)
		return
	}

	decryptedLargeData, err := os.ReadFile(largeDecryptedPath)
	if err != nil {
		fmt.Printf("❌ 读取解密大文件失败: %v\n", err)
		return
	}

	if !bytes.Equal(largeData, decryptedLargeData) {
		fmt.Println("❌ 大文件解密内容不匹配!")
		return
	}
	fmt.Println("✅ 10MB大文件解密成功，内容完全一致!")

	fmt.Println("\n=====================================")
	fmt.Println("🎉 所有加密模块测试全部通过！")
}
