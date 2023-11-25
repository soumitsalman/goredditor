package socialmediadatastore

import (
	ctx "context"
	"encoding/json"
	"fmt"
	"log"

	cosmos "github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
)

var content_store_client *cosmos.Client
var content_store_db *cosmos.DatabaseClient
var reddit_store *cosmos.ContainerClient
var user_metadata *cosmos.ContainerClient
var user_action_store *cosmos.ContainerClient

const MAX_BATCH_SIZE = 99

func InitializeContentStoreClient() *cosmos.Client {
	if content_store_client != nil {
		// no need to process further
		return content_store_client
	}

	// establish connection
	client, err := cosmos.NewClientFromConnectionString(getContentStoreConnection(), nil)
	if err != nil {
		log.Println("Failed connecting to AZ Cosmos DB instance: ", err)
		return nil
	}
	content_store_client = client

	// get instance of the DB
	db, err := content_store_client.NewDatabase(getContentStoreDB())
	if err != nil {
		log.Println("Failed finding content store DB: ", err)
		return nil
	}
	content_store_db = db

	// get instance of the reddit content store
	c0, err := content_store_db.NewContainer(getRedditStoreContainer())
	if err != nil {
		log.Println("Failed finding reddit store container: ", err)
	}
	reddit_store = c0

	// get instance of user metadata
	c1, err := content_store_db.NewContainer(getUserMetadataContainer())
	if err != nil {
		log.Println("Failed finding user metadata container: ", err)
	}
	user_metadata = c1

	// get instance of user specific content action
	c2, err := content_store_db.NewContainer(getUserActionContainer())
	if err != nil {
		log.Println("Failed finding user action container: ", err)
	}
	user_action_store = c2

	return content_store_client
}

// TODO: delete later. This is solely for debugging purpose
func addToUserContentStore__V2(items []UserActionData) {
	// throttle batch size since comosDB expects less than MAX_BATCH_SIZE operations in a batch

	for _, item := range items {
		payload, _ := json.Marshal(item)
		if resp, err := user_action_store.CreateItem(ctx.Background(), cosmos.NewPartitionKeyString(item.UserId), payload, nil); err != nil {
			log.Println("Failed inserting items: ", err)
		} else {
			log.Printf("Status %d. ActivityId %s. Consuming %v Request Units. %v\n", resp.RawResponse.StatusCode, resp.ActivityID, resp.RequestCharge, resp.Response)
		}
	}
}

// this assumes that all items are of the same kind
// this function upserts instead of insert
func AddToContentStore(items []ContentStoreData) {
	if len(items) > 0 {
		addInBatches[ContentStoreData](reddit_store, items[0].Kind, items)
	}
}

// this assumes that all items are of the same user_id
// this function upserts instead of insert
func AddToUserActionStore(items []UserActionData) {
	if len(items) > 0 {
		addInBatches[UserActionData](user_action_store, items[0].UserId, items)
	}
}

func GetExistingUserActionsContentIds(user_id string, source string) []string {
	query := "SELECT c.content_id FROM c WHERE c.user_id=@user_id AND c.source=@source"
	q_opt := cosmos.QueryOptions{
		QueryParameters: []cosmos.QueryParameter{
			{Name: "@user_id", Value: user_id},
			{Name: "@source", Value: source},
			//{Name: "@existing_list", Value: []string{"t3_182w4cz", "t3_182y2qf", "t3_182d4bi"}},
		},
	}
	var result []string
	resp_pager := user_action_store.NewQueryItemsPager(query, cosmos.NewPartitionKeyString(user_id), &q_opt)
	for resp_pager.More() {
		if resp, err := resp_pager.NextPage(ctx.Background()); err == nil {
			for _, item := range resp.Items {
				var data map[string]string
				json.Unmarshal(item, &data)
				result = append(result, data["content_id"])
				//fmt.Println(data.ContentId)
			}
		} else {
			log.Println(err)
		}
	}
	return result
}

func _TEST_query(user_id string) {

	query := "SELECT c.content_id FROM c WHERE NOT ARRAY_CONTAINS(@existing_list,  c.content_id)"
	q_opt := cosmos.QueryOptions{
		QueryParameters: []cosmos.QueryParameter{
			{Name: "@existing_list", Value: []string{"t3_182w4cz", "t3_182y2qf", "t3_182d4bi"}},
		},
	}
	resp_pager := user_action_store.NewQueryItemsPager(query, cosmos.NewPartitionKeyString(user_id), &q_opt)
	for resp_pager.More() {
		if resp, err := resp_pager.NextPage(ctx.Background()); err == nil {
			for _, item := range resp.Items {
				var data UserActionData
				json.Unmarshal(item, &data)
				fmt.Println(data.ContentId)
			}
		} else {
			log.Println(err)
		}
	}
}

func addInBatches[T any](container *cosmos.ContainerClient, partition_key string, items []T) {
	// throttle batch size since comosDB expects less than MAX_BATCH_SIZE operations in a batch
	for len(items) > 0 {
		batch := container.NewTransactionalBatch(cosmos.NewPartitionKeyString(partition_key))
		count := min(getMaxBatchSize(), len(items))
		for _, v := range items[0:count] {
			payload, _ := json.Marshal(v)
			batch.UpsertItem(payload, nil)
		}
		// precision is not a target here. If something fails, it can get picked up later for a different user
		if resp, err := container.ExecuteTransactionalBatch(ctx.Background(), batch, nil); err != nil {
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
