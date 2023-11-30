package socialmediadatastore

type ContentStoreData struct {
	// unique identifier across media source. every reddit item has one. In reddit this is the name
	// in azure cosmos DB every record/item has to have an id.
	// In case of media content the media content itself comes with an unique identifier that we can use
	Id string `json:"id"`
	// unique identifier across item Kind. In reddit this is the "name" without the "tX_" prefix
	UrlId string `json:"url_id"`
	// represents text title of the item. Applies to subreddits and posts but not comments
	Title string `json:"title,omitempty"`
	// Subreddit, Post or Comment. This is not directly serialized
	Kind string `json:"kind"`

	// display_name of the subreddit where the post or comment is in
	Channel string `json:"channel"`
	// Applies to comments and posts.
	// For comments: this represents which post or comment does this comment respond to.
	// for posts: this is the same value as the channel
	Parent string `json:"parent_id,omitempty"`

	//post text
	Text string `json:"text"`
	// for posts this is url posted by the post
	// for subreddit this is link
	Url string `json:"url,omitempty"`

	//subreddit category
	Category string `json:"category,omitempty"`

	// author of posts or comments. Empty for subreddits
	Author string `json:"author,omitempty"`
	// date of creation of the post or comment. Empty for subreddits
	CreatedDate float64 `json:"created,omitempty"`

	// Applies to posts and comments. Doesn't apply to subreddits
	Score int `json:"score,omitempty"`
	// Number of comments to a post or a comment. Doesn't apply to subreddit
	NumComments int `json:"num_comments,omitempty"`
	// Number of subscribers to a channel (subreddit). Doesn't apply to posts or comments
	NumSubscribers int `json:"num_subscribers,omitempty"`
	// Applies to subreddit posts and comments. Doesn't apply to subreddits
	UpvoteRatio float64 `json:"upvote_ratio,omitempty"`
}

type UserActionData struct {
	// in cosmos DB every item has to have an id. Here the id will be synthetic
	// other than azure cosmos DB literally no one cares about this field
	RecordId      string `json:"id"`
	ContentId     string `json:"content_id"`
	Source        string `json:"source"`
	UserId        string `json:"user_id"`
	Processed     bool   `json:"processed,omitempty"`
	Action        string `json:"action,omitempty"`
	ActionContent string `json:"content,omitempty"`
}

type UserMetadata struct {
	UserId    string   `json:"user_id"`
	Interests []string `json:"interests"`
}
