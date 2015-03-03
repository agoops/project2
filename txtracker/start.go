package main

import (
	// "encoding/hex"
	"fmt"
	"os/exec"
	"bytes"
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
	txid := "d2011b19dea6e98ec8bf78bd224856e76b6a9c460bbb347e49adb3dcf457e548"
	cmd := exec.Command("pointctl", "getrawtransaction", txid)
	// cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    if err != nil {
    	log.Fatal(err)
    }
    fmt.Printf("result: %q\n", out.String())
}