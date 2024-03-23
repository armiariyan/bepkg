package cbc

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

//PadFn type padding
type PadFn func([]byte, int) []byte

//TrimFn type unpad or trim the pad
type TrimFn func([]byte, int) ([]byte, error)

//Opts cbc options
type Opts struct {
	Pad   PadFn
	Unpad TrimFn
}

//KeyToSha256 conver key to sha256 with left trim byte
func KeyToSha256(keystr string, byteTrim int) []byte {
	converted := []byte(keystr)

	hasher := sha256.New()
	hasher.Write(converted)

	keys := hex.EncodeToString(hasher.Sum(nil))
	keys = keys[:byteTrim]

	key := make([]byte, byteTrim)
	copy(key, []byte(keys))

	return key
}

// PKCS7Padding add extra byte for encryption purpose
func PKCS7Padding(buf []byte, size int) []byte {
	bufLen := len(buf)
	padLen := size - bufLen%size
	padded := make([]byte, bufLen+padLen)
	copy(padded, buf)
	for i := 0; i < padLen; i++ {
		padded[bufLen+i] = byte(padLen)
	}
	return padded
}

// PKCS7Trimming remove extra byte pad from decrypted plaintext
func PKCS7Trimming(padded []byte, size int) ([]byte, error) {
	if len(padded)%size != 0 {
		return nil, errors.New("Invalid PKCS7 padding size")
	}

	bufLen := len(padded) - int(padded[len(padded)-1])
	buf := make([]byte, bufLen)
	copy(buf, padded[:bufLen])
	return buf, nil
}

// PKCS5Padding using pkcs 5 for padding
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS5Trimming using pkcs 5 unpad
func PKCS5Trimming(encrypt []byte, size int) ([]byte, error) {
	padding := encrypt[len(encrypt)-1]
	return encrypt[:len(encrypt)-int(padding)], nil
}
