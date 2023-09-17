package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/jhawk7/go-searchme/pkg/groupme"
	log "github.com/sirupsen/logrus"
	xurls "mvdan.cc/xurls/v2"
)

var gmClient *groupme.Client

func main() {
	gmClient = groupme.InitClient()
	router := gin.Default()
	router.Use(static.Serve("/", static.LocalFile("./frontend/dist", false)))
	router.GET("/healthcheck", HealthCheck)
	router.GET("/flights/:keyword", GetFlightDeals)
	router.Run(":8888")
}

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "up",
	})
}

func GetFlightDeals(c *gin.Context) {
	keyword := c.Param("keyword")
	groupMessages, firstMessageId, groupErr := gmClient.GetGroupMessages(nil)
	if groupErr != nil {
		ErrorHandler(groupErr, false)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
		return
	}

	//paginated call
	groupMessages2, _, groupErr2 := gmClient.GetGroupMessages(&firstMessageId)
	if groupErr2 != nil {
		ErrorHandler(groupErr2, false)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
		return
	}

	//combine group message responses
	groupMessages.Response.Messages = append(groupMessages.Response.Messages, groupMessages2.Response.Messages...)

	if keyword != "" {
		filterGroupMessages(keyword, &groupMessages, true)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": groupMessages.Response.Messages,
	})
}

func filterGroupMessages(keyword string, groupMessages *groupme.GroupMessagesResponse, highlightLinks bool) {
	messages := groupMessages.Response.Messages
	parsedMessages := []groupme.Message{}
	for _, message := range messages {
		if strings.Contains(strings.ToLower(message.Text), strings.ToLower(keyword)) {
			if highlightLinks {
				addHyperlinks(&message)
			}
			parsedMessages = append(parsedMessages, message)
		}
	}

	groupMessages.Response.Messages = parsedMessages
}

func addHyperlinks(message *groupme.Message) {
	urls := xurls.Strict().FindAllString(message.Text, -1)
	for _, url := range urls {
		hypertext := strings.ReplaceAll(message.Text, url, fmt.Sprintf(`<a href="%v">%v</a>`, url, url))
		message.Text = hypertext
	}
}

func ErrorHandler(err error, fatal bool) {
	if err != nil {
		log.Errorf("error: %v", err)

		if fatal {
			panic(err)
		}
	}
}
