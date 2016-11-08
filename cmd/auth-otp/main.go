package main

import (
  "fmt"
  "github.com/rotta-f/Requireris"
)

func main() {
  // example with non base32 secret
  //otp := Requireris.Init("elelelele")
  otp := Requireris.Init("w6f5 fky2 vf5y 2vc7 6npa ds3j 46em shts")
  fmt.Println(otp.TOPT())
}
