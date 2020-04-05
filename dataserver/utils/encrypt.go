package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"io"
)

// AesEncrypt 加密并返回加密字符串
func AesEncrypt(orig string, key string) string {
	// 转成字节数组
	origData := []byte(orig)
	k := []byte(key)
	// 分组秘钥
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	block, _ := aes.NewCipher(k)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 补全码
	origData = PKCS7Padding(origData, blockSize)
	// 加密模式
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	// 创建数组
	cryted := make([]byte, len(origData))
	// 加密
	blockMode.CryptBlocks(cryted, origData)
	return base64.StdEncoding.EncodeToString(cryted)
}

// AesDecrypt 解密
func AesDecrypt(cryted string, key string) string {
	// 转成字节数组
	crytedByte, _ := base64.StdEncoding.DecodeString(cryted)
	k := []byte(key)
	// 分组秘钥
	block, _ := aes.NewCipher(k)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	// 创建数组
	orig := make([]byte, len(crytedByte))
	// 解密
	blockMode.CryptBlocks(orig, crytedByte)
	// 去补全码
	orig = PKCS7UnPadding(orig)
	return string(orig)
}

// PKCS7Padding 补码
//AES加密数据块分组长度必须为128bit(byte[16])，密钥长度可以是128bit(byte[16])、192bit(byte[24])、256bit(byte[32])中的任意一个。
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

// AESEncrypt 加密
func AESEncrypt(src []byte, key []byte) (encrypted []byte) {
	cipher, _ := aes.NewCipher(generateKey(key))
	length := (len(src) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, src)
	pad := byte(len(plain) - len(src))
	for i := len(src); i < len(plain); i++ {
		plain[i] = pad
	}
	encrypted = make([]byte, len(plain))
	// 分组分块加密
	for bs, be := 0, cipher.BlockSize(); bs <= len(src); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Encrypt(encrypted[bs:be], plain[bs:be])
	}

	return encrypted
}

// AESDecrypt 解密
func AESDecrypt(encrypted []byte, key []byte) (decrypted []byte) {
	cipher, _ := aes.NewCipher(generateKey(key))
	decrypted = make([]byte, len(encrypted))
	//
	for bs, be := 0, cipher.BlockSize(); bs < len(encrypted); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}

	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}

	return decrypted[:trim]
}

// generateKey generateKey
func generateKey(key []byte) (genKey []byte) {
	genKey = make([]byte, 16)
	copy(genKey, key)
	for i := 16; i < len(key); {
		for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	return genKey
}

// GetHmacCode 返回hash值
func GetHmacCode(s string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// EncodeURLBase64 数据进base64编码
func EncodeURLBase64(input string) string {
	return base64.URLEncoding.EncodeToString([]byte(input))
}

// DecodeBase64 数据进base64解码
func DecodeBase64(input string) string {
	decodeBytes, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return ""
	}
	return string(decodeBytes)
}

// DecodeURLBase64 数据进base64解码
func DecodeURLBase64(input string) string {
	decodeBytes, err := base64.URLEncoding.DecodeString(input)
	if err != nil {
		return ""
	}
	return string(decodeBytes)
}

// 将json字符串转换为map
func ParseJSONStr2Map(jsonstr string) (*map[string]interface{}, error) {
	meta := make(map[string]interface{})
	err2 := json.Unmarshal([]byte(jsonstr), &meta)
	if err2 != nil {
		logs.Error("parseJSONStr2Map:解析json字符串时发生错误，%s", err2)
		return nil, err2
	}
	return &meta, nil
}

func ParseJSONBytes2Map(jsonbyte []byte) (*map[string]interface{}, error) {
	meta := make(map[string]interface{})
	err2 := json.Unmarshal(jsonbyte, &meta)
	if err2 != nil {
		logs.Error("parseJSONStr2Map:解析json字符串时发生错误，%s", err2)
		return nil, err2
	}
	return &meta, nil
}
