package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/jhawk7/go-searchme/internal/pkg/cache"
	"github.com/jhawk7/go-searchme/internal/pkg/common"
	"github.com/jhawk7/go-searchme/internal/pkg/groupme"
	xurls "mvdan.cc/xurls/v2"
)

var gmClient *groupme.Client
var cacheClient *cache.RedisClient

type Params struct {
	Filter   string `form:"filter"`
	Page     string    `form:"page"`
	PageSize string    `form:"pageSize"`
}

func main() {
	gmClient = groupme.InitClient()
	cacheClient = cache.InitClient()
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET"},
		AllowHeaders: []string{"Origin, Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"},
	}))

	router.Use(static.Serve("/", static.LocalFile("./frontend/dist", false)))
	router.GET("/healthcheck", HealthCheck)
	router.GET("/flights/:keyword", GetFlightDeals)
	router.GET("/v1/deals", GetFlightDeals) // v1/deals?filter=&page=&pageSize=20
	router.Run(":8888")
}

func CheckCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		var params Params
		if err := c.ShouldBind(&params); err != nil {
			c.Next()
		}
		cacheClient.GetValue(c, params.Filter)
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"messages":,
		})
	}
}

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "up",
	})
}

func GetFlightDeals(c *gin.Context) {
	var params Params
	var combinedMessages []groupme.Message
	var err error

	if c.ShouldBind(&params) == nil {
		common.ErrorHandler(fmt.Errorf("failed to bind query params"), false)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
		return
	}

	if hit, := checkCache(params); hit {

	}

	// retrieves last 200 messages via 2 API calls
	var offset *string
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
		common.ErrorHandler(err, false)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
		return
	}

	storeMessages("messages", &combinedMessages)

	if params.Filter != "" {
		filterGroupMessages(params.Filter, &combinedMessages, true)
	}

	storeMessages(params.Filter, &combinedMessages)

	c.JSON(http.StatusOK, gin.H{
		"messages": combinedMessages,
	})
}

func storeMessages(key string, messages *[]groupme.Message) {
	kv := cache.KVPair{
		Key:   key,
		Value: *messages,
	}

	if err := cacheClient.Store(context.Background(), kv); err != nil {
		common.ErrorHandler(fmt.Errorf("failed to cache messsages [key: %v] [err: %v]", key, err), false)
	}
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
