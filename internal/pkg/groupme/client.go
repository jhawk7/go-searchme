package groupme

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/pkg/errors"
)

const GetGroupsEndpoint string = "https://api.groupme.com/v3/groups?omit=memberships"
const GetMessagesSinceEndpoint string = "https://api.groupme.com/v3/groups/%v/messages?since_id=%v&limit=100"
const GetMessagesBeforeEndpoint string = "https://api.groupme.com/v3/groups/%v/messages?before_id=%v&limit=100"

func InitClient() *Client {
	client := Client{
		Token:       os.Getenv("TOKEN"),
		GroupId:     os.Getenv("GROUP_ID"),
		SinceId:     os.Getenv("SINCE_ID"),
		RetryClient: CreateRetryClient(),
	}

	return &client
}

func CreateRetryClient() (retryClient *retryablehttp.Client) {
	//init Retryable HTTP Client
	retryClient = retryablehttp.NewClient()
	retryClient.RetryMax = 3
	retryClient.RequestLogHook = RequestHook
	retryClient.ResponseLogHook = ResponseHook

	return retryClient
}

var RequestHook = func(logger retryablehttp.Logger, req *http.Request, retryNumber int) {
	fmt.Printf("Making Request to URL: %s; Retry Count: %v\n", req.URL, retryNumber)
}

var ResponseHook = func(logger retryablehttp.Logger, res *http.Response) {
	fmt.Printf("URL: %v responded with Status: %v\n", res.Request.URL, res.StatusCode)
}

func (client *Client) GetGroupMessages(beforeId *string) (messagesResponse GroupMessagesResponse, firstMessageId string, respErr error) {
	var url string
	// make api call based on most recent (100) messages after given sinceID
	// make subsequent api calls based on messages (100) before given beforeId if present
	if beforeId == nil {
		url = fmt.Sprintf(GetMessagesSinceEndpoint, client.GroupId, client.SinceId)
	} else {
		url = fmt.Sprintf(GetMessagesBeforeEndpoint, client.GroupId, *beforeId)
	}

	request, _ := retryablehttp.NewRequest("GET", url, nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Access-Token", client.Token)
	response, requestErr := client.RetryClient.Do(request)

	if requestErr != nil {
		respErr = fmt.Errorf("failed to get group messages; [error: %v]", requestErr)
		return
	}

	respBody, _ := io.ReadAll(response.Body)
	if response.StatusCode > 299 {
		respErr = fmt.Errorf("non-200 status received from user group response; [status: %v] [body: %v]", response.StatusCode, string(respBody))
		return
	}

	//parse response
	messagesResponse = GroupMessagesResponse{}
	parseErr := json.Unmarshal(respBody, &messagesResponse)
	if parseErr != nil {
		respErr = errors.New(fmt.Sprintf("Error parsing GetUserGroups response. [ERROR: %v]", parseErr))
		return
	}

	//can be used as before ID subsequent call to get earlier messages
	firstMessageId = messagesResponse.Response.Messages[0].Id
	return
}
