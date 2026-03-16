package crypto

import "errors"

var (
	// ErrInvalidKeyLength 密钥长度错误，必须为32字节(AES-256)
	ErrInvalidKeyLength = errors.New("invalid key length, must be 32 bytes for AES-256")
	// ErrInvalidIVLength IV长度错误，必须为12字节(GCM推荐长度)
	ErrInvalidIVLength = errors.New("invalid IV length, must be 12 bytes for GCM")
	// ErrDecryptFailed 解密失败，可能是密钥错误或数据损坏
	ErrDecryptFailed = errors.New("decryption failed, invalid key or corrupted data")
	// ErrInvalidMagic 无效的文件魔法数，不是我们的加密文件格式
	ErrInvalidMagic = errors.New("invalid file magic, not a encrypted file")
	// ErrUnsupportedVersion 不支持的文件版本
	ErrUnsupportedVersion = errors.New("unsupported file version")
)
