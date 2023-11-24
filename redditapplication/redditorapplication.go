package redditapplication

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
)

type RedditorCredentials struct {
	OauthToken             string
	ApplicationSecret      string
	ApplicationName        string
	ApplicationId          string
	ApplicationDescription string
	AboutUrl               string
	RedirectUri            string
	Username               string
	Password               string
}

type RedditorApplication struct {
	creds      *RedditorCredentials
	httpClient *http.Client
	ctx        context.Context
}

func NewClient(creds *RedditorCredentials) RedditorApplication {
	//TODO: double check if the LastAccessToken is still valid
	return RedditorApplication{
		creds:      creds,
		httpClient: &http.Client{},
		ctx:        context.Background(),
	}
}

func (client *RedditorApplication) getApplicationFullName() string {
	//Windows:My Reddit Bot:1.0 (by u/botdeveloper)
	return fmt.Sprintf("%v:%v:v0.1 (by u/%v)", runtime.GOOS, client.creds.ApplicationName, client.creds.Username)
}
