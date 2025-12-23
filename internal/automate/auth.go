package automate

import (
	"log"
	"net/http"
)

func Register(client *http.Client, user TestUser) error {
	url := UrlMap["register"]()

	_, err := DoJSONRequest(
		client,
		http.MethodPost,
		url,
		user,
	)
	if err != nil {
		return err
	}

	log.Println("Register successfully")
	return nil
}

func Login(client *http.Client, user TestUser) error {
	url := UrlMap["login"]()

	loginData := map[string]string{
		"email":    user.Email,
		"password": user.Password,
	}

	_, err := DoJSONRequest(
		client,
		http.MethodPost,
		url,
		loginData,
	)
	if err != nil {
		return err
	}

	log.Println("Login successfully")
	return nil
}
