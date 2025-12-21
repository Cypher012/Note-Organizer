package automate

import (
	"log"
	"net/http"
	"net/http/cookiejar"
)

func CreateClient() *http.Client {
	// 1. Initialize the cookie jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		// Using log.Fatal is often preferred over panic for CLI/services
		log.Fatal(err)
	}

	// 2. Create the client with the jar attached
	client := &http.Client{
		Jar: jar,
	}

	return client
}
