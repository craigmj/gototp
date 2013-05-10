package gototp

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"fmt"
	"math"
	"math/rand"
	"net/url"
	"strings"
	"time"
)

// Time-based One Time Password
type TOTP struct {
	key    []byte // Base-32 decoded Secret
	Digits int    // Number of digits for code. Defaults to 6.
	Period int    // Number of seconds for each code. Defaults to 30.
}

// Generate a new Time Based One Time Password with the given secret (a b32 encoded key)
func New(secretB32 string) (*TOTP, error) {
	// Decode the secret
	key, err := base32.StdEncoding.DecodeString(secretB32)
	if nil != err {
		return nil, fmt.Errorf("Error encountered base32 decoding secret: %v", err.Error())
	}

	return &TOTP{key, 6, 30}, nil
}

// Convert an integer into an 8 byte array
func int_to_bytestring(val int64) []byte {
	result := make([]byte, 8)
	i := len(result) - 1
	for i >= 0 {
		result[i] = byte(val & 0xff)
		i--
		val = val >> 8
	}
	return result
}

// Return the Time Based One Time Password for right now
func (totp *TOTP) Now() int32 {
	return totp.FromNow(0)
}

// Return the Time Based One Time Password for the time-period
// that is Now + the given periods.
// This is useful if you want to provide some flexibility around the acceptance
// of codes.
// For instance, you might want to accept a code that is valid in the current
// period (FromNow(0)), or that was valid in the previous period (FromNow(-1)) or that
// will be valid in the next period (FromNow(1)).
// This means that every code is therefore valid for 3 * totp.Period.
func (totp *TOTP) FromNow(periods int64) int32 {
	period := (time.Now().Unix() / int64(totp.Period)) + int64(periods)
	return totp.ForPeriod(period)
}

// Return the time-based OTP for the given period.
func (totp *TOTP) ForPeriod(period int64) int32 {
	data := int_to_bytestring(period)

	hmacHash := hmac.New(sha1.New, totp.key)
	hmacHash.Write(data)
	digest := hmacHash.Sum(nil)
	offset := int(digest[19] & 0xf)
	code := int32(digest[offset]&0x7f)<<24 |
		int32(digest[offset+1]&0xff)<<16 |
		int32(digest[offset+2]&0xff)<<8 |
		int32(digest[offset+3]&0xff)

	code = int32(int64(code) % int64(math.Pow10(totp.Digits)))
	return code
}

// Return the TOTP Secret base32 encoded
func (totp *TOTP) Secret() string {
	return base32.StdEncoding.EncodeToString(totp.key)
}

// Return the data to be contained in a QR Code for this TOTP with the given label.
func (totp *TOTP) QRCodeData(label string) string {
	// We need to URL Escape the label, but at the same time, spaces come through
	// as +'s, so we need to reverse that encoding...
	label = url.QueryEscape(label)
	label = strings.Replace(label, "+", " ", -1)
	return fmt.Sprintf("otpauth://totp/%v?secret=%v&Digits=%v&Period=%v", label, totp.Secret(), totp.Digits, totp.Period)
}

// Return a URL to generate a QRCode on Google Charts for the TOTP, with the given 
// label and width (and height equal to width).
func (totp *TOTP) QRCodeGoogleChartsUrl(label string, width int) string {
	return fmt.Sprintf("https://chart.googleapis.com/chart?cht=qr&chs=%vx%v&chl=%v", width, width, url.QueryEscape(totp.QRCodeData(label)))
}

// Generate a Random secret encoded as a b32 string
// If the length is <= 0, a default length of 10 bytes will
// be used, which will generate a secret of length 16.
func RandomSecret(length int, rnd *rand.Rand) string {
	if 0 <= length {
		length = 10
	}
	secret := make([]byte, length)
	for i, _ := range secret {
		secret[i] = byte(rnd.Int31() % 256)
	}
	return base32.StdEncoding.EncodeToString(secret)
}
