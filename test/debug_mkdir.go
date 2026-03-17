package main

import (
	"fmt"

	"github.com/pangge/baiduyun-encrypt-uploader/internal/baiduyun"
)

func main() {
	bduss := "E43eXFBZWhoUXdmQ0V6LThBc2hlRWZBMVJiMVZFQ1hNWVNWNX51aUpLfmFHTjVwRVFBQUFBJCQAAAAAAQAAAAEAAACxG5SjsNnIzLT-z7oAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAANqLtmnai7Zpd"
	client := baiduyun.NewClientWithBDUSS(bduss)

	fmt.Println("Testing mkdir /apps")
	err := client.EnsureBaseDir("/apps")
	if err != nil {
		fmt.Printf("ERROR /apps: %v\n", err)
	} else {
		fmt.Println("OK /apps")
	}

	fmt.Println("\nTesting mkdir /apps/x-video-encrypt-uploader")
	err = client.EnsureBaseDir("/apps/x-video-encrypt-uploader")
	if err != nil {
		fmt.Printf("ERROR /apps/x-video-encrypt-uploader: %v\n", err)
	} else {
		fmt.Println("OK /apps/x-video-encrypt-uploader")
	}
}
