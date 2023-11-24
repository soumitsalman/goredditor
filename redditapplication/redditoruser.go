package redditapplication

import (
	"log"
)

type RedditorUser struct {
	Id           string
	client       RedditorApplication
	existing_sr  []RedditData
	new_sr       []RedditData
	new_posts    []RedditData
	new_comments []RedditData
}

// TODO: return error if it cant find user name variables
func NewUserConnection(userId string) RedditorUser {
	var creds RedditorCredentials = RedditorCredentials{
		//TODO: in future read these from a secret store
		ApplicationName:        getAppName(),
		ApplicationDescription: getAppDescription(),
		AboutUrl:               getAboutUrl(),
		RedirectUri:            getRedirectUri(),
		ApplicationId:          getAppId(),
		ApplicationSecret:      getAppSecret(),
		Username:               getLocalUserName(),
		Password:               getLocalUserPw(),
	}
	userSession := RedditorUser{
		Id:     userId,
		client: NewClient(&creds),
	}
	initializeDataStores()

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

func (user *RedditorUser) LoadExistingSubreddits() []RedditData {
	if sr_collection, err := user.client.Subreddits(); err != nil {
		log.Printf("Getting subreddits failed: %v\n", err)
	} else {
		// TODO: save to a table
		user.existing_sr = sr_collection
	}
	return user.existing_sr
}

func (user *RedditorUser) LoadNewSubreddits() []RedditData {

	var collection []RedditData

	// search with areas of interest
	for _, area := range user.GetAreasOfInterest() {
		if res, err := user.client.SubredditSearch(area, -1); err == nil {
			collection = append(collection, res...)
		}
	}

	// collect similar subreddits
	var similar []RedditData
	for _, sr := range collection {
		if res, err := user.client.SimilarSubreddits(sr.Name); err == nil {
			similar = append(similar, res...)
		}
	}
	collection = append(collection, similar...)

	defer saveNewItemsToDB(user.Id, SUBREDDIT, collection)
	user.new_sr = collection
	return user.new_sr
}

func (user *RedditorUser) LoadNewPosts() []RedditData {

	var collection []RedditData // this is the value to be return

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
	defer saveNewItemsToDB(user.Id, POST, collection)
	user.new_posts = collection
	return user.new_posts
}

// loads comments of the posts that have already been loaded
// there is a bug in this. it loads some metadata of the parent post as well.
// Also it does not go deepder than 1 layer of comments. The reddit data doesn't work well with json marshalling
// TODO: filter out parent post

func (user *RedditorUser) LoadNewComments() []RedditData {
	var collection []RedditData

	for _, post := range user.new_posts {
		if comments, err := user.client.RetrieveComments(post); err != nil {
			log.Println("Failed retrieving comments for ", post.Name)
		} else {
			collection = append(collection, comments...)
		}
	}

	defer saveNewItemsToDB(user.Id, COMMENT, collection)
	user.new_comments = collection
	return collection
}
