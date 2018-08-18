Simple Blockchain
===================
## Description

Simple technology to keep records of user heart rate built on http.

## Quick Start

- `git clone https://github.com/jakhax/blockchain-with-go.git`
- `cd blockchain-with-go/basic-prototype`
- `touch .env && echo "PORT=8080" > .env`
- `go run main.go http_server.go` or   `go build`

## Basic Code Review

- We create a block struct and our blockchain as a slice of block.
```go
type Block struct {
	Index     int    `json:"index"`     //position of data record on blockchain
	Timestamp string `json:"timestamp"` // auto time of when data is added to the blockchain
	BPM       int    `json:"bpm"`       // beats per minute - heartrate
	Hash      string `json:"hash"`      //SHA256 identifier of this data record
	PrevHash  string `json:"prevHash"`  //SHA256 identifier of the previous data on the block chain
}
type Blockchain []Block
```

- We then create a hashing function to generate hash of a block
```go
func calculateHash(b Block) string {
	record := strconv.Itoa(b.Index) + b.Timestamp + strconv.Itoa(b.BPM) + b.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	return hex.EncodeToString(h.Sum(nil))
}
```

- Write a function to generate new block
```go
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
```

- Write a function to validate a new block hecking its `PrevHash` against the `Hash` of the previous block. 
```go
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
```

## Endpoints
- You can use a browser, curl or [Postman](https://www.getpostman.com/apps) to test the endpoints

### View Entire Chain
``` bash
GET get-blockchain
# Request sample
curl --request GET http://localhost:8080/get-blockchain
```

### Create a Block
``` bash
POST write-block
# Request sample
curl --header "Content-Type: application/json" --request POST\
 --data '{"BPM":78}' http://localhost:8080/write-block
```
