package main

import (
	"fmt"

	"github.com/pangge/baiduyun-encrypt-uploader/internal/baiduyun"
)

func main() {
	// 使用你提供的BDUSS直接测试
	bduss := "E43eXFBZWhoUXdmQ0V6LThBc2hlRWZBMVJiMVZFQ1hNWVNWNX51aUpLfmFHTjVwRVFBQUFBJCQAAAAAAQAAAAEAAACxG5SjsNnIzLT-z7oAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAANqLtmnai7Zpd"
	client := baiduyun.NewClientWithBDUSS(bduss)
	fmt.Println("☁️  测试BDUSS登录百度云盘...")

	userInfo, err := client.GetUserInfo()
	if err != nil {
		fmt.Printf("❌ 登录失败: %v\n", err)
		return
	}

	fmt.Printf("✅ 登录成功！\n")
	fmt.Printf("   用户名: %v\n", userInfo["baidu_name"])
	fmt.Printf("   用户ID: %v\n", userInfo["uk"])
	fmt.Println("🎉 BDUSS认证测试通过！")
}
