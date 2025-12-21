package automate

import (
	"log"
	"net/url"
)

var BaseApi = "http://localhost:8080/api"

var UrlMap = map[string]string{
	"register":      formatUrl("/auth/register"),
	"login":         formatUrl("/auth/login"),
	"get-folders":   formatUrl("/folder"),
	"create-folder": formatUrl("/folder"),
}

func formatUrl(endpoint string) string {
	u, _ := url.JoinPath(BaseApi, endpoint)
	log.Println(u)
	return u
}
