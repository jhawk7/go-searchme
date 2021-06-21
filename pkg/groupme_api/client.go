package groupme_api

import (
	"encoding/json"
	"fmt"
	"github.com/jhawk7/go-searchme/pkg/common"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

const GetGroupsEndpoint string = "https://api.groupme.com/v3/groups?omit=memberships"
const GetMessagesEndpoint string = "https://api.groupme.com/v3/groups/%v/messages?since_id=%v&limit=100"

func CreateGroupmeClient(config *common.Config) (*GroupmeClient) {
	client := GroupmeClient{
		Token: config.Token,
		GroupId: config.GroupId,
		SinceId: config.SinceId,
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
	retryClient.ErrorHandler = ErrorHandler

	return retryClient
}

var RequestHook = func(logger retryablehttp.Logger, req *http.Request, retryNumber int) {
	fmt.Printf("Making Request to URL: %s; Retry Count: %v\n", req.URL, retryNumber)
}

var ResponseHook = func(logger retryablehttp.Logger, res *http.Response) {
	fmt.Printf("URL: %v responded with Status: %v\n", res.Request.URL, res.StatusCode)
}

//Response will be nil if connection error occurs; with ErrorHandler we can still get the response if the endpoint responds with 500
var ErrorHandler = func(resp *http.Response, err error, numTries int) (*http.Response, error) {
	if resp != nil {
		fmt.Printf("Error. [Status: %v] [Error: %v] [Tries: %v]", resp.StatusCode, err, numTries)
	} else {
		fmt.Printf("Error. [Status: NO RESPONSE OBJECT] [Error: %v] [Tries: %v]", err, numTries)

	}
	return resp, err
}

func (client *GroupmeClient) GetUserGroups() (groupsResponse GroupsResponse, respErr error) {
	url := GetGroupsEndpoint

	request, _ := retryablehttp.NewRequest("GET", url, nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Access-Token", client.Token)
	response, requestErr := client.RetryClient.Do(request)
	respBody, _ := ioutil.ReadAll(response.Body)
	if requestErr != nil {
		respErr = errors.New(fmt.Sprintf("[ERROR: %v], [RESPONSE: %v]", requestErr, respBody))
		return
	}

	//parse response
	groupsResponse = GroupsResponse{}
	parseErr := json.Unmarshal(respBody, &groupsResponse)
	if parseErr != nil {
		respErr = errors.New(fmt.Sprintf("Error parsing GetUserGroups response. [ERROR: %v]", parseErr))
		return
	}

	return
}

func (client *GroupmeClient) GetGroupMessages() (messagesResponse GroupMessagesResponse, respErr error) {
	url := fmt.Sprintf(GetMessagesEndpoint, client.GroupId, client.SinceId)

	request, _ := retryablehttp.NewRequest("GET", url, nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Access-Token", client.Token)
	response, requestErr := client.RetryClient.Do(request)
	respBody, _ := ioutil.ReadAll(response.Body)
	if requestErr != nil {
		respErr = errors.New(fmt.Sprintf("[ERROR: %v], [RESPONSE: %v]", requestErr, respBody))
		return
	}

	//parse response
	messagesResponse = GroupMessagesResponse{}
	parseErr := json.Unmarshal(respBody, &messagesResponse)
	if parseErr != nil {
		respErr = errors.New(fmt.Sprintf("Error parsing GetUserGroups response. [ERROR: %v]", parseErr))
		return
	}

	return

}