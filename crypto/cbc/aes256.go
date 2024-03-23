package cbc

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

//Encrypt AES256 encryption using random generated IV, first block (16 byte) is for chipertext and
//the rest will be used for chipertext
func Encrypt(key, plaintext []byte, opts ...Opts) (ciphertext []byte, err error) {
	Pad := PKCS7Padding
	if len(opts) > 0 {
		Pad = opts[0].Pad
	}
	if len(plaintext)%aes.BlockSize != 0 {
		plaintext = Pad(plaintext, aes.BlockSize)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	//Bytes buff for chiper, aes.BlockSize + length plaintext
	ciphertext = make([]byte, aes.BlockSize+len(plaintext))

	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	cbc := cipher.NewCBCEncrypter(block, iv)
	cbc.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	return
}

//Decrypt AES256 decryption using IV from the chiper, first block (16 byte) is for chipertext
//and the rest is the chiper
func Decrypt(key, ciphertext []byte, opts ...Opts) (plaintext []byte, err error) {
	var block cipher.Block

	Unpad := PKCS7Trimming
	if len(opts) > 0 {
		Unpad = opts[0].Unpad
	}

	if block, err = aes.NewCipher(key); err != nil {
		return
	}

	if len(ciphertext) < aes.BlockSize {
		err = errors.New("Ciphertext too short")
		return
	}

	iv := ciphertext[:aes.BlockSize]

	ciphertext = ciphertext[aes.BlockSize:]

	cbc := cipher.NewCBCDecrypter(block, iv)
	cbc.CryptBlocks(ciphertext, ciphertext)

	plaintext, err = Unpad(ciphertext, aes.BlockSize)

	return
}

// EncryptByIV AES256 Encryption using custom IV
func EncryptByIV(key, plaintext, iv []byte, opts ...Opts) (ciphertext []byte, err error) {

	Pad := PKCS7Padding
	if len(opts) > 0 {
		Pad = opts[0].Pad
	}

	if len(plaintext)%aes.BlockSize != 0 {
		plaintext = Pad(plaintext, aes.BlockSize)
		if err != nil {
			return
		}
	}

	if len(iv)%aes.BlockSize != 0 {
		iv = Pad(iv, aes.BlockSize)
		if err != nil {
			return
		}
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	//Bytes buff for chiper, length plaintext
	ciphertext = make([]byte, len(plaintext))

	cbc := cipher.NewCBCEncrypter(block, iv)
	cbc.CryptBlocks(ciphertext, plaintext)

	return
}

// DecryptByIV AES256 decryption using custom IV
func DecryptByIV(key, ciphertext, iv []byte, opts ...Opts) (plaintext []byte, err error) {
	var block cipher.Block

	Unpad := PKCS7Trimming
	Pad := PKCS7Padding

	if len(opts) > 0 {
		Unpad = opts[0].Unpad
		Pad = opts[0].Pad
	}

	if len(iv)%aes.BlockSize != 0 {
		iv = Pad(iv, aes.BlockSize)
		if err != nil {
			return
		}
	}

	if block, err = aes.NewCipher(key); err != nil {
		return
	}

	if len(ciphertext) < aes.BlockSize {
		err = errors.New("Ciphertext too short")
		return
	}

	cbc := cipher.NewCBCDecrypter(block, iv)
	cbc.CryptBlocks(ciphertext, ciphertext)

	plaintext, err = Unpad(ciphertext, aes.BlockSize)

	return
}
