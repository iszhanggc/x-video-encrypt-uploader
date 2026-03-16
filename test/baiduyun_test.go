package main

import (
	"fmt"
	"os"

	"github.com/pangge/baiduyun-encrypt-uploader/internal/baiduyun"
)

func main() {
	fmt.Println("☁️  开始测试百度云盘客户端...")
	fmt.Println("=====================================")

	// 从环境变量获取AccessToken
	accessToken := os.Getenv("BAIDUYUN_ACCESS_TOKEN")
	if accessToken == "" {
		fmt.Println("⚠️  请设置环境变量 BAIDYUN_ACCESS_TOKEN 后再运行测试")
		fmt.Println("   export BAIDUYUN_ACCESS_TOKEN=your-token")
		return
	}

	// 创建客户端
	client := baiduyun.NewClientWithToken(accessToken)
	fmt.Println("✅ 创建百度云客户端成功")

	// 测试获取用户信息
	userInfo, err := client.GetUserInfo()
	if err != nil {
		fmt.Printf("❌ 获取用户信息失败: %v\n", err)
		return
	}
	fmt.Printf("✅ 获取用户信息成功: %v\n", userInfo["baidu_name"])	fmt.Println("=====================================")
	\n\tfmt.Println("=====================================")
	fmt.Println("🎉 百度云客户端测试通过！")
}

