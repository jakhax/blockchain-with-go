package main

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"sync"
	"time"
)

// block - items in a blockchain
type Block struct {
	Index     int    `json:"index"`     //position of data record on blockchain
	Timestamp string `json:"timestamp"` // auto time of when data is added to the blockchain
	BPM       int    `json:"bpm"`       // beats per minute - heartrate
	Hash      string `json:"hash"`      //SHA256 identifier of this data record
	PrevHash  string `json:"prevHash"`  //SHA256 identifier of the previous data on the block chain
}

type Heartrate struct {
	BPM int
}
type Blockchain []Block

var mutex = &sync.Mutex{}

var blockchain Blockchain

func calculateHash(b Block) string {
	record := strconv.Itoa(b.Index) + b.Timestamp + strconv.Itoa(b.BPM) + b.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	return hex.EncodeToString(h.Sum(nil))
}

//generate new Block function
func generateBlock(oldBlock Block, BPM int) (Block, error) {
	var newBlock Block
	t := time.Now()
	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.BPM = BPM
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)
	return newBlock, nil
}

// block validation
func isBlockvalid(newBlock, oldBlock Block) bool {
	if newBlock.Index != oldBlock.Index+1 {
		return false
	}
	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}
	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}
	return true
}

//classical blockchain problem
/* if we receive blockchains from 2 nodes from our block chain network which do we take ?
sol: the longest - its up to date..assuming
*/
func replaceChain(newBlocks []Block) {
	if len(blockchain) < len(newBlocks) {
		blockchain = newBlocks
	}
	return
}
func main() {

}
