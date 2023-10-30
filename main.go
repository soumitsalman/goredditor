package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

func main() {
	var config, err = loadCredentialsFromFile("config.json")
	if err != nil {
		log.Println(err)
		return
	}
	var client = NewClient(&config)

	if is_new_token, err := client.authenticate(); err != nil {
		log.Printf("Auth failed: %v\n", err)
		return
	} else if is_new_token {
		log.Printf("Got new auth token: \n")
		if wr_err := writeDataToJsonFile[RedditCredentials](client.creds, "config.json"); wr_err != nil {
			log.Printf("Saving auth token failed: %v\n", wr_err)
		} else {
			log.Println("Configuration file saved with new auth token")
		}
	}

	if post_collection, err := client.GetHotPosts(""); err != nil {
		log.Printf("Getting subreddits failed: %v\n", err)
	} else {
		if writeDataToJsonFile[[]map[string]any](&post_collection, "hot_response.json") == nil {
			log.Println("Saved hot post results in hot_response.json")
		} else {
			log.Println("Failed to save hot posts")
		}
	}

	if sr_collection, err := client.GetSubreddits(); err != nil {
		log.Printf("Getting subreddits failed: %v\n", err)
	} else {
		if writeDataToJsonFile[[]map[string]any](&sr_collection, "subreddit_response.json") == nil {
			log.Println("Saved subreddits in subreddit_response.json")
		} else {
			log.Println("Failed to save subreddit lists")
		}
	}
}

const REDDIT_URL = "https://oauth.reddit.com"

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

type HttpRequestFailedError int

func (err *HttpRequestFailedError) Error() string {
	return fmt.Sprintf("Request failed with status code: %v", *err)
}

func NewClient(creds *RedditCredentials) RedditClient {
	//TODO: double check if the LastAccessToken is still valid
	return RedditClient{
		creds:      creds,
		httpClient: &http.Client{},
		ctx:        context.Background(),
	}
}

func (client *RedditClient) GetMe() (map[string]any, error) {
	var req = client.buildHttpRequestHolder("GET", REDDIT_URL+"/api/v1/me", nil)
	return client.sendRequest(req)
}

func (client *RedditClient) GetSubreddits() ([]map[string]any, error) {
	var req = client.buildHttpRequestHolder("GET", REDDIT_URL+"/subreddits/mine/subscriber", nil)
	resp, err := client.sendRequest(req)
	if err != nil {
		return nil, err
	}

	children := resp["data"].(map[string]any)["children"].([]any)
	var sr_collection []map[string]any
	for _, v := range children {
		if v.(map[string]any)["kind"].(string) == "t5" {
			data := v.(map[string]any)["data"].(map[string]any)
			sr_collection = append(sr_collection, map[string]any{
				"name":               data["name"],
				"display_name":       data["display_name"],
				"title":              data["title"],
				"subscriber":         int64(data["subscribers"].(float64)),
				"public_description": data["public_description"],
				"category":           data["advertiser_category"],
				"description":        reformatStringValue(data["description"].(string)),
				"already_subscribed": data["user_is_subscriber"].(bool),
			})
		}
	}

	return sr_collection, nil
}

func (client *RedditClient) GetHotPosts(sub_reddit string) ([]map[string]any, error) {
	var url = REDDIT_URL
	if sub_reddit != "" {
		url = url + "/r/" + sub_reddit
	}
	var req = client.buildHttpRequestHolder("GET", url+"/hot", nil)
	resp, err := client.sendRequest(req)
	if err != nil {
		return nil, err
	}

	children := resp["data"].(map[string]any)["children"].([]any)
	var post_collection []map[string]any
	for _, v := range children {
		if v.(map[string]any)["kind"].(string) == "t3" {
			data := v.(map[string]any)["data"].(map[string]any)
			post_collection = append(post_collection, map[string]any{
				"subreddit":                data["subreddit"],
				"num_comments":             int(data["num_comments"].(float64)),
				"title":                    data["title"],
				"upvote_ratio":             data["upvote_ratio"].(float64),
				"contained_url":            data["url"],
				"created":                  time.Unix(int64(data["created"].(float64)), 0),
				"container_sr_subscribers": int64(data["subreddit_subscribers"].(float64)),
				"category":                 data["link_flair_text"],
				"post_score":               int(data["score"].(float64)),
				"post_content":             reformatStringValue(data["selftext"].(string)),
			})
		}
	}
	return post_collection, nil
}

func (client *RedditClient) authenticate() (bool, error) {

	if !client.isTokenExpired() {
		log.Println("Last access token is still valid. No need to reauthenticate")
		return false, nil
	}

	unpwData := url.Values{}
	unpwData.Set("grant_type", "password")
	unpwData.Set("username", client.creds.Username)
	unpwData.Set("password", client.creds.Password)

	req, _ := http.NewRequestWithContext(client.ctx, "POST", "https://www.reddit.com/api/v1/access_token", strings.NewReader(unpwData.Encode()))

	client.attachAuthorizationHeader(req, true)
	client.attachUserAgentName(req)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if respBody, err := client.sendRequest(req); err != nil {
		return false, err
	} else {
		client.creds.LastAccessToken = respBody["access_token"].(string)
		return true, nil
	}
}

func (client *RedditClient) isTokenExpired() bool {
	var jwtToken, _, err = new(jwt.Parser).ParseUnverified(client.creds.LastAccessToken, jwt.MapClaims{})
	if err != nil {
		return true
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return true
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return true
	}
	return float64(time.Now().Unix()) > exp
}

func (client *RedditClient) attachUserAgentName(req *http.Request) string {
	//Windows:My Reddit Bot:1.0 (by u/botdeveloper)
	var agentName = fmt.Sprintf("%v:%v:v0.1 (by u/%v)", runtime.GOOS, client.creds.ClientName, client.creds.Username)
	req.Header.Add("User-Agent", agentName)
	return agentName
}

func (client *RedditClient) attachAuthorizationHeader(req *http.Request, basicAuth bool) string {
	var authValue string
	if basicAuth {
		authValue = fmt.Sprintf("Basic %v", base64.StdEncoding.EncodeToString([]byte(client.creds.ClientId+":"+client.creds.ClientSecret)))
	} else {
		authValue = "bearer " + client.creds.LastAccessToken
	}
	req.Header.Add("Authorization", authValue)
	return authValue
}

// this is a dummy http request builder and does NOT actually verify if the access token is valid
func (client *RedditClient) buildHttpRequestHolder(method string, url string, body io.Reader) *http.Request {
	var req, _ = http.NewRequestWithContext(client.ctx, method, url, body)
	//standard header
	client.attachAuthorizationHeader(req, false)
	client.attachUserAgentName(req)
	if body != nil {
		req.Header.Add("Content-Type", "application/json; charset=UTF-8")
	}
	return req
}

func (client *RedditClient) sendRequest(req *http.Request) (map[string]any, error) {
	resp, err := client.httpClient.Do(req)
	if err != nil {
		//log.Println("Getting New Access Token Failed: ", err)
		return nil, err
	} else if resp.StatusCode == http.StatusUnauthorized {
		//log.Println("Getting New Access Token Failed: ")
		unauthErr := HttpRequestFailedError(resp.StatusCode)
		return nil, &unauthErr
	}
	defer resp.Body.Close()

	return deserialzeJsonBlob[map[string]any](io.Reader(resp.Body))
}

func reformatStringValue(str string) string {
	return str
	//return strings.ReplaceAll(str, "\n", "")
}

func deserialzeJsonBlob[T any](reader io.Reader) (T, error) {
	var decoder = json.NewDecoder(reader)
	var data T
	if err := decoder.Decode(&data); err != nil {
		log.Printf("Error deserializing to data of type %T: %v\n", data, err)
		return data, err
	} else {
		return data, nil
	}
}

func writeDataToJsonFile[T any](data *T, outFile string) error {
	if j_bytes, err := json.Marshal(*data); err != nil {
		return err
	} else {
		return os.WriteFile(outFile, j_bytes, os.FileMode(0644))
	}
}

// loading application configuration. In future making this retrieve from a DB
func loadCredentialsFromFile(configFilePath string) (RedditCredentials, error) {
	var configFile, err = os.Open(configFilePath)
	if err != nil {
		log.Printf("Failed loading configuration file %v. Error: %v\n", configFilePath, err)
		return RedditCredentials{}, err
	}
	defer configFile.Close()

	return deserialzeJsonBlob[RedditCredentials](configFile)
}
