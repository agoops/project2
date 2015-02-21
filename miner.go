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
	"math/rand"
	"time"
	"github.com/PointCoin/btcjson"
	"math/big"
	"github.com/PointCoin/btcutil"
	"github.com/PointCoin/btcwire"
	"github.com/PointCoin/btcrpcclient"


	// "github.com/PointCoin/pointcoind/blockchain"
)

const (
	// This should match your settings in pointcoind.conf
	rpcuser = "ankitgupta"
	rpcpass = "password"
	// This file should exist if pointcoind was setup correctly
	cert    = "/home/ubuntu/.pointcoind/rpc.cert"
)

func GetClient() *btcrpcclient.Client {
	return setupRpcClient(cert, rpcuser, rpcpass)
}
func GetTemplate(client *btcrpcclient.Client) *btcjson.GetBlockTemplateResult {

	// Get a new block template from pointcoind.
	log.Printf("Requesting a block template\n")
	template, err := client.GetBlockTemplate(&btcjson.TemplateRequest{})
	if err != nil {
		log.Fatal(err)
	}
	return template
}

func SetUpBlock(template *btcjson.GetBlockTemplateResult) (*btcwire.MsgBlock, big.Int) {
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
	// old miner address: PcvNdeZ6UsC8R34NQqi9YyxKfSXQcjoWHy
	a := "Prxy397nCyskwHwmiv3TaFG6ZgZ88Cbnju" 

	coinbaseTx := CreateCoinbaseTx(height, a, msg)

	txs = prepend(coinbaseTx.MsgTx(), txs)
	merkleRoot := createMerkleRoot(txs)

	// Finish the miner!
	// print("Looking for valid block....")
	source := rand.NewSource(time.Now().UnixNano())
	myRand := rand.New(source)
	var nonce uint32 = myRand.Uint32();

	block := CreateBlock(prevHash, merkleRoot, difficulty, nonce, txs)

	return block, difficulty
	
}



func mainn() {
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
		source := rand.NewSource(time.Now().UnixNano())
		myRand := rand.New(source)
		var nonce uint32 = myRand.Uint32();

		block := CreateBlock(prevHash, merkleRoot, difficulty, nonce, txs)

		// Naive approach, try 20 million times before checking for new template
		// Loop, checking hashes against difficulty until valid hash is found
		for i := 0; i < 20000000; i++ {
			blockSha, _ := block.Header.BlockSha()
			// print("Trying with nonce: ", block.Header.Nonce)

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
