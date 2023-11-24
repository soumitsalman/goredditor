package dataprocessingqueue

import (
	"os"
	"strconv"
)

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
