package main

import (
	"log"

	"angerproject.org/redditor/redditclient"
	"angerproject.org/redditor/utils"
)

const data_dump_folder = "C:\\Users\\soumi\\go-stuff\\reddit_data_dump\\"
const config_file = data_dump_folder + "config.json"

func authenticate(client *redditclient.RedditClient) bool {
	if is_new_token, err := client.Authenticate(); err != nil {
		log.Printf("Auth failed: %v\n", err)
		return false
	} else if is_new_token {
		log.Printf("Got new auth token: \n")
		redditclient.SaveClientToConfigFile(client, config_file)
	}
	return true
}

func collectSubscribedSubreddits(client *redditclient.RedditClient) {
	if sr_collection, err := client.Subreddits(); err != nil {
		log.Printf("Getting subreddits failed: %v\n", err)
	} else {
		var filename = data_dump_folder + "joined_subreddits.json"
		if utils.WriteDataToJsonFile(&sr_collection, filename) == nil {
			log.Println("Saved subreddits in " + filename)
		} else {
			log.Println("Failed to save subreddit lists")
		}
	}
}

func collectRecommendedSubreddits(client *redditclient.RedditClient) {
	if recommended_sr, err := client.SimilarSubreddits(); err != nil {
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

func collectPosts(client *redditclient.RedditClient) {
	var post_types = []string{"hot", "top", "best"}
	var collection []map[string]any

	for _, pt := range post_types {
		if post_collection, err := client.Posts("", pt); err != nil {
			log.Printf("Getting %v failed: %v\n", pt, err)
		} else {
			collection = append(collection, post_collection...)
			log.Printf("Retrieved %v posts\n", pt)
		}
	}

	var filename = data_dump_folder + "posts.json"
	if utils.WriteDataToJsonFile(&collection, filename) == nil {
		log.Println("Saved posts results in ", filename)
	} else {
		log.Println("Failed to save posts")
	}
}

// primary orchestrator
func main() {

	client, _ := redditclient.NewClientFromConfigFile(config_file)

	if !authenticate(&client) {
		return
	}

	//daily collection
	collectSubscribedSubreddits(&client)
	collectRecommendedSubreddits(&client)
	collectPosts(&client)

	/*
		sr_name := "reddit_api_test"
		if resp, err := client.JoinSubreddit(sr_name); !resp {
			log.Printf("Failed joining subreddit [%v]: %v\n", sr_name, err)
		} else {
			log.Printf("Successfully joined subreddit: [%v]\n", sr_name)
		}
	*/
	/*
		sr_name := "reddit_api_test_dump"
		if resp, err := client.Submit("Test Text Post from randomizer_000", "Text content", sr_name); err != nil {
			log.Printf("Failed posting text in subreddit [%v]: %v\n", sr_name, err)
		} else {
			log.Printf("Successfully posted text in subreddit [%v]: %v\n", sr_name, resp)
		}
	*/
	/*
		if resp, err := client.SubmitPost("Test Link Post from randomizer_000", "https://platform.openai.com/playground", sr_name); err != nil {
			log.Printf("Failed posting link in subreddit [%v]: %v\n", sr_name, err)
		} else {
			log.Printf("Successfully posted link subreddit: [%v]: %v\n", sr_name, resp)
		}
	*/
	/*
		post_id := "t3_17l29uh"
		if resp, err := client.Comment(fmt.Sprintf("test comment from randomizer_000 at %v", time.Now()), post_id); !resp {
			log.Printf("Failed commenting on post [%v]: %v\n", post_id, err)
		} else {
			log.Printf("Successfully commented on post [%v]\n", post_id)
		}
	*/
}
