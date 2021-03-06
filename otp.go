package Requireris

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
)

// Digits number for the end code
var Digits int64 = 6

// CounterSize is the counter byte size, defined in the RFC
const CounterSize int = 8

// OTP holds the base32 formatted secret key
type OTP struct {
	SecretKey string
}

func parseGoogleSecret(secret string) (string, error) {
	// google encodes secrets in the form of "xxxx xxxx xxxx xxxx"
	// we need to remove whitespace and uppercase it "XXXXXXXXXXXXXXXX"
	secret = strings.ToUpper(strings.Replace(secret, " ", "", -1))
	_, err := base32.StdEncoding.DecodeString(secret)

	if err != nil {
		return "", errors.New("non google secret")
	}
	return secret, nil
}

func Init(secret string, codeLength int) *OTP {
	res, err := parseGoogleSecret(secret)
	if err == nil {
		// create a base32 secret
		// use decoded google secret
		secret = res
	} else {
		// check if secret is already a base32 secret
		_, err := base32.StdEncoding.DecodeString(secret)
		if err != nil {
			// create a base32 secret
			secret = base32.StdEncoding.EncodeToString([]byte(secret))
		}
	}
	if codeLength >= 6 && codeLength <= 10 {
		Digits = int64(codeLength)
	} else {
		fmt.Println("We only support codes from 6 to 10 digits; using default 6 as value")
	}
	return &OTP{
		SecretKey: secret,
	}
}

func truncate(s []byte) uint64 {
	var code int32

	offset := int(s[19] & 0x0F)
	p := s[offset : offset+4]

	code = int32((p[0] & 0x7f)) << 24
	code |= int32((p[1] & 0xff)) << 16
	code |= int32((p[2] & 0xff)) << 8
	code |= int32((p[3] & 0xff))

	return uint64(code) & 0x7FFFFFFF
}

func (h *OTP) genOtp(c uint64) string {
	// setting counter
	counter := new([CounterSize]byte)
	binary.BigEndian.PutUint64(counter[:], c)

	byteSecret, err := base32.StdEncoding.DecodeString(h.SecretKey)
	if err != nil {
		fmt.Println(err)
		return "error"
	}

	// setting up hash
	mac := hmac.New(sha1.New, byteSecret)
	mac.Write(counter[:])

	// getting the value
	hash := mac.Sum(nil)
	res := truncate(hash)

	// getting the int
	retVal := res % uint64((math.Pow(10, float64(Digits))))
	paddedRet := fmt.Sprintf("%06d", int(retVal))

	return paddedRet
}

func (h *OTP) TOTP() string {
	secs := time.Now().Unix()
	return h.genOtp(uint64(secs / 30))
}

func (h *OTP) HOTP(C uint64) string {
	return h.genOtp(C)
}
