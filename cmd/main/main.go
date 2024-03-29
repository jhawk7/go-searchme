package main

import (
	"context"
	"encoding/json"
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
	Filter string `form:"filter"`
}

func CheckCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		var params Params
		if err := c.ShouldBind(&params); err != nil {
			c.Next()
		}

		value, cacheMiss, cacheErr := cacheClient.GetValue(c, params.Filter)
		if cacheErr != nil {
			err := fmt.Errorf("failed to retrieve value from cache; %v", cacheErr)
			common.ErrorHandler(err, false)
			c.Next()
			return
		}

		if cacheMiss {
			c.Next()
			return
		}

		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"messages": value,
		})
	}
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
	router.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "up",
		})
	})

	v1 := router.Group("/v1")
	v1.Use(CheckCache())
	{
		v1.GET("/deals", GetFlightDeals)
	}

	router.Run(":8888")
}

func GetFlightDeals(c *gin.Context) {
	var par Params

	if c.ShouldBind(&par) != nil {
		common.ErrorHandler(fmt.Errorf("failed to bind query params"), false)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
		return
	}

	combinedMessages, retrieveErr := retrieveMessages(c)
	if retrieveErr != nil {
		common.ErrorHandler(retrieveErr, false)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
		return
	}

	if par.Filter != "" {
		filterGroupMessages(par.Filter, combinedMessages, true)
	}
	storeMessages(par.Filter, combinedMessages)

	c.JSON(http.StatusOK, gin.H{
		"messages": *combinedMessages,
		"count":    len(*combinedMessages),
	})
}

func retrieveMessages(ctx context.Context) (combinedMessages *[]groupme.Message, err error) {
	// check cache for stored messages
	cachedMessages, cacheMiss, cacheErr := cacheClient.GetValue(ctx, "messages")
	if cacheErr != nil {
		common.ErrorHandler(cacheErr, false)
	}

	if !cacheMiss {
		bytes, _ := json.Marshal(cachedMessages)
		err = json.Unmarshal(bytes, combinedMessages)
		return
	}

	// retrieves last 200 messages via 2 API calls
	var messages []groupme.Message
	var offset *string
	for i := 0; i < 2; i++ {
		groupMessages, firstMessageId, groupErr := gmClient.GetGroupMessages(offset)
		if groupErr != nil {
			err = groupErr
			return
		}
		messages = append(messages, groupMessages.Response.Messages...)
		offset = &firstMessageId
	}

	storeMessages("messages", &messages)
	combinedMessages = &messages
	return
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
