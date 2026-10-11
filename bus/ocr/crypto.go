// Package crypto .
package ocr

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"trpc.group/trpc-go/trpc-go/log"
)

// pkcs7Padding
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// pkcs7UnPadding
func pkcs7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("加密字符串错误！")
	}
	unPadding := int(data[length-1])
	return data[:(length - unPadding)], nil
}

// AesEncryptToBase64 加密并编码
func AesEncryptToBase64(data string, key string) (string, error) {
	log.Infof("AesEncryptToBase64 before len:%d", len(data))
	dst, err := AesEncrypt([]byte(data), []byte(key))
	if err != nil {
		log.Errorf("AesEncryptToBase64 err:%v", err)
		return "", err
	}
	return Base64Encode(dst), nil
}

// AesEncrypt 加密
func AesEncrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	encryptBytes := pkcs7Padding(data, blockSize)
	crypted := make([]byte, len(encryptBytes))
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	blockMode.CryptBlocks(crypted, encryptBytes)
	return crypted, nil
}

// AesDecryptFromBase64 解码并解密
func AesDecryptFromBase64(data string, key string) ([]byte, error) {
	dst, err := Base64Decode(data)
	if err != nil {
		return nil, err
	}
	return AesDecrypt(dst, []byte(key))
}

// AesDecrypt 解密
func AesDecrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	// 检查输入长度
	if len(data) == 0 {
		return nil, errors.New("input data is empty")
	}
	if len(data)%blockSize != 0 {
		return nil, fmt.Errorf("input not full blocks: data length %d is not multiple of block size %d", len(data), blockSize)
	}
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	crypted := make([]byte, len(data))
	blockMode.CryptBlocks(crypted, data)
	crypted, err = pkcs7UnPadding(crypted)
	if err != nil {
		return nil, err
	}
	return crypted, nil
}

// Sha256 加密
func Sha256(data []byte, key []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}

// Base64Encode 编码
func Base64Encode(data []byte) string {
	log.Infof("Base64Encode before len:%d", len(data))
	return base64.StdEncoding.EncodeToString(data)
}

// Base64Decode 解码
func Base64Decode(string string) ([]byte, error) {
	ba, err := base64.StdEncoding.DecodeString(string)
	if err != nil {
		return nil, err
	}
	log.Infof("Base64Decode after len:%d", len(ba))
	return ba, nil
}

// CalculateMD5 计算字节数组的 MD5 哈希值并返回十六进制字符串
func CalculateMD5(bytes []byte) string {
	// 创建 MD5 哈希对象
	hash := md5.New()

	// 更新哈希对象，传入字节数组
	hash.Write(bytes)

	// 计算哈希值并转换为十六进制字符串
	return hex.EncodeToString(hash.Sum(nil))
}
