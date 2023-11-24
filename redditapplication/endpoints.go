package redditapplication

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/soumitsalman/goredditor/utils"
)

const REDDIT_DATA_URL = "https://oauth.reddit.com"
const REDDIT_AUTH_URL = "https://www.reddit.com/api/v1/access_token"

type SubmissionError struct {
	messages []any
}

func (err *SubmissionError) Error() string {
	return fmt.Sprintf("Submission Errors: %v", err.messages)
}

// authenticates client. Checks if the last known access token is still valid.
// If not it makes a re-authentication request
// returns the auth token in case of success.
// error is returned when auth was failed
func (client *RedditorApplication) Authenticate() (string, error) {

	if !utils.IsAuthTokenExpired(client.creds.OauthToken) {
		log.Println("Last access token is still valid. No need to reauthenticate")
		return client.creds.OauthToken, nil
	}

	log.Println("Getting new auth token")

	unpwData := url.Values{}
	unpwData.Set("grant_type", "password")
	unpwData.Set("username", client.creds.Username)
	unpwData.Set("password", client.creds.Password)

	req := utils.BuildHttpRequest(
		"POST",
		REDDIT_AUTH_URL,
		utils.SerializeUrlValues(unpwData),
		"url",
		utils.MakeBasicAuthToken(client.creds.ApplicationId, client.creds.ApplicationSecret),
		client.getApplicationFullName(),
		&client.ctx,
	)

	if respBody, err := utils.SendHttpRequest[AuthenticationData](req, client.httpClient); err != nil {
		return "", err
	} else {
		client.creds.OauthToken = respBody.AuthToken
		return client.creds.OauthToken, nil
	}
}

// <START Retrieval-Codes>

// gets my details
func (client *RedditorApplication) Me() (map[string]any, error) {
	var req = client.buildGetActionRequest(REDDIT_DATA_URL + "/api/v1/me")
	return utils.SendHttpRequest[map[string]any](req, client.httpClient)
}

// get subreddits that user is already part of
func (client *RedditorApplication) Subreddits() ([]RedditData, error) {
	var req = client.buildGetActionRequest(REDDIT_DATA_URL + "/subreddits/mine/subscriber")
	if resp, err := utils.SendHttpRequest[ListingData](req, client.httpClient); err != nil {
		return nil, err
	} else {
		return extractFromListing(resp), nil
	}
}

// get subreddits based on a given
// does not return unique list of items and may have duplicates
func (client *RedditorApplication) SimilarSubreddits(sr_name string) ([]RedditData, error) {
	var req = client.buildGetActionRequest(REDDIT_DATA_URL + "/api/similar_subreddits?sr_fullnames=" + sr_name)
	if resp, err := utils.SendHttpRequest[ListingData](req, client.httpClient); err != nil {
		return nil, err
	} else {
		return extractFromListing(resp), nil
	}
}

// uses the query string to look for sub-reddits
// min_users is used to filter for sub-reddits that has at least min_users number of users
func (client *RedditorApplication) SubredditSearch(search_query string, min_users int) ([]RedditData, error) {
	q, err := url.Parse(search_query)
	if err != nil {
		log.Printf("Invalid search query %v\n", err)
		return nil, err
	}
	search_str := q.String()
	var req = client.buildGetActionRequest(REDDIT_DATA_URL + "/subreddits/search?q=" + search_str)
	if resp, err := utils.SendHttpRequest[ListingData](req, client.httpClient); err != nil {
		return nil, err
	} else {
		// TODO: filter for min_users
		return extractFromListing(resp), nil
	}
}

// gets my posts: hot, best and top depending what is specified through post_type
// if sub_reddit display name is not specified it will pull from the overall list of posts instead of a specific subreddit
func (client *RedditorApplication) Posts(sub_reddit string, post_type string) ([]RedditData, error) {
	var url = REDDIT_DATA_URL

	//url correctness check
	/* TODO: check for parameter
	if !slices.Contains([]string{"hot", "best", "top"}, post_type) {
		post_type = "hot"
		log.Printf("Invalid post_type %v. Going with hot posts\n", post_type)
	}
	*/

	if utils.HasWhiteSpace(sub_reddit) {
		log.Printf("Invalid sub-reddit name: %v. Going with overall profile's %v posts\n", sub_reddit, post_type)
	} else if sub_reddit == "" {
		log.Printf("No sub-reddit specified . Going with overall profile's %v posts\n", post_type)
	} else {
		//valid subreddit name will make valid url
		url = url + "/r/" + sub_reddit
	}

	var req = client.buildGetActionRequest(url + "/" + post_type)
	if resp, err := utils.SendHttpRequest[ListingData](req, client.httpClient); err != nil {
		return nil, err
	} else {
		return extractFromListing(resp), nil
	}
}

func (client *RedditorApplication) RetrieveComments(post RedditData) ([]RedditData, error) {
	url := fmt.Sprintf("%s/r/%s/comments/%s", REDDIT_DATA_URL, post.Subreddit, post.Id)
	var req = client.buildGetActionRequest(url)

	if resp, err := utils.SendHttpRequest[[]ListingData](req, client.httpClient); err != nil {
		return nil, err
	} else {

		return extractFromListingArray(resp), nil
	}
}

// <END Retrieval-Code>

// <START State-Modifying-Code>

// joins a subreddit. sr_name should be the display name of the subreddit. NOT the unique id
func (client *RedditorApplication) Subscribe(sr_name string) (bool, error) {
	data := url.Values{}
	data.Set("action", "sub")
	data.Set("skip_inital_defaults", "true")
	data.Set("sr_name", sr_name)

	var req = client.buildPostActionReqeustWithUrlEncoding(REDDIT_DATA_URL+"/api/subscribe", data)

	if _, err := utils.SendHttpRequest[map[string]any](req, client.httpClient); err != nil {
		return false, err
	} else {
		return true, nil
	}
}

// submits a post with a specified title and text_content in a given subreddit
// sr_name should be the display name and not the unique name
// if text_context is a valid url then it will automatically submit it as a link
func (client *RedditorApplication) Submit(title string, text_content string, sr_name string) (map[string]any, error) {
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
	if resp, err := utils.SendHttpRequest[map[string]any](req, client.httpClient); err != nil {
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
func (client *RedditorApplication) Comment(comment_text string, parent_name string) (bool, error) {
	data := url.Values{}
	data.Set("parent", parent_name)
	data.Set("text", comment_text)

	var req = client.buildPostActionReqeustWithUrlEncoding(REDDIT_DATA_URL+"/api/comment", data)
	if resp, err := utils.SendHttpRequest[map[string]any](req, client.httpClient); err != nil {
		return false, err
	} else {
		return resp["success"].(bool), nil
	}
}

// <END State-Modifying-Code>

// <START wrapper on httputils>

// internal function wrapping over an httputils function
func (client *RedditorApplication) buildGetActionRequest(endpoint_url string) *http.Request {
	return utils.BuildHttpRequest(
		"GET",        //method
		endpoint_url, //uri
		nil,          //no need for payload
		"",           //no need for payload encoding
		utils.MakeBearerToken(client.creds.OauthToken), //assign auth token
		client.getApplicationFullName(),                //assign user-agent name
		&client.ctx,                                    //assigning context
	)
}

func (client *RedditorApplication) buildPostActionReqeustWithUrlEncoding(endpoint_url string, payload url.Values) *http.Request {
	return utils.BuildHttpRequest(
		"POST",
		endpoint_url,
		utils.SerializeUrlValues(payload),
		"url",
		utils.MakeBearerToken(client.creds.OauthToken),
		client.getApplicationFullName(),
		&client.ctx,
	)
}

//<END wrapper on httputils>
