package main

import (
	// "encoding/hex"
	"fmt"
	"os/exec"
	// "bytes"
	// "time"
	// "math/big"
	// "github.com/PointCoin/btcutil"
	// "github.com/PointCoin/btcwire"
	// "github.com/PointCoin/btcrpcclient"
	// "github.com/PointCoin/btcjson"

	// "strings"
	// "regexp"
	// "math/rand"
	"log"
)

func main() {
	// command = pointctl getrawtransaction d2011b19dea6e98ec8bf78bd224856e76b6a9c460bbb347e49adb3dcf457e548
	txid := "d2011b19dea6e98ec8bf78bd224856e76b6a9c460bbb347e49adb3dcf457e548"
	out, err := exec.Command("pointctl", "getrawtransaction", txid).Output()
	s := string(out[:])
	// cmd.Stdin = strings.NewReader("some input")
    // err := cmd.Run()
    if err != nil {
    	log.Fatal(err)
    }
    // fmt.Printf("result: %s\n", out)
    fmt.Println(s)


	out, err = exec.Command("pointctl", "decoderawtransaction", s).Output()
	s = string(out[:])

	if err != nil {
    	log.Fatal(err)
    }

    fmt.Println(out)

}