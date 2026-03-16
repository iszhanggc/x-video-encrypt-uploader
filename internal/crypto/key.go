package crypto

import (
	"crypto/sha256"
	"errors"
)

// 注：助记词功能需要github.com/tyler-smith/go-bip39依赖，待网络恢复后添加

// PasswordToKey 用户密码转换为32字节AES密钥
func PasswordToKey(password string) []byte {
	// 使用SHA-256对密码哈希，得到32字节密钥
	hash := sha256.Sum256([]byte(password))
	return hash[:]
}

// EncryptKey 用主密钥加密文件密钥
func EncryptKey(fileKey []byte, masterKey []byte) ([]byte, error) {
	aes, err := NewAESGCM(masterKey)
	if err != nil {
		return nil, err
	}

	iv, err := GenerateRandomIV()
	if err != nil {
		return nil, err
	}

	encrypted, err := aes.Encrypt(fileKey, iv)
	if err != nil {
		return nil, err
	}

	// 返回IV + 加密后的密钥 + 标签
	result := make([]byte, 12 + len(encrypted))
	copy(result[:12], iv)
	copy(result[12:], encrypted)

	return result, nil
}

// DecryptKey 用主密钥解密文件密钥
func DecryptKey(encryptedKey []byte, masterKey []byte) ([]byte, error) {
	if len(encryptedKey) < 12 + 16 { // IV + 至少16字节密文+标签
		return nil, errors.New("invalid encrypted key length")
	}

	iv := encryptedKey[:12]
	ciphertext := encryptedKey[12:]

	aes, err := NewAESGCM(masterKey)
	if err != nil {
		return nil, err
	}

	key, err := aes.Decrypt(ciphertext, iv)
	if err != nil {
		return nil, err
	}

	return key, nil
}
