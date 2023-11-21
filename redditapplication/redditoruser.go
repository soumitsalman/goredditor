package redditapplication

import (
	"log"
	"os"

	"angerproject.org/redditor/utils"
)

const APP_CONFIG_FILE = "appconfig.json"

type RedditorUser struct {
	Id           string
	client       RedditorApplication
	existing_sr  []SubredditData
	new_sr       []SubredditData
	new_posts    []PostData
	new_comments []CommentData
}

// TODO: return error if it cant find user name variables
func NewUserConnection(userId string) RedditorUser {
	userSession := RedditorUser{Id: userId}
	if config, err := utils.ReadDataFromJsonFile[RedditorCredentials](APP_CONFIG_FILE); err != nil {
		log.Println("Failed loading application config")
		return userSession
	} else {
		userSession.client = NewClient(&config)
	}
	//TODO: remove the dotenv loading since replit is handling the environment variable
	//loading secrets from environment variable
	//godotenv.Load()
	//TODO: in future read these from a secret store
	userSession.client.creds.ApplicationId = os.Getenv("GOREDDITOR_APP_ID")
	userSession.client.creds.ApplicationSecret = os.Getenv("GOREDDITOR_APP_SECRET")
	userSession.client.creds.Username = os.Getenv("REDDIT_LOCAL_USER_NAME")
	userSession.client.creds.Password = os.Getenv("REDDIT_LOCAL_USER_PW")
	userSession.client.creds.OauthToken = os.Getenv("REDDIT_LOCAL_USER_AUTH_TOKEN")

	return userSession
}

func (user *RedditorUser) GetAreasOfInterest() []string {
	// TODO: load it from a DB
	return []string{"cyber security", "new software products", "software development", "api integration", "generative ai", "software product management", "software program management", "autonomous vehicle", "cloud infrastructure", "information security"}
}

func (user *RedditorUser) Authenticate() string {
	if _, err := user.client.Authenticate(); err != nil {
		defer log.Printf("Auth failed: %v\n", err)
		return ""
	}
	return user.client.creds.OauthToken
}

func (user *RedditorUser) LoadExistingSubreddits() []SubredditData {
	if sr_collection, err := user.client.Subreddits(); err != nil {
		log.Printf("Getting subreddits failed: %v\n", err)
	} else {
		defer saveNewData[SubredditData](user.Id, "subscribed_subreddits", sr_collection)
		user.existing_sr = sr_collection
	}
	return user.existing_sr
}

func (user *RedditorUser) LoadNewSubreddits() []SubredditData {

	var collection []SubredditData

	// search with areas of interest
	for _, area := range user.GetAreasOfInterest() {
		if res, err := user.client.SubredditSearch(area, -1); err == nil {
			collection = append(collection, res...)
		}
	}

	// collect similar subreddits
	var similar []SubredditData
	for _, sr := range collection {
		if res, err := user.client.SimilarSubreddits(sr.Name); err == nil {
			similar = append(similar, res...)
		}
	}
	collection = append(collection, similar...)

	defer saveNewData[SubredditData](user.Id, "recommended_subreddits", collection)
	user.new_sr = collection
	return user.new_sr
}

func (user *RedditorUser) LoadNewPosts() []PostData {

	var collection []PostData // this is the value to be return

	// prepping the scope of subreddits to search for.
	var sr_in_scope = []string{""}
	for _, v := range user.existing_sr {
		sr_in_scope = append(sr_in_scope, v.DisplayName)
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
	defer saveNewData[PostData](user.Id, "posts", collection)
	user.new_posts = collection
	return user.new_posts
}

// loads comments of the posts that have already been loaded
// there is a bug in this. it loads some metadata of the parent post as well.
// Also it does not go deepder than 1 layer of comments. The reddit data doesn't work well with json marshalling
// TODO: filter out parent post

func (user *RedditorUser) LoadNewComments() []CommentData {
	var collection []CommentData

	for _, post := range user.new_posts {
		if comments, err := user.client.RetrieveComments(post); err != nil {
			log.Println("Failed retrieving comments for ", post.Name)
		} else {
			collection = append(collection, comments...)
		}
	}

	defer saveNewData[CommentData](user.Id, "comments", collection)
	user.new_comments = collection
	return collection
}
