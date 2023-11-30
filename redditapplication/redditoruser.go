package redditapplication

import (
	"log"
	"strings"
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
		log.Printf("Retrieved subscribed subreddits, Count: %d\n", len(sr_collection))
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
	user.new_sr = collection
	return user.new_sr
}

func (user *RedditorUser) LoadNewPosts() []RedditData {

	var collection []RedditData // this is the value to be return

	/* TODO: delete this. it becomes unnecessary since i am collecting the tops from each subreddit anywway
	// prepping the scope of subreddits to search for.
	var sr_in_scope = []string{""}
	for _, v := range user.existing_sr {
		sr_in_scope = append(sr_in_scope, v.DisplayName)
	}
	*/
	// for subreddit in scope each post type iterate for each
	for _, subreddit := range user.existing_sr {
		for _, pt := range []string{"hot", "top", "best"} {
			if post_collection, err := user.client.Posts(subreddit.DisplayName, pt); err != nil {
				log.Printf("Getting %s post from r/%s failed: %v\n", pt, subreddit.DisplayName, err)
			} else {
				collection = append(collection, post_collection...)
				log.Printf("Retrieved %s posts from r/%s\n", pt, subreddit.DisplayName)
			}
		}
	}
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
	user.new_comments = collection
	return user.new_comments
}

func (user *RedditorUser) SaveNewFilteredContents() {
	// TODO: do some pre-filtering like if the subreddit has less than 500 people, dump it
	if user.new_sr != nil {
		saveNewItemsToDB(user.Id, SUBREDDIT, user.new_sr)
	}
	// TODO: do some pre-filtering like if the post is older than 2 weeks, dump it
	if user.new_posts != nil {
		saveNewItemsToDB(user.Id, POST, user.new_posts)
	}
	// TODO: do some pre-filtering like if the comment is older than  weeks, dump it
	if user.new_comments != nil {
		saveNewItemsToDB(user.Id, COMMENT, user.new_comments)
	}
}

func (user *RedditorUser) TakeUserActions() {
	ua_data, content_data := getUserActionsAndContents()
	for i := 0; i < len(ua_data); i++ {
		if strings.ToLower(ua_data[i].Action) == "sub" {
			user.client.Subscribe(content_data[i].Channel)
		}
	}
}
