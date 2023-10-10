package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
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
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET"},
		AllowHeaders: []string{"Origin, Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"},
	}))
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
	var combinedMessages []groupme.Message
	var offset *string
	var err error

	if offsetParam := c.DefaultQuery("offset", ""); offsetParam == "" {
		offset = nil
	} else {
		offset = &offsetParam
	}

	// retrieves last 200 messages via 2 API calls
	for i := 0; i < 2; i++ {
		groupMessages, firstMessageId, groupErr := gmClient.GetGroupMessages(offset)
		if groupErr != nil {
			err = groupErr
			break
		}
		combinedMessages = append(combinedMessages, groupMessages.Response.Messages...)
		offset = &firstMessageId
	}

	if err != nil {
		ErrorHandler(err, false)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
		return
	}

	if keyword != "" {
		filterGroupMessages(keyword, &combinedMessages, true)
	}

	//cache offset:messages

	c.JSON(http.StatusOK, gin.H{
		"messages": combinedMessages,
		"offset":   *offset,
	})
}

func filterGroupMessages(keyword string, groupMessages *[]groupme.Message, highlightLinks bool) {
	parsedMessages := []groupme.Message{}
	for _, message := range *groupMessages {
		if strings.Contains(strings.ToLower(message.Text), strings.ToLower(keyword)) {
			if highlightLinks {
				addHyperlinks(&message)
			}
			parsedMessages = append(parsedMessages, message)
		}
	}

	*groupMessages = parsedMessages
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
