package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

var (
	slackToken  = ""
	kibelaTeam  = ""
	kibelaToken = ""
	listenPort  = 30000
)

func getEnv() {
	slackToken = os.Getenv("SLACK_TOKEN")
	if slackToken == "" {
		panic("SLACK_TOKEN is empty")
	}
	kibelaTeam = os.Getenv("KIBELA_TEAM")
	if kibelaTeam == "" {
		panic("KIBELA_TEAM is empty")
	}
	kibelaToken = os.Getenv("KIBELA_TOKEN")
	if kibelaToken == "" {
		panic("KIBELA_TOKEN is empty")
	}
	port := os.Getenv("PORT")
	if port != "" {
		var err error
		listenPort, err = strconv.Atoi(port)
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	getEnv()
	http.Handle("/", NewEventHandler(slackToken, kibelaTeam, kibelaToken))
	listenAddr := fmt.Sprintf(":%d", listenPort)
	log.Printf("Server listening on %s", listenAddr)
	if err := http.ListenAndServe(listenAddr, nil); err != nil {
		log.Println(err)
	}
}
