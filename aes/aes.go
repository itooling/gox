package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"math/rand"
	"time"
)

func GenAesKey(key []byte) []byte {
	gen := make([]byte, aes.BlockSize)
	copy(gen, key)
	for i := 16; i < len(key); {
		for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
			gen[j] ^= key[i]
		}
	}
	return gen
}

func GenAesKeyRand(l ...int) []byte {
	rd := rand.New(rand.NewSource(time.Now().Unix()))
	length := aes.BlockSize
	if len(l) > 0 {
		length = l[0]
	}
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		b := rd.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return bytes
}

func EncryptECB(src, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	length := (len(src) + aes.BlockSize) / aes.BlockSize
	dst := make([]byte, length*aes.BlockSize)
	copy(dst, src)

	padding := byte(len(dst) - len(src))
	for i := len(src); i < len(dst); i++ {
		dst[i] = padding
	}

	en := make([]byte, len(dst))
	for bs, be := 0, block.BlockSize(); bs <= len(src); bs, be = bs+block.BlockSize(), be+block.BlockSize() {
		block.Encrypt(en[bs:be], dst[bs:be])
	}

	return en, nil
}
func DecryptECB(src, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	de := make([]byte, len(src))
	for bs, be := 0, block.BlockSize(); bs < len(src); bs, be = bs+block.BlockSize(), be+block.BlockSize() {
		block.Decrypt(de[bs:be], src[bs:be])
	}

	padding := 0
	if len(de) > 0 {
		padding = len(de) - int(de[len(de)-1])
	}

	return de[:padding], nil
}

func EncryptCBC(src, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	size := block.BlockSize()

	padding := size - len(src)%size
	text := bytes.Repeat([]byte{byte(padding)}, padding)
	src = append(src, text...)

	mode := cipher.NewCBCEncrypter(block, key[:size])
	en := make([]byte, len(src))
	mode.CryptBlocks(en, src)
	return en, nil
}

func DecryptCBC(src, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	size := block.BlockSize()
	mode := cipher.NewCBCDecrypter(block, key[:size])
	de := make([]byte, len(src))
	mode.CryptBlocks(de, src)

	length := len(de)
	padding := int(de[length-1])
	de = de[:(length - padding)]

	return de, nil
}
