// This must be run on a computer with pointcoind running and pointctl installed


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
	print("\n\n")
	var txdetails string
	for {
		var inputTx string
		fmt.Printf("%s", "Enter tx hash: ")
		fmt.Scanf("%s", &inputTx)

		txdetails = getTransactionDetails(inputTx)
		if txdetails != "no" {
			break
		}
	}
	for {
		txdetailsbytes := []byte(txdetails)
		var f interface{}
		_ = json.Unmarshal(txdetailsbytes, &f)
		m := f.(map[string]interface{})
		txidreturned  := m["txid"]
		print("\n\nTransaction ID:", txidreturned)

		vinList := getVinList(m)
		voutList := getVoutList(m)

		prevOutputs := make([]string, 0)

		fmt.Println("From:")
		for i, x := range vinList {
			index := strconv.Itoa(i)

			if x.coinbase == true{
				print ("\t[" + string(index) + "] Coinbase Transaction (10 PTC)")
				prevOutputs = append(prevOutputs, "coinbase")
				continue
			}
			tx := getTransactionDetails(x.txid)
			prevOutputs = append(prevOutputs, tx)
			txjs := getTransactionJson(tx)
			txvouts := getVoutList(txjs)
			for _, y := range txvouts {

				if y.n == x.vout {
					fmt.Println("\t[" + string(index) + "]",y.addresses[0], "(" + FloatToString(y.value) + " PTC)")
					break
				}
			}
		}

		fmt.Println("To:")
		for i, x := range voutList {
			index := strconv.Itoa(i)
			fmt.Println("\t[" + string(index) + "] " + x.addresses[0] + " (" + FloatToString(x.value) + " PTC)" )
		}
		

		nextIndex := getIndexInput(len(prevOutputs),"\nEnter index of \"From\" address to see output tx:")
		savedDetails := txdetails
		txdetails = prevOutputs[nextIndex]
		if txdetails == "coinbase" {
			fmt.Println("Oh you think Pointcoin is your ally, but you merely adopted the pointcoin. This coinbase transaction was born in it. Molded by it.")
			txdetails = savedDetails
		}

	}


}



type vin struct {
	coinbase bool
    txid string
    vout  int
}

type vout struct {
	value float64
	n int
	addresses []string
}

func getIndexInput(size int, msg string)(int) {
	for {
		fmt.Printf("%s", msg)
		var nextIndex int
		fmt.Scanf("%d", &nextIndex)
		// fmt.Println(nextIndex)
		if nextIndex < 0 || nextIndex >= size {
			fmt.Println("Invalid index")
			continue
		}
		return nextIndex
	}
}
func getVinList(m map[string]interface{}) ([]vin) {
	vinList := make([]vin,0)
	vinJsonList := m["vin"]

	switch vv := vinJsonList.(type) {
	case []interface{}:
		for _, u := range vv {
			j := u.(map[string]interface{})
			var newVin vin
			if _,ok := j["coinbase"]; ok {
				// this is a coinbase transaction w/ coinbase input
				newVin = vin{coinbase:true, txid:"null", vout:0} 
			} else {
				vinTxid := j["txid"].(string)
				vinVout := int(j["vout"].(float64))
				newVin = vin{coinbase:false, txid: vinTxid, vout: vinVout}
	            // fmt.Println(i, u)
			}
			vinList = append(vinList, newVin)

        }
		// print("yes matches")
	default:
		print("nope getVinList didn't work")
	}

	// fmt.Println("vins:")
	// for _,x := range vinList {
	// 	fmt.Println(x)
	// }
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

	// fmt.Println("vouts:")
	// for _,x := range voutList {
	// 	fmt.Println(x)
	// }
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
    	fmt.Println("Invalid ash or no Connection to pointcoind. Try again.")
    	return "no"
    	// log.Fatal(err)
    }


	cmd2 := exec.Command("xargs", "pointctl", "decoderawtransaction")
	cmd2.Stdin = strings.NewReader(out.String())
	var out2 bytes.Buffer
    cmd2.Stdout = &out2
    err2 := cmd2.Run()
    if err2 != nil {
    	log.Fatal(err2)
    }
    // fmt.Println(out2.String())
	
    return out2.String()
}

func FloatToString(input_num float64) string {
    // to convert a float number to a string
    return strconv.FormatFloat(input_num, 'f', -1, 64)
}






























