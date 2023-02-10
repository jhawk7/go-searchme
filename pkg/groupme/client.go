package groupme

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/pkg/errors"
)

const GetGroupsEndpoint string = "https://api.groupme.com/v3/groups?omit=memberships"
const GetMessagesEndpoint string = "https://api.groupme.com/v3/groups/%v/messages?since_id=%v&limit=100"

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

func (client *Client) GetUserGroups() (groupsResponse GroupsResponse, respErr error) {
	url := GetGroupsEndpoint

	request, _ := retryablehttp.NewRequest("GET", url, nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Access-Token", client.Token)
	response, requestErr := client.RetryClient.Do(request)

	if requestErr != nil {
		respErr = fmt.Errorf("failed to make request for user groups; [error: %v]", requestErr)
		return
	}

	respBody, _ := ioutil.ReadAll(response.Body)
	if response.StatusCode > 299 {
		respErr = fmt.Errorf("non-200 status received from user group response; [status: %v] [body: %v]", response.StatusCode, string(respBody))
		return
	}

	//parse response
	groupsResponse = GroupsResponse{}
	parseErr := json.Unmarshal(respBody, &groupsResponse)
	if parseErr != nil {
		respErr = fmt.Errorf("failed to parse user group response; [error: %v]", parseErr)
		return
	}

	return
}

func (client *Client) GetGroupMessages() (messagesResponse GroupMessagesResponse, respErr error) {
	url := fmt.Sprintf(GetMessagesEndpoint, client.GroupId, client.SinceId)

	request, _ := retryablehttp.NewRequest("GET", url, nil)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Access-Token", client.Token)
	response, requestErr := client.RetryClient.Do(request)

	if requestErr != nil {
		respErr = fmt.Errorf("failed to get group messages; [error: %v]", requestErr)
		return
	}

	respBody, _ := ioutil.ReadAll(response.Body)
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

	return
}
