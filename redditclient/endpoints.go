package redditclient

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"runtime"
	"slices"
	"strings"

	"angerproject.org/redditor/utils"
)

const REDDIT_DATA_URL = "https://oauth.reddit.com"
const REDDIT_AUTH_URL = "https://www.reddit.com/api/v1/access_token"

func (client *RedditClient) GetMe() (map[string]any, error) {
	var req = client.buildRetrievalRequest(REDDIT_DATA_URL + "/api/v1/me")
	return utils.GetHttpResponse(req, client.httpClient)
}

// get subreddits that user is already part of
func (client *RedditClient) GetSubreddits() ([]map[string]any, error) {
	var req = client.buildRetrievalRequest(REDDIT_DATA_URL + "/subreddits/mine/subscriber")
	if resp, err := utils.GetHttpResponse(req, client.httpClient); err != nil {
		return nil, err
	} else {
		return extractSubreddits(resp), nil
	}
}

// get subreddits that are suggested based on existing subscritpion
// does not return unique list of items and may have duplicates
func (client *RedditClient) GetRecommendedSubreddits() ([]map[string]any, error) {

	current_sr_collection, err := client.GetSubreddits()
	if err != nil {
		return nil, err
	}

	var sr_collection []map[string]any
	for _, sr := range current_sr_collection {
		var req = client.buildRetrievalRequest(REDDIT_DATA_URL + "/api/similar_subreddits?sr_fullnames=" + sr["name"].(string))
		if resp, err := utils.GetHttpResponse(req, client.httpClient); err != nil {
			return nil, err
		} else {
			sr_collection = append(sr_collection, extractSubreddits(resp)[:]...)
		}
	}
	return sr_collection, nil
}

func (client *RedditClient) GetPosts(sub_reddit string, post_type string) ([]map[string]any, error) {
	var url = REDDIT_DATA_URL

	//url correctness check
	if !slices.Contains([]string{"hot", "best", "top"}, post_type) {
		post_type = "hot"
		log.Printf("Invalid post_type %v. Going with hot posts\n", post_type)
	}
	if utils.HasWhiteSpace(sub_reddit) {
		log.Printf("Invalid sub-reddit name: %v. Going with overall profile's %v posts\n", sub_reddit, post_type)
	} else if sub_reddit == "" {
		log.Printf("No sub-reddit specified . Going with overall profile's %v posts\n", post_type)
	} else {
		//valid subreddit name will make valid url
		url = url + "/r/" + sub_reddit
	}

	var req = client.buildRetrievalRequest(url + "/" + post_type)
	if resp, err := utils.GetHttpResponse(req, client.httpClient); err != nil {
		return nil, err
	} else {
		return extractPosts(resp), nil
	}
}

// authenticates client. Checks if the last known access token is still valid. If not it makes a re-authentication request
// return true if there was a new auth token or else returns false.
// error is returned when auth was failed
func (client *RedditClient) Authenticate() (bool, error) {

	if !utils.IsAuthTokenExpired(client.creds.LastAccessToken) {
		log.Println("Last access token is still valid. No need to reauthenticate")
		return false, nil
	}

	unpwData := url.Values{}
	unpwData.Set("grant_type", "password")
	unpwData.Set("username", client.creds.Username)
	unpwData.Set("password", client.creds.Password)

	req := utils.BuildHttpRequest(
		"POST",                               //this needs to be a post
		REDDIT_AUTH_URL,                      //this is the auth url
		strings.NewReader(unpwData.Encode()), //this is my username and pw, NOT the application id and secret
		"url",                                //my username and pw data has to be url encoding
		utils.MakeBasicAuthToken(client.creds.ClientId, client.creds.ClientSecret), //client id and client secret is the auth for the request
		client.getApplicationFullName(),                                            //application name
		&client.ctx,                                                                //context
	)

	if respBody, err := utils.GetHttpResponse(req, client.httpClient); err != nil {
		return false, err
	} else {
		client.creds.LastAccessToken = respBody["access_token"].(string)
		return true, nil
	}
}

func (client *RedditClient) buildRetrievalRequest(endpoint_url string) *http.Request {
	return utils.BuildHttpRequest(
		"GET",        //method
		endpoint_url, //uri
		nil,          //no need for payload
		"",           //no need for payload encoding
		utils.MakeBearerToken(client.creds.LastAccessToken), //assign auth token
		client.getApplicationFullName(),                     //assign user-agent name
		&client.ctx,                                         //assigning context
	)
}

func (client *RedditClient) getApplicationFullName() string {
	//Windows:My Reddit Bot:1.0 (by u/botdeveloper)
	return fmt.Sprintf("%v:%v:v0.1 (by u/%v)", runtime.GOOS, client.creds.ClientName, client.creds.Username)
}
