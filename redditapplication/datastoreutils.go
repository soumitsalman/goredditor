package redditapplication

import (
	"database/sql"
	"log"

	"angerproject.org/redditor/utils"
	configdb "github.com/replit/database-go"
)

// <START DATASTORE related functions.>
// TODO: push the content to DB

type RedditorDataStore struct {
}

func getDataStoreLocation() string {
	//application config
	loc, _ := configdb.Get("datastore_location")
	return loc
}

func getDataStoreConnection() *sql.DB {
	db, err := sql.Open("sqlite3", getDataStoreLocation())
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func closeDataStoreConnection(db *sql.DB) {
	db.Close()
}

func saveNewRedditData(db *sql.DB, item_kind string, value string) {
	db.ex
}

func saveNewData[T RedditData[T]](userId string, topic string, data []T) {
	//saving the blob
	var content = map[string]any{
		"topic": topic,
		topic:   data,
	}
	var filename = getDataStoreLocation() + userId + "_" + topic + ".json"
	if utils.WriteDataToJsonFile(&content, filename) == nil {
		log.Printf("Saved %s in %s\n", topic, filename)
	} else {
		log.Printf("Failed to save %s\n", topic)
	}
	//saving the state
	saveStateData[T](userId, data, NEW)
}

func readExistingData[T any](userId string, topic string) (T, error) {
	var filename = getDataStoreLocation() + userId + "_" + topic + ".json"
	return utils.ReadDataFromJsonFile[T](filename)
}

func saveStateData[T RedditData[T]](userId string, data []T, state string) {
	topic := "states"
	list := make(map[string]string)
	if states, err := readExistingData[map[string]string](userId, topic); err == nil {
		//this means there is already some states content
		list = states
	}

	for _, v := range data {
		name := v.GetUniqueName()
		if sval, ok := list[name]; ok {
			//if it does exist then push
			list[name] = newerStateValue(state, sval)
		} else {
			//update with the {state} value
			list[name] = state
		}
	}
	var filename = getDataStoreLocation() + userId + "_states" + ".json"
	if utils.WriteDataToJsonFile(&list, filename) == nil {
		log.Printf("Saved %s in %s\n", topic, filename)
	} else {
		log.Printf("Failed to save %s\n", topic)
	}
}

func newerStateValue(v1, v2 string) string {
	if v1 > v2 {
		return v1
	}
	return v2
}

const (
	NEW              = "0_new"
	INTERESTING      = "1_interesting"
	SHORT_LISTED     = "2_short_listed"
	ACTION_SUGGESTED = "3_action_suggested"
	ACTION_TAKEN     = "4_action_taken"
	IGNORE           = "9_ignore"
)

// <END DATASTORE related functionality>
