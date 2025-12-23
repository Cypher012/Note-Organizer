package automate

import (
	"log"
	"net/url"
)

var BaseApi = "http://localhost:8080/api"

var UrlMap = map[string]func(...string) string{
	"register": func(_ ...string) string {
		return BuildURL("auth", "register")
	},
	"login": func(_ ...string) string {
		return BuildURL("auth", "login")
	},
	"folders": func(args ...string) string {
		parts := append([]string{"folders"}, args...)
		return BuildURL(parts...)
	},
	"get-notes": func(_ ...string) string {
		return BuildURL("notes")
	},
}

func BuildURL(parts ...string) string {
	u, err := url.JoinPath(BaseApi, parts...)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(u)
	return u
}
