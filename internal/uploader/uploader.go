package uploader

import (
	"fmt"
	"github.com/pangge/baiduyun-encrypt-uploader/internal/baiduyun"
	"path/filepath"
)

// Uploader 定义了上传器的统一接口
type Uploader interface {
	// Name 返回上传器名称
	Name() string
	// Upload 上传本地文件到云端，返回云端完整路径
	Upload(localPath string, size int64) (string, error)
	// Delete 删除云端文件
	Delete(remotePath string) error
}

// 确保 BaiduCloudUploader 实现了 Uploader 接口
var _ Uploader = (*BaiduCloudUploader)(nil)

// BaiduCloudUploader 百度云上传器实现
type BaiduCloudUploader struct {
	client  *baiduyun.Client
	baseDir string // 百度云基础目录
}

// NewBaiduCloudUploader 创建百度云上传器
func NewBaiduCloudUploader(baseDir string, accessToken string) *BaiduCloudUploader {
	return &BaiduCloudUploader{
		client:  baiduyun.NewClientWithToken(accessToken),
		baseDir: baseDir,
	}
}

// Name 返回上传器名称
func (b *BaiduCloudUploader) Name() string {
	return "baiduyun"
}

// Upload 实现Uploader接口
func (b *BaiduCloudUploader) Upload(localPath string, size int64) (string, error) {
	// 拼接完整的云端路径
	fileName := filepath.Base(localPath)
	remotePath := filepath.Join(b.baseDir, fileName)

	// 调用百度云上传
	resp, err := b.client.UploadFile(localPath, remotePath)
	if err != nil {
		return "", err
	}

	if resp.Errno != 0 {
		return "", fmt.Errorf("上传失败，百度云错误码: %d", resp.Errno)
	}

	return resp.Path, nil
}

// Delete 实现Uploader接口 - 删除功能后续再实现，先留个桩
func (b *BaiduCloudUploader) Delete(remotePath string) error {
	// TODO: 实现删除功能
	return fmt.Errorf("delete not implemented yet")
}
