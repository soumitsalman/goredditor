package contentstore

import (
	ctx "context"
	"encoding/json"
	"log"

	cosmos "github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

var content_store_client *cosmos.Client
var content_store_db *cosmos.DatabaseClient
var reddit_store *cosmos.ContainerClient

const MAX_BATCH_SIZE = 99

//var user_interests *cosmos.ContainerClient

type ContentStoreData struct {
	// unique identifier across media source. every reddit item has one
	Name string `json:"name"`
	// unique identifier across item Kind
	Id string `json:"id"`
	// represents text title of the item. Applies to subreddits and posts but not comments
	Title string `json:"title,omitempty"`
	// Subreddit, Post or Comment. This is not directly serialized
	Kind string `json:"kind"`

	// display_name of the subreddit where the post or comment is in
	Channel string `json:"channel"`
	// Applies to comments and posts.
	// For comments: this represents which post or comment does this comment respond to.
	// for posts: this is the same value as the channel
	Parent string `json:"parent_name,omitempty"`

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

func InitializeContentStoreClient() *cosmos.Client {
	if content_store_client != nil {
		// no need to process further
		return content_store_client
	}

	client, err := cosmos.NewClientFromConnectionString(getContentStoreConnection(), nil)
	if err != nil {
		log.Println("Failed connecting to AZ Cosmos DB instance: ", err)
		return nil
	}
	content_store_client = client

	db, err := content_store_client.NewDatabase(getContentStoreDB())
	if err != nil {
		log.Println("Failed finding content store DB: ", err)
		return content_store_client
	}
	content_store_db = db

	container, err := content_store_db.NewContainer(getRedditStoreContainer())
	if err != nil {
		log.Println("Failed finding reddit container: ", err)
		return content_store_client
	}
	reddit_store = container
	return content_store_client
}

// this assumes that all items are of the same kind
// this function upserts instead of insert
func AddNewItems(kind string, items []ContentStoreData) {
	// throttle batch size since comosDB expects less than MAX_BATCH_SIZE operations in a batch
	for len(items) > 0 {
		batch := reddit_store.NewTransactionalBatch(cosmos.NewPartitionKeyString(kind))
		count := min(MAX_BATCH_SIZE, len(items))
		for _, v := range items[0:count] {
			payload, _ := json.Marshal(v)
			batch.UpsertItem(payload, nil)
		}
		// precision is not a target here. If something fails, it can get picked up later for a different user
		if resp, err := reddit_store.ExecuteTransactionalBatch(ctx.Background(), batch, nil); err != nil {
			log.Println("Failed inserting items: ", err)
		} else {
			log.Printf("Status %d. ActivityId %s. Consuming %v Request Units.\n", resp.RawResponse.StatusCode, resp.ActivityID, resp.RequestCharge)
		}
		items = items[count:]
	}

}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
