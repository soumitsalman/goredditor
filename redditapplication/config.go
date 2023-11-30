package redditapplication

import (
	"os"
	"strconv"
)

func getAppName() string {
	return os.Getenv("GOREDDITOR_APP_NAME")
}

func getAppDescription() string {
	return os.Getenv("GOREDDITOR_APP_DESCRIPTION")
}

func getAboutUrl() string {
	return os.Getenv("GOREDDITOR_ABOUT_URL")
}

func getRedirectUri() string {
	return os.Getenv("GOREDDITOR_REDIRECT_URI")
}

func getAppId() string {
	return os.Getenv("GOREDDITOR_APP_ID")
}

func getAppSecret() string {
	return os.Getenv("GOREDDITOR_APP_SECRET")
}

func getLocalUserName() string {
	return os.Getenv("REDDIT_LOCAL_USER_NAME")
}

func getLocalUserPw() string {
	return os.Getenv("REDDIT_LOCAL_USER_PW")
}

func getDataStoreLocation() string {
	return os.Getenv("DATASTORE_LOCATION")
}

func getMaxTextSize() int {
	val, _ := strconv.Atoi(os.Getenv("MAX_TEXT_SIZE"))
	return val
}
