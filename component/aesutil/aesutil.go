// Package aesutil aesutil
package aesutil

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
)

// AesEncryptStr aes加密
func AesEncryptStr(orig string, aesKey string) (string, error) {
	if orig == "" {
		return "", nil
	}
	iv, err := GetIv()
	if err != nil {
		return "", err
	}

	defer func() {
		recoverErr := recover()
		if recoverErr != nil {
			err = fmt.Errorf("panic occurred: %v", recoverErr)
		}
	}()
	ucode, err := AesEncrypt(orig, iv, aesKey)
	if err != nil || ucode == nil || len(ucode) == 0 {
		return "", errors.New("加密失败")
	}
	data := append(iv, ucode...)
	dataStr := base64.StdEncoding.EncodeToString(data)

	return dataStr, nil
}

// AesDecryptStr aes解密
func AesDecryptStr(ucode string, aesKey string) (string, error) {
	if ucode == "" {
		return "", nil
	}
	UcodeByte, err := base64.StdEncoding.DecodeString(ucode)
	if err != nil {
		return "", err
	}
	if len(UcodeByte) <= aes.BlockSize {
		return "", errors.New("参数异常")
	}
	iv := UcodeByte[:aes.BlockSize]
	data := UcodeByte[aes.BlockSize:]
	defer func() {
		recoverErr := recover()
		if recoverErr != nil {
			err = fmt.Errorf("panic occurred: %v", recoverErr)
		}
	}()
	// 解密
	ucodeInfoStr, err := AesDecrypt(data, iv, aesKey)
	if err != nil {
		return "", err
	}
	return ucodeInfoStr, nil
}

// AesEncrypt 加密
func AesEncrypt(orig string, ivByte []byte, aesKey string) ([]byte, error) {

	if orig == "" {
		return nil, nil
	}
	block, err := aes.NewCipher([]byte(aesKey))
	if err != nil {
		log.Println("Error creating AES cipher:", err)
		return nil, err
	}
	// 转成字节数组
	origData := []byte(orig)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 补全码
	origData = PKCS7Padding(origData, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, ivByte)
	// 创建数组
	cryted := make([]byte, len(origData))
	// 加密
	blockMode.CryptBlocks(cryted, origData)
	return cryted, nil
}

// AesDecrypt 解密
func AesDecrypt(crytedByte, ivByte []byte, aesKey string) (string, error) {
	block, err := aes.NewCipher([]byte(aesKey))
	if err != nil {
		log.Println("Error creating AES cipher:", err)
		return "", err
	}
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, ivByte)
	// 创建数组
	orig := make([]byte, len(crytedByte))
	// 解密
	blockMode.CryptBlocks(orig, crytedByte)
	// 去补全码
	orig = PKCS7UnPadding(orig)
	return string(orig), nil
}

// GetIv 获取随机值
func GetIv() ([]byte, error) {
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Println("Error generating random IV:", err)
		return nil, err
	}
	return iv, nil
}

// PKCS7Padding 补码 AES加密数据块分组长度必须为128bit(byte[16])，密钥长度可以是128bit(byte[16])、192bit(byte[24])、256bit(byte[32])中的任意一个。
func PKCS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS7UnPadding 去码
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
