package Requireris

import (
  "fmt"
  "math"
  "strconv"
  "encoding/base32"
  "crypto/hmac"
  "crypto/sha1"
)

type HOTP struct {
  SecretKey string
}

func GetHOTP(K string) *HOTP {
  return &HOTP{
    SecretKey : K,
  }
}

func truncate(s []byte) []byte {
  lenS := len(s)
  if (lenS != 20) {
    // TODO: Generate error, must be 20 what else can it be with hmac?
  }
  offset := s[19] & 0x0F
  P := make([]byte, 4)
  P[3] = s[offset] & 0x7F
  P[2] = s[offset + 1] & 0xFF
  P[1] = s[offset + 2] & 0xFF
  P[0] = s[offset + 3] & 0xFF
  var Pi int = 0
  Pi = ((int(s[offset]) & 0x7F) << 24 |
        (int(s[offset + 1]) & 0xFF) << 16 |
        (int(s[offset + 2]) & 0xFF) << 8 |
        (int(s[offset + 3]) & 0xFF))
  fmt.Printf("P: %X | Pi: %X\n", P, Pi)
  return P
}

func bytesToNum(b []byte) int {
  var P int = 0
  var i uint = 0
  for ; i < 4 ; i++ {
    P = P | (int(b[i] & 0xFF) << (24 - i * 8))
  }
  return P
}

func (h *HOTP) At(C int) {
  t := strconv.Itoa(C)
  enc := base32.StdEncoding.EncodeToString([]byte(h.SecretKey))
  mac := hmac.New(sha1.New, []byte(enc))
  //hh = hmac.new(base64.b32decode('base32secret3232', casefold=True), b"4", hashlib.sha1)
  mac.Write([]byte(t))
  fmt.Printf("MAC: %X\n%s\n", mac.Sum(nil), h.SecretKey + t)

  Sbits := truncate(mac.Sum(nil))
  SNum := bytesToNum(Sbits)
  fmt.Printf("HEX Sbits: %X\nSNum: %X\n", Sbits, SNum)
  fmt.Printf("DEC Sbits: %d\nSNum: %d\n", Sbits, SNum)
  D := SNum % int((math.Pow(10, 6)))
  fmt.Println(D)
}
