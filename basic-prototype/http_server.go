package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/get-blockchain", handleGetBlockChain).Methods("GET")
	muxRouter.HandleFunc("/write-block", handleWriteBlock).Methods("POST")
	return muxRouter
}

func runServer() error {
	//route handler
	r := makeMuxRouter()
	httpPort, ok := os.LookupEnv("PORT")
	if !ok {
		httpPort = "8080"
	}
	s := &http.Server{
		Addr:           ":" + httpPort,
		Handler:        handlers.LoggingHandler(os.Stdout, r),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if err := s.ListenAndServe(); err != nil {
		return err
	}
	//success
	log.Println("HTTP server listening on port:", httpPort)
	return nil
}

// route handlers
func handleGetBlockChain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(blockchain)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return
}

func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var heartrate Heartrate
	if err := json.NewDecoder(r.Body).Decode(&heartrate); err != nil {
		writeBlockJsonResponse(w, r, http.StatusBadRequest, heartrate)
		return
	}
	defer r.Body.Close()

	// lock to prevent a go routine from modifying newBlock
	mutex.Lock()
	newBlock, err := generateBlock(blockchain[len(blockchain)-1], heartrate.BPM)
	mutex.Unlock()
	if err != nil {
		writeBlockJsonResponse(w, r, http.StatusInternalServerError, heartrate)
		return
	}
	if isBlockvalid(newBlock, blockchain[len(blockchain)-1]) {
		blockchain = append(blockchain, newBlock)
		//print newBlock to console
		spew.Dump(newBlock)
	}
	writeBlockJsonResponse(w, r, http.StatusCreated, newBlock)
	return

}

//response formater for the writeblock handler
func writeBlockJsonResponse(w http.ResponseWriter, r *http.Request, status int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "	")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
	}
	w.WriteHeader(status)
	w.Write(response)
}
