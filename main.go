package main

import (
	"fmt"
	"log"
	"time"

	"angerproject.org/redditor/redditclient"
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

	/*
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
	*/
	/*
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

	post_id := "t3_17l29uh"
	if resp, err := client.Comment(fmt.Sprintf("test comment from randomizer_000 at %v", time.Now()), post_id); !resp {
		log.Printf("Failed commenting on post [%v]: %v\n", post_id, err)
	} else {
		log.Printf("Successfully commented on post [%v]\n", post_id)
	}

}
