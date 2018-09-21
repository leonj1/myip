package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"io/ioutil"
	"fmt"
	"gopkg.in/natefinch/lumberjack.v2"
	"strings"
	"regexp"
)

/**
Hold an ip address and redirect traffic to anyone that attempts to connect to me to that ip address

TODO:
- lock down to some client api
- add support to redirect other HTTP verbs
 */

type MyStruct struct {
	ip string
}

func (*MyStruct) validIP4(ipAddress string) bool {
	ipAddress = strings.Trim(ipAddress, " ")

	re, _ := regexp.Compile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
	if re.MatchString(ipAddress) {
		return true
	}
	return false
}

func (m *MyStruct) addHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response := &ErrorResponse{ErrorMessage: "Message not found"}
		respondWithJSON(w, 404, response)
		return
	}
	payload := string(body)
	if payload == "" {
		response := &ErrorResponse{ErrorMessage: "Body is empty"}
		respondWithJSON(w, 404, response)
		return
	}
	if !m.validIP4(payload) {
		response := &ErrorResponse{ErrorMessage: "Not a valid ip"}
		respondWithJSON(w, 404, response)
		return
	}
	if payload != m.ip {
		m.ip = payload
	}
	response := &PostResponse{Status: "OK"}
	respondWithJSON(w, 200, response)
}

func (m *MyStruct) redirectHandler(w http.ResponseWriter, r *http.Request) {
	redirectIp := fmt.Sprintf("http://%s", m.ip)
	http.Redirect(w, r, redirectIp, 301)
}

func main() {
	log.SetOutput(&lumberjack.Logger{
		Filename:   "/tmp/myip.log",
		MaxSize:    5, // megabytes
		MaxBackups: 3,
		MaxAge:     3, //days
	})
	foo := &MyStruct{}
	s := mux.NewRouter()
	s.HandleFunc("/", foo.addHandler).Methods("POST")
	s.HandleFunc("/", foo.redirectHandler).Methods("GET")

	log.Printf("Staring HTTPS service on %s ...\n", ":443")
	if err := http.ListenAndServe(":443", s); err != nil {
		panic(err)
	}
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

type PostResponse struct {
	Status string `json:"status"`
}

type ErrorResponse struct {
	ErrorMessage string `json:"err_msg"`
}
