package main

import (
  "fmt"
  "github.com/rotta-f/Requireris"
)

func main() {
  otp := Requireris.Init("w6f5 fky2 vf5y 2vc7 6npa ds3j 46em shts")
  fmt.Println(otp.TOPT())
}
