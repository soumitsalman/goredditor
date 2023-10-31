package redditclient

import (
	"context"
	"net/http"
)

type RedditCredentials struct {
	LastAccessToken   string `json:"access_token"`
	ClientSecret      string `json:"client_secret"`
	ClientName        string `json:"client_name"`
	ClientId          string `json:"client_id"`
	ClientDescription string `json:"client_description"`
	AboutUrl          string `json:"about_url"`
	RedirectUri       string `json:"redirect_uri"`
	Username          string `json:"user_name"`
	Password          string `json:"user_password"`
}

type RedditClient struct {
	creds      *RedditCredentials
	httpClient *http.Client
	ctx        context.Context
}
