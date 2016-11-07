package main

import (
  "time"
  "fmt"
  "github.com/rotta-f/Requireris"
)

func main() {
  hotp := Requireris.GetHOTP("base32secret3232")
  now := time.Now()
  secs := now.Unix()
  fmt.Println("secs: ", secs, " | secs/30: ", secs / 30)
  hotp.At(uint64(secs / 30))
}
