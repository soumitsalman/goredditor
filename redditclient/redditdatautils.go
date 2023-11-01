package redditclient

import "time"

func extractPosts(resp map[string]any) []map[string]any {
	children := resp["data"].(map[string]any)["children"].([]any)
	var post_collection []map[string]any
	for _, v := range children {
		if v.(map[string]any)["kind"].(string) == "t3" {
			data := v.(map[string]any)["data"].(map[string]any)
			post_collection = append(post_collection, map[string]any{
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
	return post_collection
}

func extractSubreddits(resp map[string]any) []map[string]any {
	children := resp["data"].(map[string]any)["children"].([]any)
	var sr_collection []map[string]any
	for _, v := range children {
		if v.(map[string]any)["kind"].(string) == "t5" {
			data := v.(map[string]any)["data"].(map[string]any)
			sr_collection = append(sr_collection, map[string]any{
				"name":               data["name"],                         //unique name with t5_ prefix
				"display_name":       data["display_name"],                 //url name
				"title":              data["title"],                        //fancy long name
				"subscriber":         int64(data["subscribers"].(float64)), //number of subscribers
				"public_description": data["public_description"],           //short description
				"category":           data["advertiser_category"],          //optional category of the subreddit
				"description":        data["description"].(string),         //long description
				"already_subscribed": data["user_is_subscriber"].(bool),    //if i am already subscribed
			})
		}
	}
	return sr_collection
}
