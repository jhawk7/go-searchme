package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/jhawk7/go-searchme/pkg/common"
	"github.com/jhawk7/go-searchme/pkg/groupme_api"
	"strings"
)

func main() {
	router := mux.NewRouter()
	fmt.Println("Listening on port 8888...")
	router.HandleFunc("/", RootHandler).Methods("GET")
	router.HandleFunc("/issearchmeup", HealthCheck).Methods("GET")
	router.HandleFunc("/groupme/groups", GetUserGroups).Methods("GET")
	router.HandleFunc("/groupme/group/messages", GetGroupMessages).Methods("GET")
	router.HandleFunc("/groupme/group/messages/{keyword:[a-zA-Z0-9 ]+}", GetGroupMessages).Methods("GET")
	log.Fatal(http.ListenAndServe(":8888", router))
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	response := &common.Response{
		Status: http.StatusOK,
		Body:   "Welcome to go-searchme!",
	}
	data, jsonErr := json.Marshal(response)
	if jsonErr != nil {
		fmt.Printf("Error creating response")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(jsonErr)
	} else {
		w.WriteHeader(response.Status)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := &common.Response{
		Status: http.StatusOK,
		Body:   "go-searchme is up!",
	}
	data, jsonErr := json.Marshal(response)
	if jsonErr != nil {
		fmt.Printf("Error creating response")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(jsonErr)
	} else {
		w.WriteHeader(response.Status)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}

func GetUserGroups(w http.ResponseWriter, r *http.Request) {
	config := common.GetConfig()
	gmClient := groupme_api.CreateGroupmeClient(config)
	groupsResponse, err := gmClient.GetUserGroups()
	common.ErrorHandler(w, err)
	data, jsonErr := json.Marshal(&groupsResponse)
	common.ErrorHandler(w, jsonErr)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func GetGroupMessages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	config := common.GetConfig()
	gmClient := groupme_api.CreateGroupmeClient(config)
	groupMessages, err := gmClient.GetGroupMessages()
	if vars["keyword"] != "" {
		parsedMessages := ParseGroupMessages(vars["keyword"], &groupMessages)
		groupMessages = parsedMessages
	}
	common.ErrorHandler(w, err)
	data, jsonErr := json.Marshal(&groupMessages)
	common.ErrorHandler(w, jsonErr)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func ParseGroupMessages(keyword string, groupMessages *groupme_api.GroupMessagesResponse) (parsedGroupMessages groupme_api.GroupMessagesResponse) {
	messages := groupMessages.Response.Messages
	parsedMessages := []groupme_api.Message{}
	for _, message := range messages {
		text := message.Text
		if strings.Contains(strings.ToLower(text), strings.ToLower(keyword)) {
			parsedMessages = append(parsedMessages, message)
		}
	}

	parsedGroupMessages.Response.Messages = parsedMessages
	return
}
