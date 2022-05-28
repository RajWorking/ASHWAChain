package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

// web server
func run() error {
	mux := makeMuxRouter()
	httpPort := PORT
	log.Println("HTTP Server Listening on port :", httpPort)
	s := &http.Server{
		Addr:           ":" + httpPort,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

// create handlers
func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleGetPowChain).Methods("GET")
	muxRouter.HandleFunc("/txs", handleGetTxChain).Methods("GET")
	muxRouter.HandleFunc("/committee", handleGetComChain).Methods("GET")

	muxRouter.HandleFunc("/", handleWritePowBlock).Methods("POST")
	muxRouter.HandleFunc("/tx", handleWriteTransaction).Methods("POST")
	return muxRouter
}
