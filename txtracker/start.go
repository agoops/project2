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
	"strconv"
	// "strings"
	// "regexp"
	// "math/rand"
	"log"
)

func main() {
	
	s := "\"Hello\""
	fmt.Println(s)

	o,_ := strconv.Unquote(s)
	fmt.Println(o)









	// command = pointctl getrawtransaction d2011b19dea6e98ec8bf78bd224856e76b6a9c460bbb347e49adb3dcf457e548
	txid := "d2011b19dea6e98ec8bf78bd224856e76b6a9c460bbb347e49adb3dcf457e548"
	cmd := exec.Command("pointctl", "getrawtransaction", txid)

	var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    if err != nil {
    	log.Fatal(err)
    }
    // fmt.Printf("result: %s\n", out)
    fmt.Println(out.String())


	cmd2 := exec.Command("pointctl", "decoderawtransaction", out.String())
	var out2 bytes.Buffer
    cmd2.Stdout = &out2
    err2 := cmd2.Run()
    if err2 != nil {
    	log.Fatal(err2)
    }
	

    fmt.Println(out2)

}