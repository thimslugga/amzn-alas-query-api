package main

import (
	"log"
	"net/http"
)

func main() {
	// Get the config, global var in config.go
	config = NewConfig()

	log.Println("Initializing redis database")
	var err error

	// db is a global in redisdb.go
	db, err = NewDatabase()
	if err != nil {
		log.Fatal(err)
	}

	// kick off the refresh loop to keep ALAS in sync with Redis
	go db.RefreshALASLoop()

	// Listen and serve requests
	router := NewRouter()
	log.Println("Starting to serve requests on", config.ListenAddr)
	log.Fatal(http.ListenAndServe(config.ListenAddr, router))
}
