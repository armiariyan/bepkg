package totp

import (
	"crypto/hmac"
	"crypto/subtle"
	"encoding/binary"
	"math"
	"strings"
	"time"

	"gitlab.com/gobang/bepkg/otp"
)

// DefaultInterval set 30 as default
const DefaultInterval = 30

// DefaultSkew 1 as default
const DefaultSkew = 1

// DefaultDigits 6 as default
const DefaultDigits = otp.DigitsSix

// Opts provides options for Generate().  The default values
// are compatible with Google-Authenticator.
type Opts struct {
	// Number of seconds a TOTP hash is valid for. Defaults to 30 seconds.
	Interval uint
	// Secret to store. Defaults to a randomly generated secret of SecretSize.  You should generally leave this empty.
	Secret string
	// Digits as part of the input. Defaults to 6.
	Digits otp.Digits
	// Algorithm to use for HMAC. Defaults to SHA1.
	Algorithm otp.Algorithm
	// Starting time 0
	T0 int
	// Periods before or after the current time to allow.  Value of 1 allows up to Period
	// of either side of the specified time.  Defaults to 0 allowed skews.  Values greater
	// than 1 are likely sketchy.
	Skew uint
}

// Generate Create new otp passcode
func Generate(opts Opts) (string, error) {
	if opts.Interval == 0 {
		opts.Interval = DefaultInterval
	}

	if opts.Digits == 0 {
		opts.Digits = DefaultDigits
	}

	//timestamp := time.Now().UTC().Unix()
	//counter := uint64((timestamp - int64(opts.T0)) / int64(opts.Interval))
	timestamp := time.Now().UTC().UnixNano() / int64(time.Millisecond)
	counter := uint64(math.Floor((float64(timestamp/1000) - float64(opts.T0)) / float64(opts.Interval)))

	return generateCode(opts.Secret, counter, opts)

}

func generateCode(secret string, counter uint64, opts Opts) (string, error) {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, counter)
	mac := hmac.New(opts.Algorithm.Hash, []byte(secret))

	mac.Write(buf)
	sum := mac.Sum(nil)

	// "Dynamic truncation" in RFC 4226
	// http://tools.ietf.org/html/rfc4226#section-5.4
	offset := sum[len(sum)-1] & 0xf
	bin := int64(((int(sum[offset]) & 0x7f) << 24) |
		((int(sum[offset+1] & 0xff)) << 16) |
		((int(sum[offset+2] & 0xff)) << 8) |
		(int(sum[offset+3]) & 0xff))

	l := opts.Digits.Length()
	mod := int32(bin % int64(math.Pow10(l)))

	return opts.Digits.Format(mod), nil
}

// Validate examine if existing passcode still valid
func Validate(passcode string, secret string, t time.Time, opts Opts) (bool, error) {
	if opts.Interval == 0 {
		opts.Interval = DefaultInterval
	}
	if opts.Digits == 0 {
		opts.Digits = DefaultDigits
	}
	if opts.Skew == 0 {
		opts.Skew = DefaultSkew
	}

	counters := []uint64{}
	//counter := uint64((t.Unix() - int64(opts.T0)) / int64(opts.Interval))
	//counter := uint64(math.Floor(float64(t.Unix()) / float64(opts.Interval)))
	timestamp := t.UnixNano() / int64(time.Millisecond)
	counter := uint64(math.Floor((float64(timestamp/1000) - float64(opts.T0)) / float64(opts.Interval)))

	counters = append(counters, uint64(counter))
	for i := 1; i <= int(opts.Skew); i++ {
		counters = append(counters, uint64(counter+uint64(i)))
		counters = append(counters, uint64(counter-uint64(i)))
	}

	for _, counter := range counters {
		rv, _ := validateCode(passcode, counter, secret, opts)

		if rv == true {
			return true, nil
		}
	}

	return false, otp.ErrValidatePasscode
}

func validateCode(passcode string, counter uint64, secret string, opts Opts) (bool, error) {
	passcode = strings.TrimSpace(passcode)

	//fmt.Printf("Counter: %d, Passcode: %s (%d), Digits Len: %d\n", counter, passcode, len(passcode), opts.Digits.Length())

	if len(passcode) != opts.Digits.Length() {
		return false, otp.ErrValidateInputInvalidLength
	}

	otpstr, err := generateCode(secret, counter, opts)
	if err != nil {
		return false, err
	}
	//fmt.Printf("otpstr: %s passcode: %s\n", otpstr, passcode)
	if subtle.ConstantTimeCompare([]byte(otpstr), []byte(passcode)) == 1 {
		return true, nil
	}

	return false, otp.ErrValidatePasscode
}
