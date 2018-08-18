Simple Blockchain Network 
===================
## Description

Simple blockchain network built on top of TCP with no decentralization, nodes add blocks to a tcp server.

## Basic Code Review

- assuming we have all the code from the previous basic prototype minus the http implementation
- first lets create a tcp server
```go
func createTcpServer() {
	tcpPort, ok := os.LookupEnv("PORT")
	if !ok {
		tcpPort = "9999"
	}
	s, err := net.Listen("tcp", ":"+tcpPort)
	if err != nil {
		log.Fatal("ERROR", err)
	}
	defer s.Close()
	for {
		conn, err := s.Accept()
		if err != nil {
			log.Fatal("ERROR:", err)
		}
		go handleTcpConn(conn)
	}
}
```
### Adding a pulse to the blockchain via tcp

- prompt the client to enter their BPM
- scan the client’s input from stdin
- create a new block with this data, using the generateBlock, isBlockValid, and replaceChain functions we created previously
- put the new blockchain in the channel we created to be broadcast to the network
- allow the client to enter a new BPM
```go
func handleTcpConn(conn net.Conn) {
	defer conn.Close()
	io.WriteString(conn, "Enter a new BPM: ")
	scanner := bufio.NewScanner(conn)
	go func() {
		for scanner.Scan() {
			bpm, err := strconv.Atoi(scanner.Text())
			if err != nil {
				log.Println("ERROR:", err)
				io.WriteString(conn, "\nInvalid input, expected an integer\nEnter a new BPM: ")
				continue
			}
			newBlock, err := generateBlock(blockchain[len(blockchain)-1], bpm)
			if err != nil {
				log.Println("ERROR:", err)
				io.WriteString(conn, "\nServer error! Try again, \nEnter a new BPM: ")
				continue
			}
			if isBlockvalid(newBlock, blockchain[len(blockchain)-1]) {
				newBlockchain := append(blockchain, newBlock)
				replaceChain(newBlockchain)
			}
			btChan <- blockchain
			io.WriteString(conn, "\nEnter a new BPM:")
		}
    }()
    //broadcasting code goes here
    for _ = range btChan {
		spew.Dump(blockchain)
	}
```

### Broadcast blockchain to all nodes
- We will new blockchain to all the connections being served by our TCP server. Since we’re coding this up on one computer, we’ll be simulating how data gets transmitted to all clients. We will be first convert the blockchain to json  then we print it to clients terminals. We do this periodically after n seconds.
- Under the last line of code, add the following.
```go
	go func() {
		for {
			time.Sleep(30 * time.Second)
			//use a mutex to stop access to the blockchain
			mutex.Lock()
			j, err := json.Marshal(blockchain)
			if err != nil {
				log.Fatal(err)
			}
			mutex.Unlock()
			io.WriteString(conn, string(j)+"\n")
		}
	}()
```