package otp

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"fmt"
	"hash"
)

// ErrValidateSecretInvalidBase32 Error when attempting to convert the secret from base32 to raw bytes.
var ErrValidateSecretInvalidBase32 = errors.New("Decoding of secret as base32 failed")

// ErrValidateInputInvalidLength The user provided passcode length was not expected.
var ErrValidateInputInvalidLength = errors.New("Code length unexpected")

// ErrValidatePasscode passcode is not valid
var ErrValidatePasscode = errors.New("Code is not valid")

// Algorithm represents the hashing function to use in the HMAC
// operation needed for OTPs.
type Algorithm int

const (
	//AlgorithmSHA1 default hash algorithm refer to SHA1
	//DO NOT CHANGE THIS ALGORITHM ORDER
	AlgorithmSHA1 Algorithm = iota
	//AlgorithmSHA256 refer to SHA256
	AlgorithmSHA256
	//AlgorithmSHA512 refer to SHA512
	AlgorithmSHA512
	//AlgorithmMD5 refer to MD5
	AlgorithmMD5
)

func (a Algorithm) String() string {
	switch a {
	case AlgorithmSHA1:
		return "SHA1"
	case AlgorithmSHA256:
		return "SHA256"
	case AlgorithmSHA512:
		return "SHA512"
	case AlgorithmMD5:
		return "MD5"
	}
	panic("Unknown OTP Algorithm")
}

//Hash create new hash object based on algorithm
func (a Algorithm) Hash() hash.Hash {
	switch a {
	case AlgorithmSHA1:
		return sha1.New()
	case AlgorithmSHA256:
		return sha256.New()
	case AlgorithmSHA512:
		return sha512.New()
	case AlgorithmMD5:
		return md5.New()
	}
	panic("Unknown OTP Algorithm")
}

// Digits represents the number of digits present in the
// user's OTP passcode. Six and Eight are the most common values.
type Digits int

const (
	//DigitsSix default digit is 6
	DigitsSix Digits = 6
	//DigitsEight option to use 8 digit code
	DigitsEight Digits = 8
)

// Format converts an integer into the zero-filled size for this Digits.
func (d Digits) Format(in int32) string {
	f := fmt.Sprintf("%%0%dd", d)
	return fmt.Sprintf(f, in)
}

// Length returns the number of characters for this Digits.
func (d Digits) Length() int {
	return int(d)
}

func (d Digits) String() string {
	return fmt.Sprintf("%d", d)
}
