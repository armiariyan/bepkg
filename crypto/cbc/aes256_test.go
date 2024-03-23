package cbc

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestAES256EncryptDecrypt(t *testing.T) {

	key := []byte("32423423n432423423n432423423n412")

	plaintext := []byte("P|1fd94d876c3c8c55016c747cc647000a|14840627")

	ciphertext, err := Encrypt(key, plaintext)
	if err != nil {
		t.Errorf("Something happen, %s", err)
	}

	cleartext := base64.StdEncoding.EncodeToString(ciphertext)

	fmt.Printf("CBC: %s\n", cleartext)

	plaintextDecripted, err := Decrypt(key, ciphertext)

	if string(plaintext) != string(plaintextDecripted) {
		t.Errorf("wrong plaintext decrypted, should be\n %v \ninstead of\n %v", plaintext, plaintextDecripted)
	}
	fmt.Printf("Decrypted Plaintext: %s\n", plaintextDecripted)

}

func TestAES256EncryptDecryptByIV(t *testing.T) {

	resultToCompareChipertext := "fWJ3Z97ffg3vjkQYqy2gORFh9zXfvm7gNvmSTCCIswEx6X4USj3FejvdOPi9fAvf"

	key := []byte("32423423n432423423n432423423n412")

	plaintext := []byte("P|1fd94d876c3c8c55016c747cc647000a|14840627")
	iv, _ := hex.DecodeString("31353636383038383337383539000000")

	ciphertext, err := EncryptByIV(key, plaintext, iv)
	if err != nil {
		t.Errorf("Something happen, %s", err)
	}

	cleartext := base64.StdEncoding.EncodeToString(ciphertext)

	if cleartext != resultToCompareChipertext {
		t.Errorf("wrong chipertext, should be %s instead of %s", resultToCompareChipertext, cleartext)
	}
	fmt.Printf("CBC: %s\n", cleartext)

	plaintextDecripted, err := DecryptByIV(key, ciphertext, iv)

	if string(plaintext) != string(plaintextDecripted) {
		t.Errorf("wrong plaintext decrypted, should be\n %v \ninstead of\n %v", plaintext, plaintextDecripted)
	}
	fmt.Printf("Decrypted Plaintext: %s\n", plaintextDecripted)

}

func TestSec(t *testing.T) {
	converted := []byte("Tc4sh2018!123456") //[]byte(keyStr)

	key := KeyToSha256(string(converted), 32)

	iv := []byte("1568637457542")
	plaintext := []byte("1234" + "111111" + string(iv))

	ciphertext, err := EncryptByIV(key, plaintext, iv, Opts{
		Pad: PKCS5Padding,
	})
	if err != nil {
		t.Errorf("Something happen, %s", err)
	}

	cleartext := base64.StdEncoding.EncodeToString(ciphertext)
	fmt.Printf("KEY: %s\n", key)
	fmt.Printf("Chipertext: %s\n", cleartext)

	x, _ := DecryptByIV(key, ciphertext, iv, Opts{
		Unpad: PKCS5Trimming,
		Pad:   PKCS5Padding,
	})
	fmt.Printf("X Text: %s\n", x)
	fmt.Printf("X Textv(byre): %v\n", []byte(x))
}
