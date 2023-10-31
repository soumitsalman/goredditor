package main

import (
	"log"

	"angerproject.org/redditor/redditclient"
	"angerproject.org/redditor/utils"
)

const data_dump_folder = "C:\\Users\\soumi\\go-stuff\\reddit_data_dump\\"
const config_file = data_dump_folder + "config.json"

func main() {

	client, _ := redditclient.NewClientFromConfigFile(config_file)

	if is_new_token, err := client.Authenticate(); err != nil {
		log.Printf("Auth failed: %v\n", err)
		return
	} else if is_new_token {
		log.Printf("Got new auth token: \n")
		redditclient.SaveClientToConfigFile(&client, config_file)
	}

	var post_type = "hot"
	if post_collection, err := client.GetPosts("", post_type); err != nil {
		log.Printf("Getting subreddits failed: %v\n", err)
	} else {
		var filename = data_dump_folder + post_type + "_posts.json"
		if utils.WriteDataToJsonFile(&post_collection, filename) == nil {
			log.Println("Saved hot post results in ", filename)
		} else {
			log.Println("Failed to save hot posts")
		}
	}

	if sr_collection, err := client.GetSubreddits(); err != nil {
		log.Printf("Getting subreddits failed: %v\n", err)
	} else {
		var filename = data_dump_folder + "joined_subreddits.json"
		if utils.WriteDataToJsonFile(&sr_collection, filename) == nil {
			log.Println("Saved subreddits in " + filename)
		} else {
			log.Println("Failed to save subreddit lists")
		}
	}

	if recommended_sr, err := client.GetRecommendedSubreddits(); err != nil {
		log.Printf("Getting recommended subreddits failed: %v\n", err)
	} else {
		var filename = data_dump_folder + "recommended_subreddits.json"
		if utils.WriteDataToJsonFile(&recommended_sr, filename) == nil {
			log.Println("Saved recommened subreddits in " + filename)
		} else {
			log.Println("Failed to save recommended subreddit lists")
		}
	}
}
