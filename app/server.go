package main

import (
	"log"
	"net/http"
	"os"
)

func serve() {
	http.HandleFunc("/endpoint", sendResponse())

	if err := http.ListenAndServe(":****", nil); err != nil {
		panic(err)
	}
}

func sendResponse() func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		fileData, err := os.ReadFile("config.json")
		if err != nil {
			http.Error(w, "status unavailable", http.StatusInternalServerError)
			log.Print("Error: ", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(fileData)
	}
}
