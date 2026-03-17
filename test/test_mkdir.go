package main

import (
	"fmt"

	"github.com/pangge/baiduyun-encrypt-uploader/internal/baiduyun"
)

func main() {
	bduss := "E43eXFBZWhoUXdmQ0V6LThBc2hlRWZBMVJiMVZFQ1hNWVNWNX51aUpLfmFHTjVwRVFBQUFBJCQAAAAAAQAAAAEAAACxG5SjsNnIzLT-z7oAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAANqLtmnai7Zpd"
	client := baiduyun.NewClientWithBDUSS(bduss)

	fmt.Println("🧪 测试创建目录 /apps")
	err := client.EnsureBaseDir("/apps")
	if err != nil {
		fmt.Printf("❌ 创建目录失败: %v\n", err)
		return
	}
	fmt.Println("✅ 创建 /apps 成功")

	fmt.Println("\n🧪 测试创建目录 /apps/x-video-encrypt-uploader")
	err = client.EnsureBaseDir("/apps/x-video-encrypt-uploader")
	if err != nil {
		fmt.Printf("❌ 创建目录失败: %v\n", err)
		return
	}
	fmt.Println("✅ 创建 /apps/x-video-encrypt-uploader 成功")

	fmt.Println("\n🎉 所有目录创建测试通过！")
}
