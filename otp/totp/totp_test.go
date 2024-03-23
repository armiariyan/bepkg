package totp

import (
	"fmt"
	"testing"
	"time"
)

func doTaskEvery(d time.Duration, f func(time.Time) bool) {
	for x := range time.Tick(d) {
		if ok := f(x); !ok {
			break
		}
	}
}

func TestGenerateCode(t *testing.T) {
	opts := Opts{
		Secret: "8cbefbe60ab72710758c954c1fcdc15c853559ff1c43d958fb816fb35c819053b95745682bed10cf957ea427004786d24be74ae675ddf331f8d720812c6a383b620213f144e6dd425c63a069f13418e4f6a6037d9bb2e678d8473fb40cb9658ea8a7bec56183f68b141a809356868a7c5671a26bc43465f267bc38ad6a44fe2a1909361878ee66e4661bbf11a81cf5b1c5e353051d4dbcee5d6fc7dee604cca60d41e1dce7acb2ae64c7e16da1be0cc77c1ec99000e268e3565a589790dc5e7801aba4dc5244f2b67d6b134e9bee5b8287eb757a0da1c496e195f54255e829ca023a67ebac0f81c63f096bba6a91e174b87264bd250650685010339c50f2cd868c3d95d16ad81969001d3bec3cfd87c1a9be54ed65f6f5a6b7d2a7065df1a5fc061ea573f54b08bb07ff96103c171d8aac5adba48560fbe3fd3e57a0c582b0fd944607d9e5d63686bf64bbe65ada6cd7176550eb15508359ccd812010734a67d131479b275acc98cea94eeef5851cc3e62c44e288448e35ddeac21c14ec3d2",
	}
	code, err := Generate(opts)
	if err != nil {
		t.Errorf("Err, %s", err)
	}

	fmt.Printf("TOTP: %s\n", code)
}

func TestValidateCode(t *testing.T) {
	opts := Opts{
		Secret:   "5cbefbe60ab72710758c954c1fcdc15c853559ff1c43d958fb816fb35c819053b95745682bed10cf957ea427004786d24be74ae675ddf331f8d720812c6a383b620213f144e6dd425c63a069f13418e4f6a6037d9bb2e678d8473fb40cb9658ea8a7bec56183f68b141a809356868a7c5671a26bc43465f267bc38ad6a44fe2a1909361878ee66e4661bbf11a81cf5b1c5e353051d4dbcee5d6fc7dee604cca60d41e1dce7acb2ae64c7e16da1be0cc77c1ec99000e268e3565a589790dc5e7801aba4dc5244f2b67d6b134e9bee5b8287eb757a0da1c496e195f54255e829ca023a67ebac0f81c63f096bba6a91e174b87264bd250650685010339c50f2cd868c3d95d16ad81969001d3bec3cfd87c1a9be54ed65f6f5a6b7d2a7065df1a5fc061ea573f54b08bb07ff96103c171d8aac5adba48560fbe3fd3e57a0c582b0fd944607d9e5d63686bf64bbe65ada6cd7176550eb15508359ccd812010734a67d131479b275acc98cea94eeef5851cc3e62c44e288448e35ddeac21c14ec3d2",
		Digits:   6,
		Interval: 60,
		Skew:     1,
	}

	code, err := Generate(opts)
	if err != nil {
		t.Errorf("Err, %s", err)
	}
	//code = "243578"

	fmt.Printf("TOTP: %s\n", code)

	isValidated, err := Validate(code, opts.Secret, time.Now().UTC(), opts)
	if err != nil {
		t.Errorf("Err, %s", err)
	}
	if !isValidated {
		t.Errorf("Err, Code %s is not valid", code)
	}
}

// func TestGenerateCodePerSec(t *testing.T) {
// 	secret, _ := new(big.Int).SetString("8cbefbe60ab72710758c954c1fcdc15c853559ff1c43d958fb816fb35c819053b95745682bed10cf957ea427004786d24be74ae675ddf331f8d720812c6a383b620213f144e6dd425c63a069f13418e4f6a6037d9bb2e678d8473fb40cb9658ea8a7bec56183f68b141a809356868a7c5671a26bc43465f267bc38ad6a44fe2a1909361878ee66e4661bbf11a81cf5b1c5e353051d4dbcee5d6fc7dee604cca60d41e1dce7acb2ae64c7e16da1be0cc77c1ec99000e268e3565a589790dc5e7801aba4dc5244f2b67d6b134e9bee5b8287eb757a0da1c496e195f54255e829ca023a67ebac0f81c63f096bba6a91e174b87264bd250650685010339c50f2cd868c3d95d16ad81969001d3bec3cfd87c1a9be54ed65f6f5a6b7d2a7065df1a5fc061ea573f54b08bb07ff96103c171d8aac5adba48560fbe3fd3e57a0c582b0fd944607d9e5d63686bf64bbe65ada6cd7176550eb15508359ccd812010734a67d131479b275acc98cea94eeef5851cc3e62c44e288448e35ddeac21c14ec3d2", 16)
// 	interval := 60 * time.Second
// 	digit := 6
// 	t0 := 0

// 	topts := Opts{
// 		Interval:  uint(interval.Seconds()),
// 		Digits:    otp.Digits(digit),
// 		T0:        t0,
// 		Secret:    fmt.Sprintf("%x", secret),
// 		Algorithm: 2,
// 	}

// 	var initPasscode string

// 	doTaskEvery(1*time.Second, func(t time.Time) bool {
// 		fmt.Printf("Time: %v \n", t)
// 		passcode, _ := Generate(topts)
// 		fmt.Println(passcode)

// 		if initPasscode == "" {
// 			initPasscode = passcode
// 		}

// 		if initPasscode != passcode {
// 			return false
// 		}

// 		return true
// 	})

// }

// func TestValidateCodePerSec(t *testing.T) {

// 	secret, _ := new(big.Int).SetString("8cbefbe60ab72710758c954c1fcdc15c853559ff1c43d958fb816fb35c819053b95745682bed10cf957ea427004786d24be74ae675ddf331f8d720812c6a383b620213f144e6dd425c63a069f13418e4f6a6037d9bb2e678d8473fb40cb9658ea8a7bec56183f68b141a809356868a7c5671a26bc43465f267bc38ad6a44fe2a1909361878ee66e4661bbf11a81cf5b1c5e353051d4dbcee5d6fc7dee604cca60d41e1dce7acb2ae64c7e16da1be0cc77c1ec99000e268e3565a589790dc5e7801aba4dc5244f2b67d6b134e9bee5b8287eb757a0da1c496e195f54255e829ca023a67ebac0f81c63f096bba6a91e174b87264bd250650685010339c50f2cd868c3d95d16ad81969001d3bec3cfd87c1a9be54ed65f6f5a6b7d2a7065df1a5fc061ea573f54b08bb07ff96103c171d8aac5adba48560fbe3fd3e57a0c582b0fd944607d9e5d63686bf64bbe65ada6cd7176550eb15508359ccd812010734a67d131479b275acc98cea94eeef5851cc3e62c44e288448e35ddeac21c14ec3d2", 16)
// 	interval := 60 * time.Second
// 	digit := 6
// 	t0 := 0

// 	opts := Opts{
// 		Interval:  uint(interval.Seconds()),
// 		Digits:    otp.Digits(digit),
// 		T0:        t0,
// 		Secret:    fmt.Sprintf("%x", secret),
// 		Algorithm: 2,
// 		Skew:      2,
// 	}

// 	code, err := Generate(opts)
// 	if err != nil {
// 		t.Errorf("Err, %s", err)
// 	}

// 	fmt.Printf("TOTP: %s\n", code)

// 	doTaskEvery(1*time.Second, func(t time.Time) bool {
// 		fmt.Printf("Time: %v\n", t)
// 		isValidated, err := Validate(code, opts.Secret, time.Now().UTC(), opts)
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 		return isValidated
// 	})

// }
