package redditapplication

import (
	"fmt"
	"log"
	"os"

	"angerproject.org/redditor/utils"
	"github.com/joho/godotenv"
)

const APP_CONFIG_FILE = "appconfig.json"

type RedditorUser struct {
	userId      string
	client      RedditorApplication
	existing_sr []map[string]any
	new_sr      []map[string]any
	new_post    []map[string]any
}

// TODO: return error if it cant find user name variables
func NewUserConnection(userId string) RedditorUser {
	userSession := RedditorUser{userId: userId}
	if config, err := utils.ReadDataFromJsonFile[RedditorCredentials](APP_CONFIG_FILE); err != nil {
		log.Println("Failed loading application config")
		return userSession
	} else {
		userSession.client = NewClient(&config)
	}
	//loading secrets from environment variable
	godotenv.Load()
	//TODO: in future read these from a secret store
	userSession.client.creds.ApplicationId = os.Getenv("GOREDDITOR_APP_ID")
	userSession.client.creds.ApplicationSecret = os.Getenv("GOREDDITOR_APP_SECRET")
	userSession.client.creds.Username = os.Getenv("REDDIT_LOCAL_USER_NAME")
	userSession.client.creds.Password = os.Getenv("REDDIT_LOCAL_USER_PW")
	userSession.client.creds.LastAccessToken = os.Getenv("REDDIT_LOCAL_USER_AUTH_TOKEN")

	return userSession
}

func (user *RedditorUser) GetAreasOfInterest() []string {
	// TODO: load it from a DB
	return []string{"cyber security", "new software products", "software development", "api integration", "generative ai", "software product management", "software program management", "autonomous vehicle", "cloud infrastructure", "information security"}
}

// TODO: remove
// var areas_of_interest = []string{"cyber security", "new software products", "software development", "api integration", "generative ai", "software product management", "software program management", "autonomous vehicle", "cloud infrastructure", "information security"}

func (user *RedditorUser) Authenticate() bool {
	if is_new_token, err := user.client.Authenticate(); err != nil {
		defer log.Printf("Auth failed: %v\n", err)
		return false
	} else if is_new_token {
		defer log.Printf("Got new auth token: \n")
		//save to local env variable for now
		fmt.Println(user.client.creds.LastAccessToken)
		//os.Setenv("REDDIT_LOCAL_USER_AUTH_TOKEN", user.client.creds.LastAccessToken)
	}
	return true
}

func (user *RedditorUser) GetExistingSubreddits() []map[string]any {
	if sr_collection, err := user.client.Subreddits(); err != nil {
		log.Printf("Getting subreddits failed: %v\n", err)
		user.existing_sr = []map[string]any{}
	} else {
		defer utils.SaveToFile(
			user.userId,
			"subscribed_subreddits",
			&sr_collection,
		)
		user.existing_sr = sr_collection
	}
	return user.existing_sr
}

func (user *RedditorUser) GetNewSubreddits() []map[string]any {

	var collection []map[string]any
	if user.existing_sr == nil {
		collection = user.GetExistingSubreddits()
	} else {
		//TODO: used cached subreddits
		collection = user.existing_sr
	}

	// search with areas of interest
	for _, area := range user.GetAreasOfInterest() {
		if res, err := user.client.SubredditSearch(area, -1); err == nil {
			collection = append(collection, res...)
		}
	}

	// collect similar subreddits
	var similar []map[string]any
	for _, sr := range collection {
		if res, err := user.client.SimilarSubreddits(sr["name"].(string)); err == nil {
			similar = append(similar, res...)
		}
	}
	collection = append(collection, similar...)

	defer utils.SaveToFile(
		user.userId,
		"recommended_subreddits",
		&collection,
	)
	user.new_sr = collection
	return user.new_sr
}

func (user *RedditorUser) GetNewPosts() []map[string]any {

	var collection []map[string]any // this is the value to be return

	// prepping the scope of subreddits to search for.
	var sr_in_scope = []string{""}
	for _, v := range user.existing_sr {
		sr_in_scope = append(sr_in_scope, v["display_name"].(string))
	}

	// for subreddit in scope each post type iterate for each
	for _, subreddit := range sr_in_scope {
		for _, pt := range []string{"hot", "top", "best"} {
			if post_collection, err := user.client.Posts(subreddit, pt); err != nil {
				log.Printf("Getting %v post from r/%v failed: %v\n", pt, subreddit, err)
			} else {
				collection = append(collection, post_collection...)
				log.Printf("Retrieved %v posts from r/%v\n", pt, subreddit)
			}
		}
	}
	// save it in a file
	defer utils.SaveToFile(
		user.userId,
		"posts",
		&collection,
	)
	user.new_post = collection
	return user.new_post
}
