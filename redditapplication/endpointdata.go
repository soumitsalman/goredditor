package redditapplication

type AuthenticationData struct {
	AuthToken string `json:"access_token"`
}

type SubredditData struct {
	Name                  string `json:"name"`                //unique name with t5_ prefix
	DisplayName           string `json:"display_name"`        //url name
	Title                 string `json:"title"`               //fancy long name
	SubscriberCount       int    `json:"subscriber"`          //number of subscribers
	PublicDesciption      string `json:"public_description"`  //short description
	Description           string `json:"description"`         //long description
	Category              string `json:"advertiser_category"` //optional category of the subreddit	
	UserAlreadySubscribed bool `json:"already_subscribed"`
}

type PostData struct {
	Name     string `json:"name"` //name with t3_ prefix
	Id string `json:id` //this the name without the t3_ prefix
	Title        string            `json:"title"` //title of the post
	Subreddit    string            `json:"subreddit"` //subreddit that contains this post
	Author   string `json:"author"` //author of the post	
	Text  string `json:"selftext"` //text content of the post
	Url      string `json:"url"` //either the URL of the post itself of the URL that was submitted as the post (which applies when Text is blank)	
	CreatedDate   float64 `json:"created"` //time it was created on. this is a floating point representation will need to be transformed
	UpvoteRatio float32           `json:"upvote_ratio"` //up_vote to total vote ratio
	Score    int    `json:"score"` //post score. the higher the better
	NumComments  int               `json:"num_comments"` //number of comments in the post
	Category string `json:"link_flair_text"` //optional category of the post defined by the author
}

type CommentData struct {
	Id string `json:id` //comment id. This is the {name} without the t1_ prefix
	Name string `json:name` //comment name with t1_{id}
	Text   string `json:"body"` //comment content
	Author string `json:"author"` //authoer of the comment
	Score  int    `json:"score"` //score of the comment. the higher the better
	CreatedDate float64 `json:"created"` //floating point representation of the datetime the comment was created
	PostId string `json:"link_id"` //the overarching post this comment or its parent comments are part of
}

type ListingData[T any] struct {
	Data struct {
		Children []struct {
			Data T `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

type UnmappedData map[string]any

func extractFromOneListing[T any](resp ListingData[T]) []T {
	items := make([]T, len(resp.Data.Children))
	for i, v := range resp.Data.Children {
		items[i] = v.Data
	}
	return items
}

func extractFromMultipleListing[T any](resp []ListingData[T]) []T {
	var collection []T

	for _,listing := range resp {
		collection = append(collection, extractFromOneListing[T](listing)...)
	}
	return collection
}

func getChildren(resp map[string]any) []any {
	var collection []any
	if resp != nil {
		if data, ok := resp["data"]; ok {
			if children, ok := data.(map[string]any)["children"]; ok {
				collection = children.([]any)
			}
		}
	}
	return collection
}

/*
func extractPosts(resp map[string]any) []map[string]any {
	var collection []map[string]any
	for _, ch := range getChildren(resp) {
		v := ch.(map[string]any)
		if v["kind"].(string) == "t3" {
			if data, ok := v["data"].(map[string]any); ok {
				collection = append(collection, map[string]any{
					"subreddit":                data["subreddit"],
					"num_comments":             int(data["num_comments"].(float64)),
					"title":                    data["title"],
					"upvote_ratio":             data["upvote_ratio"].(float64),
					"created":                  time.Unix(int64(data["created"].(float64)), 0),
					"container_sr_subscribers": int64(data["subreddit_subscribers"].(float64)),
					"category":                 data["link_flair_text"],
					"post_score":               data["score"],
					"post_content":             data["selftext"],
					"name":                     data["name"],
					"author":                   data["author"],
					"url":                      data["url"],
					"id": data["id"],
				})
			}
		}
	}

	return collection
}

func extractSubreddits(resp map[string]any) []map[string]any {
	var collection []map[string]any
	for _, ch := range getChildren(resp) {
		v := ch.(map[string]any)
		if v["kind"].(string) == "t5" {
			data := v["data"].(map[string]any)

			collection = append(collection, map[string]any{
				"name":               data["name"],                //unique name with t5_ prefix
				"display_name":       data["display_name"],        //url name
				"title":              data["title"],               //fancy long name
				"subscriber":         data["subscribers"],         //number of subscribers
				"public_description": data["public_description"],  //short description
				"category":           data["advertiser_category"], //optional category of the subreddit
				"description":        data["description"],         //long description
				"already_subscribed": data["user_is_subscriber"],  //if i am already subscribed
			})
		}
	}
	return collection
}
*/