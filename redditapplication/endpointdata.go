package redditapplication

import (
	"strings"

	"angerproject.org/redditor/utils"
)

const REDDIT_PRETTY_URL = "www.reddit.com"

type AuthenticationData struct {
	AuthToken string `json:"access_token"`
}

type SubredditData struct {
	Name                  string `json:"name"`                    //unique name with t5_ prefix
	DisplayName           string `json:"display_name"`            //url name
	Title                 string `json:"title"`                   //fancy long name
	SubscriberCount       int    `json:"subscriber"`              //number of subscribers
	PublicDescription     string `json:"public_description_html"` //short description
	Description           string `json:"description_html"`        //long description
	Category              string `json:"advertiser_category"`     //optional category of the subreddit
	UserAlreadySubscribed bool   `json:"already_subscribed"`
	Link                  string `json:"url"` //link to the subreddit page
}

type PostData struct {
	Name        string  `json:"name"`            //name with t3_ prefix
	Id          string  `json:"id"`              //this the name without the t3_ prefix
	Title       string  `json:"title"`           //title of the post
	Subreddit   string  `json:"subreddit"`       //subreddit that contains this post
	Author      string  `json:"author"`          //author of the post
	Text        string  `json:"selftext_html"`   //text content of the post
	Url         string  `json:"url"`             //URL that was submitted as the post (which applies when Text is blank)
	CreatedDate float64 `json:"created"`         //time it was created on. this is a floating point representation will need to be transformed
	UpvoteRatio float32 `json:"upvote_ratio"`    //up_vote to total vote ratio
	Score       int     `json:"score"`           //post score. the higher the better
	NumComments int     `json:"num_comments"`    //number of comments in the post
	Category    string  `json:"link_flair_text"` //optional category of the post defined by the author
	Link        string  `json:"permalink"`       //reddit.com url of the post
}

type CommentData struct {
	Id          string  `json:"id"`        //comment id. This is the {name} without the t1_ prefix
	Name        string  `json:"name"`      //comment name with t1_{id}
	Text        string  `json:"body_html"` //comment content
	Author      string  `json:"author"`    //authoer of the comment
	Score       int     `json:"score"`     //score of the comment. the higher the better
	CreatedDate float64 `json:"created"`   //floating point representation of the datetime the comment was created
	PostId      string  `json:"link_id"`   //the overarching post this comment or its parent comments are part of
	Link        string  `json:"permalink"` //reddit.com url of the comment
	//Replies     any     `json:"replies"` //the reply can be either string or an array of Listing

}

// this is primarily a wrapper in which reddit sends the response
type ListingData[T any] struct {
	Data struct {
		Children []struct {
			Data T `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

type UnmappedData map[string]any

func extractFromListing[T any](resp ListingData[T]) []T {
	items := make([]T, len(resp.Data.Children))
	for i, v := range resp.Data.Children {
		items[i] = v.Data
	}
	return items
}

func extractFromListingArray[T any](resp []ListingData[T]) []T {
	var collection []T
	for _, listing := range resp {
		collection = append(collection, extractFromListing[T](listing)...)
	}
	return collection
}

//<START data retrieval post processing >

type RedditData[T any] interface {
	GetUniqueName() string
	GetNormalizedData() T
}

// returns the name field which is the unique identifier across reddit
func (post PostData) GetUniqueName() string {
	return post.Name
}

// returns the name field which is the unique identifier across reddit
func (sr SubredditData) GetUniqueName() string {
	return sr.Name
}

// returns the name field which is the unique identifier across reddit
func (comm CommentData) GetUniqueName() string {
	return comm.Name
}

// removes html tags in comment text and prefixes www.reddit.com in front of the link
func (comment CommentData) GetNormalizedData() CommentData {
	comment.Text = utils.ExtractTextFromHtml(comment.Text)
	comment.Link = ensureRedditDotCom(comment.Link)
	return comment
}

// removes html tags in post text and prefixes www.reddit.com in front of the link
func (post PostData) GetNormalizedData() PostData {
	post.Text = utils.ExtractTextFromHtml(post.Text)
	post.Link = ensureRedditDotCom(post.Link)
	return post
}

// removes html tags in subreddit descriptions and prefixes www.reddit.com in front of the link
func (subreddit SubredditData) GetNormalizedData() SubredditData {
	subreddit.PublicDescription = utils.ExtractTextFromHtml(subreddit.PublicDescription)
	subreddit.Description = utils.ExtractTextFromHtml(subreddit.Description)
	subreddit.Link = ensureRedditDotCom(subreddit.Link)
	return subreddit
}

// running data normalization on an array of items
func normalizeDataList[T RedditData[T]](items []T) []T {
	for i := range items {
		items[i] = items[i].GetNormalizedData()
	}
	return items
}

// prefixes www.reddit.com to the URL it it is not there
func ensureRedditDotCom(url string) string {
	if !strings.HasPrefix(url, REDDIT_PRETTY_URL) {
		return REDDIT_PRETTY_URL + url
	}
	return url
}

//<END data retrieval post processing >
