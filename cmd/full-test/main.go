package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/pangge/baiduyun-encrypt-uploader/internal/baiduyun"
	"github.com/pangge/baiduyun-encrypt-uploader/internal/crypto"
)

func main() {
	fmt.Println("🔐 完整加密上传测试")
	fmt.Println("=====================================")

	// 测试文件路径
	localFilePath := "/root/.openclaw/workspace/xiage/GX011345.MP4"
	// 直接上传到根目录
	remoteBaseDir := ""

	// 你的BDUSS
	bduss := "E43eXFBZWhoUXdmQ0V6LThBc2hlRWZBMVJiMVZFQ1hNWVNWNX51aUpLfmFHTjVwRVFBQUFBJCQAAAAAAQAAAAEAAACxG5SjsNnIzLT-z7oAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAANqLtmnai7Zpd"

	// 1. 生成加密密钥
	fmt.Println("\n1️⃣ 生成加密密钥...")
	key, err := crypto.GenerateRandomKey()
	if err != nil {
		fmt.Printf("❌ 生成密钥失败: %v\n", err)
		return
	}
	fmt.Printf("✅ 密钥生成成功，请务必保存以下密钥！\n")
	fmt.Printf("   密钥(Base64): %s\n", crypto.EncodeKeyToBase64(key))
	fmt.Println("   ⚠️  丢失密钥将无法解密文件，请妥善保存！")

	// 2. 加密文件
	fmt.Println("\n2️⃣ 开始加密文件...")
	encryptedFilePath := "/tmp/GX011345.encrypted.bin"

	// 如果已经有加密文件，先删除
	os.Remove(encryptedFilePath)

	startTime := time.Now()
	header, err := crypto.EncryptFile(localFilePath, encryptedFilePath, key)
	if err != nil {
		fmt.Printf("❌ 加密失败: %v\n", err)
		return
	}
	elapsed := time.Since(startTime)
	fmt.Printf("✅ 加密完成！\n")
	fmt.Printf("   原文件大小: %d bytes (%.2f MB)\n", header.OriginalSize, float64(header.OriginalSize)/1024/1024)
	fmt.Printf("   加密后大小: %d bytes (%.2f MB)\n", header.EncryptedSize, float64(header.EncryptedSize)/1024/1024)
	fmt.Printf("   耗时: %v\n", elapsed)
	fmt.Printf("   加密速度: %.2f MB/s\n", float64(header.OriginalSize)/1024/1024/elapsed.Seconds())

	// 3. 创建百度云客户端
	fmt.Println("\n3️⃣ 初始化百度云客户端...")
	client := baiduyun.NewClientWithBDUSS(bduss)
	fmt.Println("✅ 百度云客户端初始化完成")

	// 4. 上传加密文件
	fmt.Println("\n4️⃣ 开始上传加密文件到百度云盘...")
	fileInfo, err := os.Stat(encryptedFilePath)
	if err != nil {
		fmt.Printf("❌ 获取加密文件信息失败: %v\n", err)
		return
	}

	var remotePath string
	if remoteBaseDir == "" {
		remotePath = filepath.Base(encryptedFilePath)
	} else {
		remotePath = filepath.Join(remoteBaseDir, filepath.Base(encryptedFilePath))
	}
	startTime = time.Now()
	resp, err := client.UploadFile(encryptedFilePath, remotePath)
	elapsed = time.Since(startTime)

	if err != nil {
		fmt.Printf("❌ 上传失败: %v\n", err)
		return
	}

	if resp.Errno != 0 {
		fmt.Printf("❌ 上传返回错误码: %d\n", resp.Errno)
		return
	}

	fmt.Printf("✅ 上传完成！\n")
	fmt.Printf("   云端路径: /%s\n", resp.Path)
	fmt.Printf("   文件大小: %d bytes\n", resp.Size)
	fmt.Printf("   耗时: %v\n", elapsed)
	fmt.Printf("   上传速度: %.2f MB/s\n", float64(fileInfo.Size())/1024/1024/elapsed.Seconds())

	fmt.Println("\n=====================================")
	fmt.Println("🎉 完整加密上传测试成功！")
	fmt.Println("\n📝 关键信息请保存：")
	fmt.Printf("   密钥: %s\n", crypto.EncodeKeyToBase64(key))
	fmt.Printf("   云端文件: /%s\n", resp.Path)
}
