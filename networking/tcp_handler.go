package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"
)

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
				//@todo -implement proof of work
				replaceChain(newBlockchain)
			}
			btChan <- blockchain
			io.WriteString(conn, "\nEnter a new BPM:")
		}

	}()
	// broadcast blockchain to all nodes after every n seconds
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
	for _ = range btChan {
		spew.Dump(blockchain)
	}
}
