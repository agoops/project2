package main

import (
	// "encoding/hex"
	"fmt"
	"time"
	"math/big"
	"github.com/PointCoin/btcutil"
	"github.com/PointCoin/btcwire"
	"github.com/PointCoin/btcrpcclient"
	"github.com/PointCoin/btcjson"

	// "strings"
	// "regexp"
	// "math/rand"
	// "log"
)
// func workerWork() {
// 	for {
// 		fmt.Println("hashing")
// 		time.Sleep(time.Second * 3)

// 	}
// }
func main() {

	// var maxNonce uint32 = 4294967294
	// _ = maxNonce
	client := GetClient()
	for {
		template := GetTemplate(client)
		block,difficulty := SetUpBlock(template)
		

		newTemplateChan := make(chan bool, 1)
		done := make(chan bool, 1)
		go FindValidBlock(newTemplateChan, done, block,difficulty,client)
		go TemplateChecker(newTemplateChan, template, client)

		<- done	
	}
	
}


func TemplateChecker(newTemplateChan chan bool, template *btcjson.GetBlockTemplateResult, client *btcrpcclient.Client) {
	prevHash := template.PreviousHash

	for i := 0; i < 100000000; i++ {
		time.Sleep(time.Second * 10)
		otherTemplate := GetTemplate(client)
		if prevHash != otherTemplate.PreviousHash {
			newTemplateChan <- true
		}
	}
	
}

func FindValidBlock(newTemplateChan chan bool, done chan bool, block *btcwire.MsgBlock, difficulty big.Int, client *btcrpcclient.Client) {
	hashesThenCheck := 500000

	for {
		select {
		case <- newTemplateChan:
			fmt.Println("New Template detected before we found block. Resetting...")
			done <- true
			return
		default:
			for i := 0; i < hashesThenCheck; i++ {
				blockSha, _ := block.Header.BlockSha()

				if lessThanDiff(blockSha, difficulty) {
					// submit block
					print("valid hash found")
					err := client.SubmitBlock(btcutil.NewBlock(block), nil)
					print(err)
					done <- true
					return
				} else {
					if block.Header.Nonce == 4294967294 {
						block.Header.Nonce = 0
					} else {
						block.Header.Nonce += 1
					}
				}
			}
			fmt.Println(hashesThenCheck, "hashes done. Checking for new template")
		}

	}
}








