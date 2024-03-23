package otp

import (
	"testing"
)

func TestAlgorithmOrder(t *testing.T) {

	if AlgorithmSHA1 != 0 {
		t.Errorf("Algorithm SHA1 in wrong order should be at first order")
	}
	if AlgorithmSHA256 != 1 {
		t.Errorf("AlgorithmSHA256 in wrong order should be at second order")
	}
	if AlgorithmSHA512 != 2 {
		t.Errorf("AlgorithmSHA512 in wrong order should be at third order")
	}
	if AlgorithmMD5 != 3 {
		t.Errorf("AlgorithmMD5 in wrong order should be at fourth order")
	}

}

func TeshStringHash(t *testing.T) {
	if AlgorithmSHA1.String() != "SHA1" {
		t.Errorf("Wrong hash for SHA1")
	}
	if AlgorithmSHA256.String() != "SHA256" {
		t.Errorf("Wrong hash for SHA256")
	}
	if AlgorithmSHA512.String() != "SHA512" {
		t.Errorf("Wrong hash for SHA512")
	}
	if AlgorithmMD5.String() != "MD5" {
		t.Errorf("Wrong hash for MD5")
	}
}
