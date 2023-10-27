package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
)

type Task struct {
	UserId     int    `json:"userId"`
	TaskId     int    `json:"id"`
	Title      string `json:"title"`
	IsComplete bool   `json:"completed"`
}

func main() {

	/*
		const CLIENT_ID = "fBny617OCBBOQTBgsrMnxg"
		const CLIENT_SECRET = "55Fb5NhxI9Uc8Xze9vOUom0FuyAiUQ"
		const USER_NAME = "randomizer_000"
		const PASSWORD = "4EvjH^D&R5CZtIBN"

		var redditAuthConfig = &oauth2.Config{
			ClientID:     CLIENT_ID,
			ClientSecret: CLIENT_SECRET,
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://www.reddit.com/api/v1/authorize",
				TokenURL: "https://www.reddit.com/api/v1/access_token",
			},
		}

		var ctx = context.Background()
		var redditAuthToken, authErr = redditAuthConfig.PasswordCredentialsToken(
			ctx,
			USER_NAME,
			PASSWORD)

		if authErr != nil {
			fmt.Println("Auth Failed: ", authErr)
			return
		} else {
			fmt.Println("Auth worked: ", redditAuthToken)
		}

		var redditClient = redditAuthConfig.Client(ctx, redditAuthToken)
		var req, _ = http.NewRequestWithContext(ctx, "GET", "https://oauth.reddit.com/api/v1/me", nil)
		req.Header.Set("User-Agent", "script:https://github.com/soumitsalman/goredditor#readme:v0.1 (by /u/randomizer_000)")
		var resp, reqErr = redditClient.Do(req)
		if reqErr != nil {
			fmt.Println("Can't find me: ", reqErr)
			return
		} else {
			fmt.Println("X-Ratelimit-Remaining: ", resp.Header.Get("X-Ratelimit-Remaining"))
		}

		defer resp.Body.Close()


			var t Task
			var dec_err = json.NewDecoder(resp.Body).Decode(&t)
			if dec_err != nil {
				fmt.Println(dec_err)
			} else {
				fmt.Println(t)
			}

		var bodyString = new(strings.Builder)
		io.Copy(bodyString, resp.Body)
		fmt.Println(bodyString.String())
	*/
	var config, err = loadConfigurationFromFile("config.json")
	if err != nil {
		log.Println(err)
		return
	}
	var session = CreateNewSession(&config)
	var meData, meErr = session.GetMe()

	if meErr != nil {
		log.Println("Failed to load me: ", meErr)
	} else {
		fmt.Println(meData["name"], " has ", meData["total_karma"], " karma")
	}

}

type RedditAccessConfiguration struct {
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

type RedditSession struct {
	configuration *RedditAccessConfiguration
	client        *http.Client
	ctx           context.Context
}

func CreateNewSession(configuration *RedditAccessConfiguration) RedditSession {
	//TODO: double check if the LastAccessToken is still valid
	return RedditSession{
		configuration: configuration,
		client:        &http.Client{},
		ctx:           context.Background(),
	}
}

func (session *RedditSession) GetMe() (map[string]any, error) {
	var req = session.getHttpRequestHolder("GET", "https://oauth.reddit.com/api/v1/me")
	var resp, err = session.client.Do(req)
	if err != nil {
		log.Println("Getting me failed: ", err)
		return nil, err
	}
	defer resp.Body.Close()

	return deserialzeJsonBlob[map[string]any](io.Reader(resp.Body))
}

func (configuration *RedditAccessConfiguration) getUserAgentName() string {
	//Windows:My Reddit Bot:1.0 (by u/botdeveloper)
	return fmt.Sprintf("%v:%v:v0.1 (by u/%v)", runtime.GOOS, configuration.ClientName, configuration.Username)
}

func (configuration *RedditAccessConfiguration) getAuthorizationToken() string {
	return "bearer " + configuration.LastAccessToken
}

func (session *RedditSession) getHttpRequestHolder(method string, url string) *http.Request {
	var req, _ = http.NewRequestWithContext(session.ctx, method, url, nil)
	//standard header
	req.Header.Add("User-Agent", session.configuration.getUserAgentName())
	req.Header.Add("Authorization", session.configuration.getAuthorizationToken())
	return req
}

func loadConfigurationFromFile(configFilePath string) (RedditAccessConfiguration, error) {
	var configFile, err = os.Open(configFilePath)
	if err != nil {
		log.Printf("Failed loading configuration file %v. Error: %v\n", configFilePath, err)
		return RedditAccessConfiguration{}, err
	}
	defer configFile.Close()

	return deserialzeJsonBlob[RedditAccessConfiguration](configFile)
}

func deserialzeJsonBlob[T any](reader io.Reader) (T, error) {
	var decoder = json.NewDecoder(reader)
	var data T
	var err = decoder.Decode(&data)
	if err != nil {
		log.Printf("Error deserializing to data of type %T: %v\n", data, err)
		return data, err
	} else {
		return data, nil
	}
}
