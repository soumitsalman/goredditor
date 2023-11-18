package main

import (
	"fmt"

	rapp "angerproject.org/redditor/redditapplication"
)

// this is for pure data collection
func collectContents(user *rapp.RedditorUser) {
	user.LoadExistingSubreddits()
	user.LoadNewSubreddits()
	user.LoadNewPosts()
	user.LoadNewComments()
}

// this is for ONLY making posts and subscribing to new subreddits
func takeActions(user *rapp.RedditorUser) {
	fmt.Println("TODO: implement it")

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

// primary orchestrator
func main() {

	user := rapp.NewUserConnection("soumitsr@gmail.com")

	if !user.Authenticate() {
		return
	}

	//daily collection
	collectContents(&user)
	takeActions(&user)

}
