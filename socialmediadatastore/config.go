package socialmediadatastore

import (
	"os"
	"strconv"
	"time"
)

func getContentStoreConnection() string {
	return os.Getenv("AZ_COSMOSDB_CONNECTION")
}

func getContentStoreDB() string {
	return os.Getenv("CONTENT_STORE_DB")
}

func getRedditStoreContainer() string {
	return os.Getenv("REDDIT_STORE_CONTAINER")
}

func getUserMetadataContainer() string {
	return os.Getenv("USER_METADATA_CONTAINER")
}

func getUserActionContainer() string {
	return os.Getenv("USER_CONTENT_CONTAINER")
}

func getServiceBusConnection() string {
	return os.Getenv("AZ_SERVICE_BUS_CONNECTION")
}

func getNewItemsQueue() string {
	return os.Getenv("NEW_ITEMS_QUEUE")
}

func getInterestingItemsQueue() string {
	return os.Getenv("INTERESTING_ITEMS_QUEUE")
}

func getShortListedItemsQueue() string {
	return os.Getenv("SHORTLISTED_ITEMS_QUEUE")
}

func getUserActionsQueue() string {
	return os.Getenv("USER_ACTION_QUEUE")
}

func getMaxBatchSize() int {
	val, _ := strconv.Atoi(os.Getenv("MAX_BATCH_SIZE"))
	return val
}

func getMaxWaitTime() time.Duration {
	val, _ := strconv.Atoi(os.Getenv("MAX_WAIT_TIME"))
	return time.Duration(val) * time.Second
}
