# GetPocket API Golang SDK

### Example usage:

```go
package main

import (
	pocket "github.com/ainurqa95/pocket"
)

func main() {
	client := NewPocketClient(consumerKey)
	token, err := client.GetRequestToken("https://example.com")
	fmt.Println(token, err)
	authUrl, err := client.DefineAuthorizationUrl(token, "https://example.com")
	fmt.Println(authUrl, err)
	time.Sleep(60 * time.Second)
	// only after authorization
	accessTokenRes, err := client.AuthAndGetAccessToken(token)
	fmt.Println(accessTokenRes, err)
	err = client.AddItem(AddInput{
		accessToken: accessTokenRes.accessToken,
		url:         "https://github.com",
	})
	fmt.Println(err)
}

```