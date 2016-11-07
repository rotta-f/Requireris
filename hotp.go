package Requireris

import (
  "fmt"
  "crypto/hmac"
  "crypto/sha1"
  "encoding/binary"
  "math"
)

const Digits int64 = 6
const CounterSize int = 8

type HOTP struct {
  SecretKey string
}

func GetHOTP(K string) *HOTP {
  return &HOTP{
    SecretKey : K,
  }
}

func truncate(s []byte) uint64 {
  var code int32

  offset := int(s[19] & 0x0F)
  p := s[offset : offset + 4]

  code = int32((p[0] & 0x7f)) << 24
  code |= int32((p[1] & 0xff)) << 16
  code |= int32((p[2] & 0xff)) << 8
  code |= int32((p[3] & 0xff))

  return uint64(code) & 0x7FFFFFFF
}

func bytesToNum(b []byte) int {
  var P int = 0
  var i uint = 0
  for ; i < 4; i++ {
    P = P | (int(b[i] & 0xFF) << (24 - i * 8))
  }
  return P
}

func (h *HOTP) At(C uint64) uint64 {
  fmt.Println(h.SecretKey)
  // setting counter
  counter := new([CounterSize]byte)
  binary.BigEndian.PutUint64(counter[:], C)

  // setting up hash
  mac := hmac.New(sha1.New, []byte(h.SecretKey))
  mac.Write(counter[:])

  // getting the value
  hash := mac.Sum(nil)
  res := truncate(hash)

  // getting the int and % to Digits length
  retVal := res % uint64((math.Pow(10, float64(Digits))))
  fmt.Println(retVal)
  return retVal
  /*
    fmt.Printf("secret '%s'\n", []byte(secret))
    //hh = hmac.new(base64.b32decode('base32secret3232', casefold=True), b"4", hashlib.sha1)
    mac.Write([]byte(t))
    fmt.Printf("MAC: %X\n%s\n", mac.Sum(nil), h.SecretKey + t)

    Sbits := truncate(mac.Sum(nil))
    SNum := bytesToNum(Sbits)
    fmt.Printf("HEX Sbits: %X\nSNum: %X\n", Sbits, SNum)
    fmt.Printf("DEC Sbits: %d\nSNum: %d\n", Sbits, SNum)
    D := SNum % int((math.Pow(10, 6)))
    fmt.Println(D)
    */
}
