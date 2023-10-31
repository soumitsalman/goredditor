package redditclient

import (
	"context"
	"log"
	"net/http"

	"angerproject.org/redditor/utils"
)

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
