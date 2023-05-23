package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jhawk7/go-searchme/pkg/groupme"
	log "github.com/sirupsen/logrus"
)

var gmClient *groupme.Client

func main() {
	gmClient = groupme.InitClient()
	router := gin.Default()
	router.GET("/healthcheck", HealthCheck)
	router.GET("/groupme/groups", GetUserGroups)
	router.GET("/groupme/group/messages", GetGroupMessages)
	router.GET("/groupme/group/messages/:keyword", GetGroupMessages)
	router.Run(":8888")
}

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "up",
	})
}

func GetUserGroups(c *gin.Context) {
	groups, groupsErr := gmClient.GetUserGroups()
	if groupsErr != nil {
		ErrorHandler(groupsErr, false)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad reqeust",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": groups,
	})
}

func GetGroupMessages(c *gin.Context) {
	keyword := c.Param("keyword")
	groupMessages, groupErr := gmClient.GetGroupMessages()
	if groupErr != nil {
		ErrorHandler(groupErr, false)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
		return
	}

	if keyword != "" {
		filterGroupMessages(keyword, &groupMessages)
	}

	//return text from message only
	textMessages := []string{}
	for _, message := range groupMessages.Response.Messages {
		textMessages = append(textMessages, message.Text)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": textMessages,
	})
}

func filterGroupMessages(keyword string, groupMessages *groupme.GroupMessagesResponse) {
	messages := groupMessages.Response.Messages
	parsedMessages := []groupme.Message{}
	for _, message := range messages {
		if strings.Contains(strings.ToLower(message.Text), strings.ToLower(keyword)) {
			parsedMessages = append(parsedMessages, message)
		}
	}

	groupMessages.Response.Messages = parsedMessages
}

func ErrorHandler(err error, fatal bool) {
	if err != nil {
		log.Errorf("error: %v", err)

		if fatal {
			panic(err)
		}
	}
}
