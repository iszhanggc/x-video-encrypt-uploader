package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

// AESGCM 封装AES-256-GCM加密解密操作
type AESGCM struct {
	key []byte // 32字节密钥，对应AES-256
}

// NewAESGCM 创建AES-256-GCM实例
// key必须是32字节长度
func NewAESGCM(key []byte) (*AESGCM, error) {
	if len(key) != 32 {
		return nil, ErrInvalidKeyLength
	}
	return &AESGCM{key: key}, nil
}

// GenerateRandomKey 生成32字节随机AES-256密钥
func GenerateRandomKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// GenerateRandomIV 生成12字节随机IV（GCM推荐长度）
func GenerateRandomIV() ([]byte, error) {
	iv := make([]byte, 12)
	_, err := io.ReadFull(rand.Reader, iv)
	if err != nil {
		return nil, err
	}
	return iv, nil
}

// Encrypt 加密数据
// 返回加密后的数据（包含GCM标签）和可能的错误
func (a *AESGCM) Encrypt(plaintext []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 加密，结果格式：ciphertext + tag
	ciphertext := gcm.Seal(nil, iv, plaintext, nil)
	return ciphertext, nil
}

// Decrypt 解密数据
// ciphertext包含加密数据和GCM标签
func (a *AESGCM) Decrypt(ciphertext []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := gcm.Open(nil, iv, ciphertext, nil)
	if err != nil {
		return nil, ErrDecryptFailed
	}

	return plaintext, nil
}

// EncryptStream 流式加密，适合大文件
// reader: 明文输入流
// writer: 密文输出流
// iv: 12字节随机IV
// 返回加密后的总字节数和错误
func (a *AESGCM) EncryptStream(reader io.Reader, writer io.Writer, iv []byte) (int64, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return 0, err
	}

	// TODO: GCM流式加密实现，当前先用CTR模式
	stream := cipher.NewCTR(block, iv)
	writer = &cipher.StreamWriter{S: stream, W: writer}

	// 拷贝数据
	n, err := io.Copy(writer, reader)
	if err != nil {
		return n, err
	}

	return n, nil
}

// DecryptStream 流式解密，适合大文件
func (a *AESGCM) DecryptStream(reader io.Reader, writer io.Writer, iv []byte) (int64, error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return 0, err
	}

	// TODO: GCM流式解密实现，当前先用CTR模式
	stream := cipher.NewCTR(block, iv)
	writer = &cipher.StreamWriter{S: stream, W: writer}

	n, err := io.Copy(writer, reader)
	if err != nil {
		return n, err
	}

	return n, nil
}
