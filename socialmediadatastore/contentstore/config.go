package contentstore

import "os"

func getContentStoreConnection() string {
	return os.Getenv("AZ_COSMOSDB_CONNECTION")
}

func getContentStoreDB() string {
	return os.Getenv("CONTENT_STORE_DB")
}

func getRedditStoreContainer() string {
	return os.Getenv("REDDIT_STORE_CONTAINER")
}
