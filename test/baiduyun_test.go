package main

import (
	"fmt"
	"os"

	"github.com/pangge/baiduyun-encrypt-uploader/internal/baiduyun"
)

func main() {
	fmt.Println("☁️  开始测试百度云盘客户端...")
	fmt.Println("=====================================")

	// 优先尝试BDUSS方式（推荐）
	bduss := os.Getenv("BAIDUYUN_BDUSS")
	if bduss != "" {
		// 使用BDUSS创建客户端
		client := baiduyun.NewClientWithBDUSS(bduss)
		fmt.Println("✅ 使用BDUSS创建百度云客户端成功")

		// 测试获取用户信息
		userInfo, err := client.GetUserInfo()
		if err != nil {
			fmt.Printf("❌ 获取用户信息失败: %v\n", err)
			return
		}
		fmt.Printf("✅ 获取用户信息成功: %v\n", userInfo["baidu_name"])
		fmt.Println("=====================================")
		fmt.Println("🎉 百度云客户端(BDUSS方式)测试通过！")
		return
	}

	// 回退到AccessToken方式
	accessToken := os.Getenv("BAIDUYUN_ACCESS_TOKEN")
	if accessToken != "" {
		// 创建客户端
		client := baiduyun.NewClientWithToken(accessToken)
		fmt.Println("✅ 使用AccessToken创建百度云客户端成功")

		// 测试获取用户信息
		userInfo, err := client.GetUserInfo()
		if err != nil {
			fmt.Printf("❌ 获取用户信息失败: %v\n", err)
			return
		}
		fmt.Printf("✅ 获取用户信息成功: %v\n", userInfo["baidu_name"])
		fmt.Println("=====================================")
		fmt.Println("🎉 百度云客户端(AccessToken方式)测试通过！")
		return
	}

	// 都没设置
	fmt.Println("⚠️  请设置环境变量后再运行测试:")
	fmt.Println("   推荐BDUSS方式: export BAIDUYUN_BDUSS=your-bduss")
	fmt.Println("   或者AccessToken方式: export BAIDUYUN_ACCESS_TOKEN=your-token")
}
