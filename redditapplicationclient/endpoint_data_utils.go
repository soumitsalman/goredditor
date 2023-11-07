package redditapplicationclient

import "time"

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
					"contained_url":            data["url"],
					"created":                  time.Unix(int64(data["created"].(float64)), 0),
					"container_sr_subscribers": int64(data["subreddit_subscribers"].(float64)),
					"category":                 data["link_flair_text"],
					"post_score":               int(data["score"].(float64)),
					"post_content":             data["selftext"].(string),
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
