package redditapplication

type AuthenticationData struct {
	AuthToken string `json:"access_token"`
}

const (
	SUBREDDIT = "subreddit"
	POST      = "post"
	COMMENT   = "comment"
)

// this is primarily a wrapper in which reddit sends the response
type ListingData struct {
	Data struct {
		Children []struct {
			Kind string     `json:"kind"`
			Data RedditData `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

// represents Subreddit, Channels, Forums, Users, Posts, Comments
type RedditData struct {
	Name        string `json:"name"`         // unique identifier across media source. every reddit item has one
	DisplayName string `json:"display_name"` //url name for subreddits
	Id          string `json:"id"`           // unique identifier across item Kind
	Title       string `json:"title"`        // represents text title of the item. Applies to subreddits and posts but not comments
	// Subreddit, Post or Comment. This is not directly serialized
	Kind string

	// display_name of the subreddit where the post or comment is in
	Subreddit string `json:"subreddit"`
	// Applies to comments and posts.
	// For comments: this represents which post or comment does this comment respond to.
	// for posts: this is the same value as the channel
	Parent string `json:"link_id"`

	// comment body
	CommentBody string `json:"body_html"`
	// post text
	PostText string `json:"selftext_html"`
	// for posts this is url posted by the post
	// for subreddit this is link
	Url string `json:"url"`
	//subreddit short description
	PublicDescription string `json:"public_description_html"`
	//subreddit long description
	Description string `json:"description_html"`
	//subreddit category
	SubredditCategory string `json:"advertiser_category"`
	// optional author or creator defined category of the post topic or subreddit topic
	PostCategory string `json:"link_flair_text"`
	// url or link to the content item in the media source
	Link string `json:"permalink"`

	// author of posts or comments. Empty for subreddits
	Author string `json:"author"`
	// date of creation of the post or comment. Empty for subreddits
	CreatedDate float64 `json:"created"`

	// Applies to posts and comments. Doesn't apply to subreddits
	Score int `json:"score"`
	// Number of comments to a post or a comment. Doesn't apply to subreddit
	NumComments int `json:"num_comments"`
	// Number of subscribers to a channel (subreddit). Doesn't apply to posts or comments
	NumSubscribers int `json:"subscribers"`
	// this applies to posts and comments to indicate the same thing as above
	SubredditSubscribers int `json:"subreddit_subscribers"`
	// Applies to subreddit posts and comments. Doesn't apply to subreddits
	UpvoteRatio float64 `json:"upvote_ratio"`
}

func extractFromListing(resp ListingData) []RedditData {
	items := make([]RedditData, len(resp.Data.Children))
	for i, v := range resp.Data.Children {
		items[i] = v.Data
		items[i].Kind = convertKind(v.Kind)
	}
	return items
}

func extractFromListingArray(resp []ListingData) []RedditData {
	var collection []RedditData
	for _, listing := range resp {
		collection = append(collection, extractFromListing(listing)...)
	}
	return collection
}

func convertKind(kind string) string {
	switch kind {
	case "t5":
		return SUBREDDIT
	case "t3":
		return POST
	case "t1":
		return COMMENT
	default:
		return ""
	}
}
