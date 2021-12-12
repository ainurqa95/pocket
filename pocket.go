package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const consumerKey = "99907-059c27eeb9c6dbfd43d936f0"
const (
	host            = "https://getpocket.com/v3/"
	requestTokenUri = "oauth/request"
	authorizeUrl    = "https://getpocket.com/auth/authorize?request_token=%s&redirect_uri=%s"
	authUri         = "oauth/authorize"
)

type PocketMananger interface {
	getRequestToken(redirectUri string) (string, error)
}

type PocketClient struct {
	consumerKey string
}

func NewPocketClient(consumerKey string) *PocketClient {
	return &PocketClient{consumerKey: consumerKey}
}

func (client *PocketClient) getRequestToken(redirectUri string) (string, error) {
	requestUrl := host + requestTokenUri

	body := map[string]string{
		"consumer_key": client.consumerKey,
		"redirect_uri": redirectUri,
	}

	codeBody, err := client.makeRequest(requestUrl, body)
	if codeBody["code"] == "" {
		return "", errors.New("empty request token code in API response")
	}

	return codeBody["code"], err
}

func (client *PocketClient) getAuthorizationUrl(tokenCode string, redirectUri string) (string, error) {
	if tokenCode == "" || redirectUri == "" {
		return "", errors.New("empty params")
	}

	return fmt.Sprintf(authorizeUrl, tokenCode, redirectUri), nil
}

type AccesTokenResponse struct {
	accessToken string
	username    string
}

func (client *PocketClient) authAndFinalAccessToken(tokenCode string) (*AccesTokenResponse, error) {
	reqUrl := host + authUri
	body := map[string]string{
		"consumer_key": client.consumerKey,
		"code":         tokenCode,
	}
	accessTokenResp, err := client.makeRequest(reqUrl, body)
	if err != nil {
		return nil, err
	}
	if accessTokenResp["access_token"] == "" {
		return nil, errors.New("empty access token in API response")
	}
	if accessTokenResp["username"] == "" {
		return nil, errors.New("empty username")
	}

	return &AccesTokenResponse{
		accessToken: accessTokenResp["access_token"],
		username:    accessTokenResp["username"],
	}, nil

}

func (client *PocketClient) makeRequest(url string, reqBody map[string]string) (map[string]string, error) {

	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}
	jsonReqBody, _ := json.Marshal(reqBody)
	reqBytesBody := bytes.NewBuffer(jsonReqBody)
	req, err := http.NewRequest("POST", url, reqBytesBody)
	if err != nil {
		return make(map[string]string), err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("X-Accept", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		return make(map[string]string), nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return make(map[string]string), err
	}
	respBody := string(body)
	var resultBody map[string]string
	json.Unmarshal([]byte(respBody), &resultBody)

	return resultBody, nil
}

func main() {
	client := NewPocketClient(consumerKey)
	token, err := client.getRequestToken("pocketapp1234:authorizationFinished")
	fmt.Println(token, err)
	authUrl, err := client.getAuthorizationUrl(token, "https://example.com")
	fmt.Println(authUrl, err)
	time.Sleep(60 * time.Second)
	// only after authorization
	accessTokenRes, err := client.authAndFinalAccessToken(token)
	fmt.Println(accessTokenRes, err)

}
