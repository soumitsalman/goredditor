package redditapplicationclient

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"runtime"

	"angerproject.org/redditor/utils"
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

func NewClientFromConfigFile(config_file string) (RedditClient, error) {
	if config, err := utils.ReadDataFromJsonFile[RedditCredentials](config_file); err != nil {
		log.Println(err)
		return RedditClient{}, err
	} else {
		return NewClient(&config), nil
	}
}

func NewClient(creds *RedditCredentials) RedditClient {
	//TODO: double check if the LastAccessToken is still valid
	return RedditClient{
		creds:      creds,
		httpClient: &http.Client{},
		ctx:        context.Background(),
	}
}

func SaveClientToConfigFile(client *RedditClient, config_file string) error {
	if err := utils.WriteDataToJsonFile[RedditCredentials](client.creds, config_file); err != nil {
		log.Printf("Saving auth token failed: %v\n", err)
		return err
	} else {
		log.Println("Configuration file saved with new auth token")
		return nil
	}
}

func (client *RedditClient) getApplicationFullName() string {
	//Windows:My Reddit Bot:1.0 (by u/botdeveloper)
	return fmt.Sprintf("%v:%v:v0.1 (by u/%v)", runtime.GOOS, client.creds.ClientName, client.creds.Username)
}
