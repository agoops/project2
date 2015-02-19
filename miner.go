//
// miner.go
// Staring template for PointCoint miner.
//
// cs4501: Cryptocurrency Cafe
// University of Virginia, Spring 2015
// Project 2
//

package main

import (
	"fmt"
	"log"
	// "math/big"
	"math/rand"

	"github.com/PointCoin/btcjson"
	"github.com/PointCoin/btcutil"
	// "github.com/PointCoin/pointcoind/blockchain"
)

const (
	// This should match your settings in pointcoind.conf
	rpcuser = "ankitgupta"
	rpcpass = "password"
	// This file should exist if pointcoind was setup correctly
	cert    = "/home/ubuntu/.pointcoind/rpc.cert"
)

func main() {
	print := fmt.Println
	// Setup the client using application constants, fail horribly if there's a problem
	client := setupRpcClient(cert, rpcuser, rpcpass)

	for { // Loop forever (you may want to do something smarter!)
		// Get a new block template from pointcoind.
		log.Printf("Requesting a block template\n")
		template, err := client.GetBlockTemplate(&btcjson.TemplateRequest{})
		if err != nil {
			log.Fatal(err)
		}

		// The template returned by GetBlockTemplate provides these fields that 
		// you will need to use to create a new block:

		// The hash of the previous block
		prevHash := template.PreviousHash

		// The difficulty target
		difficulty := formatDiff(template.Bits)

		// The height of the next block (number of blocks between genesis block and next block)
		height := template.Height

		// The transactions from the network	
		txs := formatTransactions(template.Transactions) 

		// These are configurable parameters to the coinbase transaction
		msg := "ag7bf" // replace with your UVa Computing ID (e.g., "dee2b")
		
		// This is my address generated with wallet of passphrase "mywallet", 
		a := "PcvNdeZ6UsC8R34NQqi9YyxKfSXQcjoWHy" // replace with the address you want mining fees to go to (or leave it like this and Nick gets them)

		coinbaseTx := CreateCoinbaseTx(height, a, msg)

		txs = prepend(coinbaseTx.MsgTx(), txs)
		merkleRoot := createMerkleRoot(txs)

		// Finish the miner!
		print("Looking for valid block....")
		var nonce uint32 = rand.Uint32();

		block := CreateBlock(prevHash, merkleRoot, difficulty, nonce, txs)

		// Loop, checking hashes against difficulty until valid hash is found
		for {
			blockSha, _ := block.Header.BlockSha()
			print("Trying with nonce: ", block.Header.Nonce)

			if lessThanDiff(blockSha, difficulty) {
				// submit block
				print("valid hash found")
				err := client.SubmitBlock(btcutil.NewBlock(block), nil)
				print(err)
				break
			} else {
				if block.Header.Nonce == 4294967294 {
					block.Header.Nonce = 0
				} else {
					block.Header.Nonce += 1
				}
			}
		}
	}
}
