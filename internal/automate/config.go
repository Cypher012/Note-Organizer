package automate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/brianvoe/gofakeit/v6"
)

type TestUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewTestUser(t bool) TestUser {
	if t {
		return TestUser{
			Username: gofakeit.Username(),
			Email:    gofakeit.Email(),
			Password: "Cipher2017",
		}
	}

	return TestUser{
		Username: "cipher",
		Email:    "ayoojoade@gmail.com",
		Password: "Cipher2017",
	}
}

func strPtr(s string) *string { return &s }

func DoJSONRequest(
	client *http.Client,
	method string,
	url string,
	payload any,
) ([]byte, error) {
	var body io.Reader

	if payload != nil {
		b, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(b)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Println("STATUS:", resp.Status)
	log.Println("BODY:", string(respBody))

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("request failed: %s", resp.Status)
	}

	return respBody, nil
}
