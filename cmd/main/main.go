package main

import (
	"net/http"
	"strings"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jhawk7/go-searchme/pkg/groupme"
	log "github.com/sirupsen/logrus"
	xurls "mvdan.cc/xurls/v2"
)

var gmClient *groupme.Client

func main() {
	gmClient = groupme.InitClient()
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.GET("/healthcheck", HealthCheck)
	router.GET("/groupme/groups", GetUserGroups)
	router.GET("/groupme/group/messages", GetGroupMessages)
	router.GET("/groupme/group/messages/:keyword", GetGroupMessages)
	router.GET("/flights/:keyword", DisplayFlightDeals)
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
	groupMessages, _, groupErr := gmClient.GetGroupMessages(nil)
	if groupErr != nil {
		ErrorHandler(groupErr, false)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
		return
	}

	if keyword != "" {
		filterGroupMessages(keyword, &groupMessages, false)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": groupMessages.Response.Messages,
	})
}

func DisplayFlightDeals(c *gin.Context) {
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
		filterGroupMessages(keyword, &groupMessages, false)
	}

	// Render HTML template
	c.HTML(http.StatusOK, "flight_deals.tmpl", gin.H{
		"Messages": groupMessages.Response.Messages,
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

func addHyperlinks(message *groupme.Message ) {
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
