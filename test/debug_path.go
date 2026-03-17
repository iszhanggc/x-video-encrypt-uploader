package main

import (
	"fmt"
	"path/filepath"

	"github.com/pangge/baiduyun-encrypt-uploader/internal/baiduyun"
)

func main() {
	bduss := "E43eXFBZWhoUXdmQ0V6LThBc2hlRWZBMVJiMVZFQ1hNWVNWNX51aUpLfmFHTjVwRVFBQUFBJCQAAAAAAQAAAAEAAACxG5SjsNnIzLT-z7oAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAANqLtmnai7Zpd"
	client := baiduyun.NewClientWithBDUSS(bduss)

	remotePath := "apps/x-video-encrypt-uploader/GX011345.encrypted.bin"
	dirPath := filepath.Dir(remotePath)
	fmt.Printf("remotePath = %q\n", remotePath)
	fmt.Printf("dirPath = %q\n", dirPath)

	err := client.EnsureBaseDir(dirPath)
	if err != nil {
		fmt.Printf("EnsureBaseDir failed: %v\n", err)
	} else {
		fmt.Println("EnsureBaseDir success!")
	}
}
