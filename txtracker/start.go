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
	"strings"
	"encoding/json"
	// "regexp"
	// "math/rand"
	"log"
)


func main() {
	print := fmt.Println
	s := "\"Hello\""
	fmt.Println(s)

	o,_ := strconv.Unquote(s)
	fmt.Println(o)

	// command = pointctl getrawtransaction d2011b19dea6e98ec8bf78bd224856e76b6a9c460bbb347e49adb3dcf457e548
	txid := "d2011b19dea6e98ec8bf78bd224856e76b6a9c460bbb347e49adb3dcf457e548"
	txdetails := getTransactionDetails(txid)
	txdetailsbytes := []byte(txdetails)


	var f interface{}
	_ = json.Unmarshal(txdetailsbytes, &f)
	m := f.(map[string]interface{})


	// for k, v := range m {
	//     switch vv := v.(type) {
	// 	    case string:
	// 	        fmt.Println(k, "is string", vv)
	// 	    case int:
	// 	        fmt.Println(k, "is int", vv)
	// 	    case []interface{}:
	// 	        fmt.Println(k, "is an array:")
	// 	        for i, u := range vv {
	// 	            fmt.Println(i, u)
	// 	        }
	// 	    default:
	// 	        fmt.Println(k, "is of a type I don't know how to handle")
	// 	    }
	// }

	txidreturned  := m["txid"]

	vinList := make([]vin,0)
	vinJsonList := m["vin"]

	switch vv := vinJsonList.(type) {
	case []interface{}:
		for _, u := range vv {
					j := u.(map[string]interface{})
					vinTxid := j["txid"]
					vinVout := j["vout"]
					newVin := vin{txid: vinTxid.(string), vout: vinVout.(int)}
					vinList = append(vinList, newVin)
		            // fmt.Println(i, u)
		        }
		// print("yes matches")
	default:
		print("nope didn't work")
	}

	for x := range vinList {
		print(x)
	}
	



	_ = vinList
	print("\n\ngot txid", txidreturned)

}


type vin struct {
    txid string
    vout  int
}

func getTransactionDetails(txhash string) (string){
	// command = pointctl getrawtransaction d2011b19dea6e98ec8bf78bd224856e76b6a9c460bbb347e49adb3dcf457e548
	cmd := exec.Command("pointctl", "getrawtransaction", txhash)

	var out bytes.Buffer
    cmd.Stdout = &out
    err := cmd.Run()
    if err != nil {
    	log.Fatal(err)
    }
    // fmt.Printf("result: %s\n", out)
    // fmt.Println(out.String())


	cmd2 := exec.Command("xargs", "pointctl", "decoderawtransaction")
	cmd2.Stdin = strings.NewReader(out.String())
	var out2 bytes.Buffer
    cmd2.Stdout = &out2
    err2 := cmd2.Run()
    if err2 != nil {
    	log.Fatal(err2)
    }
	
    return out2.String()
    // fmt.Println(out2.String())
}




























