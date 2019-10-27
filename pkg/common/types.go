package common

type Response struct {
	Status int    `json:"status"`
	Body   string `json:"body"`
}

type Config struct {
	Token string
	SinceId string
	GroupId string
}