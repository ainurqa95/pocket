package pocket

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	host            = "https://getpocket.com/v3/"
	requestTokenUri = "oauth/request"
	authorizeUrl    = "https://getpocket.com/auth/authorize?request_token=%s&redirect_uri=%s"
	authUri         = "oauth/authorize"
	addUri          = "add"
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

func (client *PocketClient) GetRequestToken(redirectUri string) (string, error) {
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

func (client *PocketClient) DefineAuthorizationUrl(tokenCode string, redirectUri string) (string, error) {
	if tokenCode == "" || redirectUri == "" {
		return "", errors.New("empty params")
	}

	return fmt.Sprintf(authorizeUrl, tokenCode, redirectUri), nil
}

type AccesTokenResponse struct {
	AccessToken string
	Username    string
}

func (client *PocketClient) AuthAndGetAccessToken(tokenCode string) (*AccesTokenResponse, error) {
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
		AccessToken: accessTokenResp["access_token"],
		Username:    accessTokenResp["username"],
	}, nil

}

type AddInput struct {
	Url         string
	AccessToken string
	Title       string
	Tags        []string
	Tweet_id    string
}

func (addInput *AddInput) validate() error {
	if addInput.AccessToken == "" {
		return errors.New("Add input request token is empty")
	}
	if addInput.Url == "" {
		return errors.New("Add input url is empty")
	}

	return nil
}

func (client *PocketClient) AddItem(addInput AddInput) error {
	err := addInput.validate()
	if err != nil {
		return err
	}
	reqUrl := host + addUri
	body := map[string]string{
		"consumer_key": client.consumerKey,
		"access_token": addInput.AccessToken,
		"title":        addInput.Title,
		"url":          addInput.Url,
		"tags":         strings.Join(addInput.Tags, ", "),
		"tweet_id":     addInput.Tweet_id,
	}

	_, err = client.makeRequest(reqUrl, body)

	return err
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
	if resp.StatusCode != 200 {
		return make(map[string]string), errors.New("Err status code :" + resp.Status)
	}
	var resultBody map[string]string
	json.Unmarshal([]byte(respBody), &resultBody)

	return resultBody, nil
}
