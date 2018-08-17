Simple Blockchain
===================
## Description

Simple technology to keep records of user heart rate

## Quick Start

- `git clone https://github.com/jakhax/blockchain-with-go.git`
- `cd blockchain-with-go/basic-prototype`
- `touch .env && echo "PORT=8080" > .env`
- `go run main.go http_server.go` or   `go build`

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
