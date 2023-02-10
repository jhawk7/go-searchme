package groupme

import "github.com/hashicorp/go-retryablehttp"

type Client struct {
	Token       string
	GroupId     string
	SinceId     string
	RetryClient *retryablehttp.Client
}

type GroupsResponse struct {
	Response []Group `json:"response"`
}

type Group struct {
	Id            string      `json:"id"`
	GroupId       string      `json:"group_id"`
	Name          string      `json:"name"`
	PhoneNumber   string      `json:"phone_number"`
	Type          string      `json:"type"`
	Description   string      `json:"description"`
	ImageUrl      string      `json:"image_url"`
	CreatorUserId string      `json:"creator_user_id"`
	CreatedAt     int         `json:"created_at"`
	UpdatedAt     int         `json:"update_at"`
	LastMessage   LastMessage `json:"messages"`
	MaxMembers    int         `json:"max_members"`
}

type LastMessage struct {
	Count                int     `json:"count"`
	LastMessageId        string  `json:"last_message_id"`
	LastMessageCreatedAt int     `json:"last_message_created_at"`
	Preview              Preview `json:"preview"`
}

type Preview struct {
	Nickname string `json:"nickname"`
	Text     string `json:"text"`
	ImageUrl string `json:"image_url"`
}

type GroupMessagesResponse struct {
	Response Response `json:"response"`
}

type Response struct {
	Messages []Message `json:"messages"`
}

type Message struct {
	AvatarUrl  string `json:"avatar_url"`
	CreatedAt  int    `json:"created_at"`
	GroupId    string `json:"group_id"`
	Id         string `json:"id"`
	Name       string `json:"name"`
	SenderId   string `json:"sender_id"`
	SenderType string `json:"sender_type"`
	SourceGUID string `json:"source_guid"`
	System     bool   `json:"system"`
	Text       string `json:"text"`
	UserId     string `json:"user_id"`
	Platform   string `json:"platform"`
}
