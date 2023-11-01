package redditclient

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"slices"

	"angerproject.org/redditor/utils"
)

const REDDIT_DATA_URL = "https://oauth.reddit.com"
const REDDIT_AUTH_URL = "https://www.reddit.com/api/v1/access_token"

type SubmissionError struct {
	messages []any
}

func (err *SubmissionError) Error() string {
	return fmt.Sprintf("Submission Errors: %v", err.messages)
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

	req := client.buildPostActionReqeustWithUrlEncoding(REDDIT_AUTH_URL, unpwData)
	if respBody, err := utils.SendHttpRequest(req, client.httpClient); err != nil {
		return false, err
	} else {
		client.creds.LastAccessToken = respBody["access_token"].(string)
		return true, nil
	}
}

// gets my details
func (client *RedditClient) Me() (map[string]any, error) {
	var req = client.buildGetActionRequest(REDDIT_DATA_URL + "/api/v1/me")
	return utils.SendHttpRequest(req, client.httpClient)
}

// get subreddits that user is already part of
func (client *RedditClient) Subreddits() ([]map[string]any, error) {
	var req = client.buildGetActionRequest(REDDIT_DATA_URL + "/subreddits/mine/subscriber")
	if resp, err := utils.SendHttpRequest(req, client.httpClient); err != nil {
		return nil, err
	} else {
		return extractSubreddits(resp), nil
	}
}

// get subreddits that are suggested based on existing subscritpion
// does not return unique list of items and may have duplicates
func (client *RedditClient) SimilarSubreddits() ([]map[string]any, error) {

	current_sr_collection, err := client.Subreddits()
	if err != nil {
		return nil, err
	}

	var sr_collection []map[string]any
	for _, sr := range current_sr_collection {
		// this has to queried using unique name of the subreddit (e.g t5_ffffff)
		var req = client.buildGetActionRequest(REDDIT_DATA_URL + "/api/similar_subreddits?sr_fullnames=" + sr["name"].(string))
		if resp, err := utils.SendHttpRequest(req, client.httpClient); err != nil {
			return nil, err
		} else {
			sr_collection = append(sr_collection, extractSubreddits(resp)[:]...)
		}
	}
	return sr_collection, nil
}

// gets my posts: hot, best and top depending what is specified through post_type
// if sub_reddit display name is not specified it will pull from the overall list of posts instead of a specific subreddit
func (client *RedditClient) Posts(sub_reddit string, post_type string) ([]map[string]any, error) {
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

	var req = client.buildGetActionRequest(url + "/" + post_type)
	if resp, err := utils.SendHttpRequest(req, client.httpClient); err != nil {
		return nil, err
	} else {
		return extractPosts(resp), nil
	}
}

// joins a subreddit. sr_name should be the display name of the subreddit. NOT the unique id
func (client *RedditClient) Subscribe(sr_name string) (bool, error) {
	data := url.Values{}
	data.Set("action", "sub")
	data.Set("skip_inital_defaults", "true")
	data.Set("sr_name", sr_name)

	var req = client.buildPostActionReqeustWithUrlEncoding(REDDIT_DATA_URL+"/api/subscribe", data)

	if _, err := utils.SendHttpRequest(req, client.httpClient); err != nil {
		return false, err
	} else {
		return true, nil
	}
}

// submits a post with a specified title and text_content in a given subreddit
// sr_name should be the display name and not the unique name
// if text_context is a valid url then it will automatically submit it as a link
func (client *RedditClient) Submit(title string, text_content string, sr_name string) (map[string]any, error) {
	data := url.Values{}
	data.Set("api_type", "json")
	data.Set("sr", sr_name)
	data.Set("title", title)
	if utils.IsValidUrl(text_content) {
		//it is a valid url and so post the URL
		data.Set("url", text_content)
		data.Set("kind", "link")
	} else {
		//it is a self-text
		data.Set("text", text_content)
		data.Set("kind", "self")
	}

	var req = client.buildPostActionReqeustWithUrlEncoding(REDDIT_DATA_URL+"/api/submit", data)
	if resp, err := utils.SendHttpRequest(req, client.httpClient); err != nil {
		return nil, err
	} else {
		resp_data := resp["json"].(map[string]any)
		sub_errors := resp_data["errors"].([]any)
		if len(sub_errors) == 0 {
			return resp_data["data"].(map[string]any), nil
		} else {
			return nil, &SubmissionError{messages: sub_errors}
		}
	}
}

// posts comment_text as comment to a given post or comment (parent_name)
// parent_name has to be the unique id with t3_ or t1_
// TODO: currently it is returning bool. Change to return the metadata of the comment
func (client *RedditClient) Comment(comment_text string, parent_name string) (bool, error) {
	data := url.Values{}
	data.Set("parent", parent_name)
	data.Set("text", comment_text)

	var req = client.buildPostActionReqeustWithUrlEncoding(REDDIT_DATA_URL+"/api/comment", data)
	if resp, err := utils.SendHttpRequest(req, client.httpClient); err != nil {
		return false, err
	} else {
		return resp["success"].(bool), nil
	}
}

// internal function wrapping over an httputils function
func (client *RedditClient) buildGetActionRequest(endpoint_url string) *http.Request {
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

func (client *RedditClient) buildPostActionReqeustWithUrlEncoding(endpoint_url string, payload url.Values) *http.Request {
	return utils.BuildHttpRequest(
		"POST",
		endpoint_url,
		utils.SerializeUrlValues(payload),
		"url",
		utils.MakeBearerToken(client.creds.LastAccessToken),
		client.getApplicationFullName(),
		&client.ctx,
	)
}

/*
func (client *RedditClient) buildPostActionReqeust(endpoint_url string, payload io.Reader, payload_encoding string) *http.Request {
	return utils.BuildHttpRequest(
		"POST",
		endpoint_url,
		payload,
		payload_encoding,
		utils.MakeBearerToken(client.creds.LastAccessToken),
		client.getApplicationFullName(),
		&client.ctx,
	)
}
*/
