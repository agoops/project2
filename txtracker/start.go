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
	// address: Prxy397nCyskwHwmiv3TaFG6ZgZ88Cbnju
	// command = pointctl getrawtransaction c1de1be883834d733d096b3e14674978459f111f90d9dfbc5a82c9fa20db60a7
	
	txid := "c1de1be883834d733d096b3e14674978459f111f90d9dfbc5a82c9fa20db60a7"
	txdetails := getTransactionDetails(txid)
	// m := getTransactionJson(txdetails)
	txdetailsbytes := []byte(txdetails)

	var f interface{}
	_ = json.Unmarshal(txdetailsbytes, &f)
	m := f.(map[string]interface{})
	txidreturned  := m["txid"]
	print("\n\ngot txid", txidreturned)

	vinList := getVinList(m)
	voutList := getVoutList(m)

	_ , _ = vinList,voutList



	// Start with transaction

	// See input addresses of transaction as well as amounts
	// For each vin, going to have to 
	fmt.Println("Outputs")
	for i, x := range vinList {
		tx := getTransactionDetails(x.txid)
		txjs := getTransactionJson(tx)
		txvouts := getVoutList(txjs)
		for _, y := range txvouts {

			if y.n == x.vout {
				fmt.Println("\t[",i,"]",y.addresses[0],y.value)
				break
			}
		}
	}
		// 1) Get tx, 
		// 2) Get 'n'th output
		// 3) Get address and amount of that output


	// See output addresses as well as amounts
	// For each vout
		// 1)Print address and amount



}


type vin struct {
    txid string
    vout  int
}

type vout struct {
	value float64
	n int
	addresses []string
}

func getVinList(m map[string]interface{}) ([]vin) {
	vinList := make([]vin,0)
	vinJsonList := m["vin"]

	switch vv := vinJsonList.(type) {
	case []interface{}:
		for _, u := range vv {
			j := u.(map[string]interface{})
			vinTxid := j["txid"].(string)
			vinVout := int(j["vout"].(float64))
			newVin := vin{txid: vinTxid, vout: vinVout}
			vinList = append(vinList, newVin)
            // fmt.Println(i, u)
        }
		// print("yes matches")
	default:
		print("nope getVinList didn't work")
	}

	fmt.Println("vins:")
	for _,x := range vinList {
		fmt.Println(x)
	}
	return vinList

}

func getVoutList(m map[string] interface{}) ([]vout) {
	voutList := make([]vout,0)
	voutJsonList := m["vout"]

	switch oo := voutJsonList.(type) {
	case []interface{}:
		for _,u := range oo {
			j := u.(map[string]interface{})
			voutVal := j["value"].(float64)
			voutN := int(j["n"].(float64))

			vScriptPubKey := j["scriptPubKey"].(map[string]interface{})
			vAddresses := vScriptPubKey["addresses"].([]interface{})
			vAddressesStrings := make([]string, 0)
			for _,u := range vAddresses {
				addr := u.(string)
				vAddressesStrings = append(vAddressesStrings, addr)
			}

			newVout := vout{value: voutVal, n: voutN, addresses: vAddressesStrings}
			voutList = append(voutList, newVout)
		}
	}

	fmt.Println("vouts:")
	for _,x := range voutList {
		fmt.Println(x)
	}
	return voutList
}

func getTransactionJson(txdetails string) (map[string]interface{}){
	txdetailsbytes := []byte(txdetails)

	var f interface{}
	_ = json.Unmarshal(txdetailsbytes, &f)
	m := f.(map[string]interface{})
	return m
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




























