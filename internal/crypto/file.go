package crypto

import (
	"encoding/binary"
	"io"
	"os"
)

const (
	// MagicNumber 加密文件的魔法数，用于识别文件格式
	MagicNumber = "BECU" // Baidu Encrypt Cloud Uploader
	// CurrentVersion 当前加密文件格式版本
	CurrentVersion = 1
	// HeaderSize 固定头部大小
	// 4字节magic + 1字节version + 12字节IV + 48字节加密的文件密钥 + 16字节内容标签 = 81字节
	HeaderSize = 4 + 1 + 12 + 48 + 16
)

// FileHeader 加密文件头部结构
type FileHeader struct {
	Magic          [4]byte  // 魔法数，固定为BECU
	Version        uint8    // 文件格式版本
	IV             [12]byte // AES-GCM的IV
	EncryptedKey   [48]byte // 主密钥加密后的文件密钥（32字节密钥 + 16字节GCM标签）
	ContentTag     [16]byte // 文件内容加密的GCM标签
	EncryptedSize  int64    // 加密后内容的大小
	OriginalSize   int64    // 原文件大小
	OriginalHash   [32]byte // 原文件SHA256哈希
}

// EncryptFile 加密整个文件
// srcPath: 原文件路径
// dstPath: 加密后文件路径
// masterKey: 主密钥，用于加密文件密钥
// 返回加密文件的头部信息和错误
func EncryptFile(srcPath, dstPath string, masterKey []byte) (*FileHeader, error) {
	// 1. 生成随机文件密钥
	fileKey, err := GenerateRandomKey()
	if err != nil {
		return nil, err
	}

	// 2. 生成随机IV
	iv, err := GenerateRandomIV()
	if err != nil {
		return nil, err
	}

	// 3. 用主密钥加密文件密钥
	masterAES, err := NewAESGCM(masterKey)
	if err != nil {
		return nil, err
	}
	encryptedKey, err := masterAES.Encrypt(fileKey, iv[:12]) // 用前12字节作为IV加密文件密钥
	if err != nil {
		return nil, err
	}

	// 4. 打开原文件
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return nil, err
	}
	defer srcFile.Close()

	// 5. 创建目标文件
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return nil, err
	}
	defer dstFile.Close()

	// 6. 先写入占位的头部（后续填充）
	header := &FileHeader{
		Version: CurrentVersion,
	}
	copy(header.Magic[:], MagicNumber)
	copy(header.IV[:], iv)
	copy(header.EncryptedKey[:], encryptedKey)

	// 写入空头部占位
	if err := binary.Write(dstFile, binary.LittleEndian, header); err != nil {
		return nil, err
	}

	// 7. 加密文件内容
	fileAES, err := NewAESGCM(fileKey)
	if err != nil {
		return nil, err
	}

	// 读取文件内容
	plaintext, err := io.ReadAll(srcFile)
	if err != nil {
		return nil, err
	}

	// 加密内容
	ciphertext, err := fileAES.Encrypt(plaintext, iv)
	if err != nil {
		return nil, err
	}

	// 分离内容和标签
	content := ciphertext[:len(ciphertext)-16]
	tag := ciphertext[len(ciphertext)-16:]
	copy(header.ContentTag[:], tag)

	// 写入加密内容
	if _, err := dstFile.Write(content); err != nil {
		return nil, err
	}

	// 填充头部信息
	header.OriginalSize = int64(len(plaintext))
	header.EncryptedSize = int64(len(content))

	// 回到文件开头，写入完整头部
	if _, err := dstFile.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	if err := binary.Write(dstFile, binary.LittleEndian, header); err != nil {
		return nil, err
	}

	return header, nil
}

// DecryptFile 解密整个文件
// srcPath: 加密文件路径
// dstPath: 解密后文件路径
// masterKey: 主密钥
// 返回错误
func DecryptFile(srcPath, dstPath string, masterKey []byte) error {
	// 1. 打开加密文件
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// 2. 读取头部
	var header FileHeader
	if err := binary.Read(srcFile, binary.LittleEndian, &header); err != nil {
		return err
	}

	// 3. 校验魔法数
	if string(header.Magic[:]) != MagicNumber {
		return ErrInvalidMagic
	}

	// 4. 校验版本
	if header.Version != CurrentVersion {
		return ErrUnsupportedVersion
	}

	// 5. 解密文件密钥
	masterAES, err := NewAESGCM(masterKey)
	if err != nil {
		return err
	}

	fileKey, err := masterAES.Decrypt(header.EncryptedKey[:], header.IV[:])
	if err != nil {
		return err
	}

	// 6. 读取加密内容
	ciphertext := make([]byte, header.EncryptedSize + 16) // 内容 + 标签
	if _, err := srcFile.Read(ciphertext[:header.EncryptedSize]); err != nil {
		return err
	}
	copy(ciphertext[header.EncryptedSize:], header.ContentTag[:])

	// 7. 解密内容
	fileAES, err := NewAESGCM(fileKey)
	if err != nil {
		return err
	}

	plaintext, err := fileAES.Decrypt(ciphertext, header.IV[:])
	if err != nil {
		return err
	}

	// 8. 写入解密后的文件
	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := dstFile.Write(plaintext); err != nil {
		return err
	}

	return nil
}

// ReadFileHeader 读取加密文件的头部信息
func ReadFileHeader(srcPath string) (*FileHeader, error) {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return nil, err
	}
	defer srcFile.Close()

	var header FileHeader
	if err := binary.Read(srcFile, binary.LittleEndian, &header); err != nil {
		return nil, err
	}

	if string(header.Magic[:]) != MagicNumber {
		return nil, ErrInvalidMagic
	}

	return &header, nil
}
