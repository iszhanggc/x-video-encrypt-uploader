package main

import (
	"fmt"

	"github.com/pangge/baiduyun-encrypt-uploader/internal/baiduyun"
)

func main() {
	bduss := "E43eXFBZWhoUXdmQ0V6LThBc2hlRWZBMVJiMVZFQ1hNWVNWNX51aUpLfmFHTjVwRVFBQUFBJCQAAAAAAQAAAAEAAACxG5SjsNnIzLT-z7oAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAANqLtmnai7Zpd"
	client := baiduyun.NewClientWithBDUSS(bduss)

	remotePath := "/apps/x-video-encrypt-uploader/GX011345.encrypted.bin"
	fileSize := int64(723342661)
	blockMd5s := []string{}

	fmt.Printf("DEBUG: 测试预上传: %s, size=%d\n", remotePath, fileSize)
	resp, err := client.Preupload(remotePath, fileSize, blockMd5s, true)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	fmt.Printf("SUCCESS: resp=%+v\n", *resp)
}
